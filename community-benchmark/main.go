package main

import (
	"community-benchmark/cli"
	"log"
	"os"
)

func main() {
	app := cli.CreateApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		return
	}
}
