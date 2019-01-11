package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	blockchain string
	server     string
	nodes      int
	image      string
	cpus       int
	memory     int
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

	viper.AddConfigPath("./.config/whiteblock")
	viper.AddConfigPath(home + "/.config/whiteblock")
	viper.SetConfigName("config")

	b, err := ioutil.ReadFile(home + "/.config/whiteblock/config.json")
	if err != nil {
		//fmt.Print(err)
	}
	var config Config
	json.Unmarshal(b, &config)

	err = viper.ReadInConfig()
	if err != nil {
		//fmt.Println("Default parameters are not set. Please continue and enter build fields.")
		return
	}

	if len(config.Servers) > 0 {
		server = strconv.Itoa(config.Servers[0])
	}
	
	blockchain = config.Blockchain
	nodes = config.Nodes
	image = config.Image

	viper.WatchConfig()
	viper.AutomaticEnv()
}
