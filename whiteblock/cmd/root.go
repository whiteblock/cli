package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	blockchain string
	server     string
)

type Iface struct {
	Ip      string `json:"ip"`
	Gateway string `json:"gateway"`
	Subnet  int    `json:"subnet"`
}

type Switch struct {
	Addr  string `json:"addr"`
	Iface string `json:"iface"`
	Brand int    `json:"brand"`
	Id    int    `json:"id"`
}

type Server struct {
	Addr     string   `json:"addr"`  //IP to access the server
	Iaddr    Iface    `json:"iaddr"` //Internal IP of the server for NIC attached to the vyos
	Nodes    int      `json:"nodes"`
	Max      int      `json:"max"`
	Id       int      `json:"id"`
	ServerID int      `json:"serverID"`
	Iface    string   `json:"iface"`
	Switches []Switch `json:"switches"`
	Ips      []string `json:"ips"`
}

var RootCmd = &cobra.Command{
	Use:   "whiteblock",
	Short: "Create and test blockchains",
	Long: `This application will deploy a blockchain, create nodes, and allow those nodes to interact in the network. Documentation, usages, and exmaples can be found at www.whiteblock.io/docs/cli.
	`,
}

func writeConfigFile(configFile string) {
	cwd := os.Getenv("HOME")
	err := os.MkdirAll(cwd+"/cli/whiteblock/config", 0755)
	if err != nil {
		log.Fatalf("could not create directory: %s", err)
	}

	file, err := os.Create(cwd + "/cli/whiteblock/config/config.json")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(configFile)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath("$HOME/cli/whiteblock/config")
	viper.AddConfigPath("./cli/whiteblock/config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(home + ".config/whiteblock")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		println("No config file could be found. Please follow the steps to create one.")
		fmt.Println(err)
		configArr := make([]string, 0)
		configOpt := [1]string{"blockchain"}

		idList := make([]string, 0)

		scanner := bufio.NewScanner(os.Stdin)
		tmp := 0
		for {
			if tmp == len(configOpt) {
				break
			}
			for i := 0; i < len(configOpt); i++ {
				fmt.Print(configOpt[i] + ": ")
				scanner.Scan()

				text := scanner.Text()
				if len(text) == 0 {
					println("invalid")
					break
				}
				configArr = append(configArr, text)
				tmp = i + 1
			}
		}

		getServerAddr := "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		command := "get_servers"
		results := []byte(wsEmitListen(getServerAddr, command, ""))
		var result map[string]Server
		err := json.Unmarshal(results, &result)
		if err != nil {
			panic(err)
		}

		serverID := 0
		for _, v := range result {
			serverID = v.ServerID
			//move this and take out break statement if instance has multiple servers
			idList = append(idList, fmt.Sprintf("%d", serverID))
			break
		}

		server = strings.Join(idList, ",")

		blockchain := configArr[0]
		param := "{\"blockchain\":\"" + blockchain + "\",\"server\":\"" + fmt.Sprintf(server) + "\"}"
		println(param)
		writeConfigFile(param)

		viper.ReadInConfig()
	}

	blockchain = viper.GetString("blockchain")
	if !viper.IsSet("blockchain") {
		blockchain = "ethereum"
	}
	server = viper.GetString("server")
	if !viper.IsSet("server") {
		server = "1"
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	viper.AutomaticEnv()
}
