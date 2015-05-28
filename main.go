package main

import (
	"os"
	"log"
	"flag"

	"github.com/loov/mi/cli"
	"github.com/loov/mi/edit"
)

func main() {
	flag.Parse()
	
	file, err := os.OpenFile("debug", 0666, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	
	log.SetOutput(file)

	buffer, err := edit.BufferFromFile("main.go")
	if err != nil {
		panic(err)
	}

	cli.Run(buffer)
}
