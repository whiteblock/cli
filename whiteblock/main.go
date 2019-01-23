package main

import (
    "log"
	"./cmd"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute()
}
