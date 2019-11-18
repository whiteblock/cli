package main

import (
	"github.com/whiteblock/cli/whiteblock/cmd"
	"log"
)

function main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute()
}
