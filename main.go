package main

import (
	"github.com/nsf/termbox-go"

	"github.com/loov/mi/edit"
	"github.com/loov/mi/render"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	buffer, err := edit.BufferFromFile("main.go")
	if err != nil {
		panic(err)
	}

	render.Buffer(buffer)
	termbox.Flush()

inputloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break inputloop
			}
			switch ev.Ch {
			case 'q', 'Q':
				break inputloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		render.Buffer(buffer)
		termbox.Flush()
	}
}
