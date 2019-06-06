package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	util "github.com/whiteblock/cli/whiteblock/util"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"sync"
	"syscall"
)

func handleFetchChunk(testnetID string, node Node, logName string, chunk string) (string, error) {
	ep := fmt.Sprintf("%s/testnets/%s/nodes/%s/logs/%s/chunks/%s", conf.APIURL, testnetID, node.ID, logName, chunk)
	log.WithFields(log.Fields{"ep": ep, "chunk": chunk}).Trace("fetching the log chunk")
	return util.JwtHTTPRequest("GET", ep, "")
}

func handleChunks(testnetID string, node Node, logName string, rawChunks string, sem *semaphore.Weighted) []string {
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
				res, err = handleFetchChunk(testnetID, node, logName, chunk)
				if err == nil {
					break
				}
			}
			if err != nil {
				util.PrintErrorFatal(err)
			}
			log.WithFields(log.Fields{"chunk": chunk, "num": i}).Debug("fetched a chunk")
			err = ioutil.WriteFile(fmt.Sprintf("./%s/%s/%s", node.ID, logName, chunk), []byte(res), 0664)
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

	for _, logName := range logs {
		os.RemoveAll(fmt.Sprintf("./%s/%v", node.ID, logName))
		os.MkdirAll(fmt.Sprintf("./%s/%v", node.ID, logName), 0755)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(logs))
	mux := sync.Mutex{}

	for _, logName := range logs {
		go func(logName string) {
			defer wg.Done()
			ep := fmt.Sprintf("%s/testnets/%s/nodes/%s/logs/%v/chunks", conf.APIURL, testnetID, node.ID, logName)
			log.WithFields(log.Fields{"ep": ep}).Debug("fetching the log chunks")
			res, err := util.JwtHTTPRequest("GET", ep, "")
			if err != nil {
				util.PrintErrorFatal(err)
			}
			mux.Lock()
			out[logName] = handleChunks(testnetID, node, logName, res, sem)
			mux.Unlock()
		}(logName.(string))

	}
	wg.Wait()
	return res["nextPageToken"], out
}

func convertBlockNumber(blockNumber interface{}) int64 {
	switch num := blockNumber.(type) {
	case float64:
		return int64(num)
	case string:
		res, err := strconv.ParseInt(num, 0, 64)
		if err != nil {
			util.PrintErrorFatal(err)
		}
		return res
	default:
		util.PrintErrorFatal(fmt.Errorf("blocknumber is of unknown type"))
	}
	panic("shouldn't reach")
}

func handleExportBlocks(testnetID string, node string, rawRes string, coveredBlockNumbers *map[int64]struct{}, sem *semaphore.Weighted) (interface{}, []string) {
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

		num := convertBlockNumber(blockNumber)
		if _, ok := (*coveredBlockNumbers)[num]; ok {
			wg.Done()
			continue
		}
		(*coveredBlockNumbers)[num] = struct{}{}
		sem.Acquire(ctx, 1)
		go func(blockNumber interface{}, i int) {
			defer wg.Done()
			defer sem.Release(1)
			ep := fmt.Sprintf("%s/testnets/%s/nodes/%s/blocks/%v", conf.APIURL, testnetID, node, blockNumber)
			log.WithFields(log.Fields{"ep": ep}).Debug("fetching the block data")
			var res string
			var err error
			for it := 0; it < 10; it++ {
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
			log.WithFields(log.Fields{"num": i, "blockNumber": blockNumber}).Trace("fetched a block")
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

func appendBlocks(items []string, finish bool, firstCall bool, f *os.File) {

	fInfo, err := f.Stat()
	if err != nil {
		util.PrintErrorFatal(err)
	}
	if fInfo.Size() == 0 {
		f.Write([]byte("[")) //TODO check err
	}
	first := true
	for _, item := range items {
		if len(item) == 0 {
			continue
		}
		if first && firstCall && fInfo.Size() < 2 {
			first = false
		} else {
			_, err := f.Write([]byte(",")) //TODO check err
			if err != nil {
				util.PrintErrorFatal(err)
			}
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

		spinner := Spinner{txt: "fetching the block and log data"}
		spinner.Run(100)
		defer spinner.Kill()

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
		ep := fmt.Sprintf("%s/testnets/%s/nodes", conf.APIURL, testnetID)
		res, err := util.JwtHTTPRequest("GET", ep, "")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		err = json.Unmarshal([]byte(res), &nodes)
		if err != nil {
			util.PrintErrorFatal(err)
		}
		for _, node := range nodes {
			os.RemoveAll(fmt.Sprintf("./%s", node.ID))
			os.MkdirAll(fmt.Sprintf("./%s", node.ID), 0755)
		}
		log.Trace("removed the files")
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
						ep = fmt.Sprintf("%s/testnets/%s/nodes/%s/logs", conf.APIURL, testnetID, node.ID)
					} else {
						ep = fmt.Sprintf("%s/testnets/%s/nodes/%s/logs?next=%v", conf.APIURL,
							testnetID, node.ID, url.QueryEscape(nextToken.(string)))
					}
					res, err := util.JwtHTTPRequest("GET", ep, "")
					if err != nil {
						util.PrintErrorFatal(err)
					}
					log.WithFields(log.Fields{"ep": ep, "res": res}).Debug("fetching the logs")
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

			for i, node := range nodes {
				wg.Add(1)
				go func(node Node, i int) {
					defer wg.Done()
					coveredBlockNumbers := map[int64]struct{}{}
					first := true
					for {
						var ep string
						if nextToken == nil {
							ep = fmt.Sprintf("%s/testnets/%s/nodes/%s/blocks", conf.APIURL, testnetID, node.ID)
						} else {
							ep = fmt.Sprintf("%s/testnets/%s/nodes/%s/blocks?next=%v", conf.APIURL, testnetID, node.ID,
								url.QueryEscape(nextToken.(string)))
						}
						log.WithFields(log.Fields{"ep": ep}).Debug("fetching the log chunks")
						res, err := util.JwtHTTPRequest("GET", ep, "")
						if err != nil {
							util.PrintErrorFatal(err)
						}
						log.WithFields(log.Fields{"ep": ep, "res": res}).Debug("fetched the blocks")

						nextToken, blocks = handleExportBlocks(testnetID, node.ID, res, &coveredBlockNumbers, sem)
						appendBlocks(blocks, nextToken == nil, first, files[i])
						first = false

						if nextToken == nil {
							break
						}
					}
					files[i].Close()
				}(node, i)
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
