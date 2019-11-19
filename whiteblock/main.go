package main

import (
	"github.com/whiteblock/cli/whiteblock/cmd"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute()
}
