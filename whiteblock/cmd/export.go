package cmd

import (
	util "../util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"net/url"
	"os"
	"sync"
	"syscall"
)

func handleFetchChunk(testnetID string, node Node, log string, chunk string) (string, error) {
	ep := fmt.Sprintf("https://api.whiteblock.io/testnets/%s/nodes/%s/logs/%s/chunks/%s", testnetID, node.ID, log, chunk)
	fmt.Println(ep)
	return util.JwtHTTPRequest("GET", ep, "")
}

func handleChunks(testnetID string, node Node, log string, rawChunks string, sem *semaphore.Weighted) []string {
	var res map[string]interface{}
	err := json.Unmarshal([]byte(rawChunks), &res)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	chunks := res["items"].([]interface{})
	wg := sync.WaitGroup{}
	wg.Add(len(chunks))

	ctx := context.TODO()

	outChunks := make([]string, len(chunks))

	for i, chunk := range chunks {
		sem.Acquire(ctx, 1)
		go func(chunk string, i int) {
			defer sem.Release(1)
			defer wg.Done()
			var err error
			var res string
			for j := 0; j < 10; j++ {
				res, err = handleFetchChunk(testnetID, node, log, chunk)
				if err == nil {
					break
				}
			}
			if err != nil {
				util.PrintErrorFatal(err)
			}
			fmt.Printf("%d is done\n", i)
			err = ioutil.WriteFile(fmt.Sprintf("./%s/%s/%s", node.ID, log, chunk), []byte(res), 0664)
			if err != nil {
				util.PrintErrorFatal(err)
			}

		}(chunk.(string), i)
		outChunks[i] = chunk.(string)
	}
	wg.Wait()

	return outChunks
}

func handleExportLogs(testnetID string, node Node, rawRes string, sem *semaphore.Weighted) (interface{}, map[string][]string) {
	var res map[string]interface{}
	err := json.Unmarshal([]byte(rawRes), &res)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	logs := res["items"].([]interface{})
	out := map[string][]string{}

	for _, log := range logs {
		os.RemoveAll(fmt.Sprintf("./%s/%v", node.ID, log))
		os.MkdirAll(fmt.Sprintf("./%s/%v", node.ID, log), 0755)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(logs))
	mux := sync.Mutex{}

	for _, log := range logs {
		go func(log string) {
			defer wg.Done()
			ep := fmt.Sprintf("https://api.whiteblock.io/testnets/%s/nodes/%s/logs/%v/chunks", testnetID, node.ID, log)
			fmt.Println(ep)
			res, err := util.JwtHTTPRequest("GET", ep, "")
			if err != nil {
				util.PrintErrorFatal(err)
			}
			mux.Lock()
			out[log] = handleChunks(testnetID, node, log, res, sem)
			mux.Unlock()
		}(log.(string))

	}
	wg.Wait()
	return res["nextPageToken"], out
}

func handleExportBlocks(testnetID string, node string, rawRes string, sem *semaphore.Weighted) (interface{}, []string) {
	var res map[string]interface{}
	err := json.Unmarshal([]byte(rawRes), &res)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	blockNumbers := res["items"].([]interface{})
	out := make([]string, len(blockNumbers))

	wg := sync.WaitGroup{}
	wg.Add(len(blockNumbers))
	mux := sync.Mutex{}

	ctx := context.TODO()
	for i, blockNumber := range blockNumbers {
		sem.Acquire(ctx, 1)
		go func(blockNumber interface{}, i int) {
			defer wg.Done()
			defer sem.Release(1)
			ep := fmt.Sprintf("https://api.whiteblock.io/testnets/%s/nodes/%s/blocks/%v", testnetID, node, blockNumber)
			fmt.Println(ep)
			var res string
			var err error
			for i := 0; i < 10; i++ {
				res, err = util.JwtHTTPRequest("GET", ep, "")
				if err == nil {
					break
				}
			}
			if err != nil {
				util.PrintErrorFatal(err)
			}
			mux.Lock()
			out[i] = res
			fmt.Printf("%d is done\n", i)
			mux.Unlock()
		}(blockNumber, i)

	}
	wg.Wait()
	return res["nextPageToken"], out
}

func mergeDown(node Node, files map[string][]string) {
	for file, chunks := range files {
		tmpFileName := fmt.Sprintf("./%s/%s.tmp", node.ID, file)
		fd, err := os.Create(tmpFileName)
		if err != nil {
			util.PrintErrorFatal(err)
		}
		defer fd.Close()
		for _, chunk := range chunks {
			data, err := ioutil.ReadFile(fmt.Sprintf("./%s/%s/%s", node.ID, file, chunk))
			if err != nil {
				util.PrintErrorFatal(err)
			}
			_, err = fd.Write(data)
			if err != nil {
				util.PrintErrorFatal(err)
			}
		}
		os.RemoveAll(fmt.Sprintf("./%s/%s", node.ID, file))

		syscall.Rename(fmt.Sprintf("./%s/%s.tmp", node.ID, file), fmt.Sprintf("./%s/%s", node.ID, file))
	}

}

func appendBlocks(items []string, finish bool, f *os.File) {

	fInfo, err := f.Stat()
	if err != nil {
		util.PrintErrorFatal(err)
	}
	if fInfo.Size() == 0 {
		f.Write([]byte("[")) //TODO check err
	}

	for i, item := range items {
		if fInfo.Size() > 1 || i != 0 {
			f.Write([]byte(",")) //TODO check err
		}
		_, err := f.Write([]byte(item))
		if err != nil {
			util.PrintErrorFatal(err)
		}
	}

	if finish {
		f.Write([]byte("]"))
	}
}

var exportCmd = &cobra.Command{
	Hidden: true,
	Use:    "export [testnet id]",
	Short:  "Export stuff",
	Long:   "Export stuff",
	Run: func(cmd *cobra.Command, args []string) {
		var testnetID string
		var err error
		if len(args) == 0 {
			testnetID, err = getPreviousBuildId()
			if err != nil {
				util.PrintErrorFatal(err)
			}
		} else {
			testnetID = args[0]
		}
		nodes := []Node{}
		sem := semaphore.NewWeighted(200)
		ep := fmt.Sprintf("https://api.whiteblock.io/testnets/%s/nodes", testnetID)
		res, err := util.JwtHTTPRequest("GET", ep, "")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		err = json.Unmarshal([]byte(res),&nodes)
		for _, node := range nodes {
			os.RemoveAll(fmt.Sprintf("./%s", node.ID))
			os.MkdirAll(fmt.Sprintf("./%s", node.ID), 0755)
		}
		fmt.Println("removed the files")
		wg := sync.WaitGroup{}
		wg.Add(len(nodes))

		for _, node := range nodes {
			go func(node Node) {
				defer wg.Done()
				var nextToken interface{}
				var files map[string][]string
				for {
					var ep string
					if nextToken == nil {
						ep = fmt.Sprintf("https://api.whiteblock.io/testnets/%s/nodes/%s/logs", testnetID, node.ID)
					} else {
						ep = fmt.Sprintf("https://api.whiteblock.io/testnets/%s/nodes/%s/logs?next=%v",
							testnetID, node.ID, url.QueryEscape(nextToken.(string)))
					}

					fmt.Println(ep)
					res, err := util.JwtHTTPRequest("GET", ep, "")
					if err != nil {
						util.PrintErrorFatal(err)
					}
					fmt.Println(res)
					nextToken, files = handleExportLogs(testnetID, node, res, sem)
					mergeDown(node, files)
					if nextToken == nil {
						break
					}
				}
			}(node)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			var nextToken interface{}
			var blocks []string

			files := []*os.File{}
			for _, node := range nodes {

				f, err := os.Create(fmt.Sprintf("./%s/blocks.json", node.ID))
				if err != nil {
					util.PrintErrorFatal(err)
				}
				files = append(files, f)
			}
			defer func() {
				for _, file := range files {
					file.Close()
				}
			}()

			for {
				var ep string
				if nextToken == nil {
					ep = fmt.Sprintf("https://api.whiteblock.io/testnets/%s/nodes/00000000-0000-0000-0000-000000000000/blocks", testnetID)
				} else {
					ep = fmt.Sprintf("https://api.whiteblock.io/testnets/%s/nodes/00000000-0000-0000-0000-000000000000/blocks?next=%v",
						testnetID, url.QueryEscape(nextToken.(string)))
				}

				fmt.Println(ep)
				res, err := util.JwtHTTPRequest("GET", ep, "")
				if err != nil {
					util.PrintErrorFatal(err)
				}
				nextToken, blocks = handleExportBlocks(testnetID, "00000000-0000-0000-0000-000000000000", res, sem)
				for _, file := range files {
					appendBlocks(blocks, nextToken == nil, file)
				}

				if nextToken == nil {
					break
				}
			}
		}()
		wg.Wait()
		/*for _,node := range nodes {
			ep := fmt.Sprintf("https://api.whiteblock.io/testnets/%s/nodes/%s/blocks",testnetID,node.ID)
			fmt.Println(ep)
			res,err := util.JwtHTTPRequest("GET",ep,"")
			fmt.Println(res,err)
		}*/
		//fmt.Printf("https://api.whiteblock.io/testnets/%s/nodes/%s/blocks\n",testnetID,node.ID)
		//	fmt.Printf("https://api.whiteblock.io/testnets/%s/nodes/%s/logs\n",testnetID,node.ID)
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
