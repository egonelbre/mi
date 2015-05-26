package main

import (
	"github.com/loov/mi/cli"
	"github.com/loov/mi/edit"
)

func main() {
	buffer, err := edit.BufferFromFile("main.go")
	if err != nil {
		panic(err)
	}

	cli.Run(buffer)
}
