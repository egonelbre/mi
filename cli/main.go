package cli

import (
	"github.com/nsf/termbox-go"

	"github.com/loov/mi/cli/render"
	"github.com/loov/mi/edit"
)

func Run(buffer *edit.Buffer) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	render.Buffer(buffer)
	render.Flush()

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
		render.Flush()
	}
}