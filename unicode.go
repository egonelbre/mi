// +build ignore

package main

import (
	"flag"
	"fmt"
)

var (
	start = flag.Int("s", 0x20, "start at unicode")
	width = flag.Int("w", 0x20, "number of characters per line")
	count = flag.Int("c", 0x200, "count of characters")
)

func main() {
	flag.Parse()

	r := rune(0)
	e := rune(*start)
	w := rune(*width)

	fmt.Println()
	fmt.Printf("%04x ", e)
	for {
		r = e
		e = r + rune(*count)
		for r < e {
			fmt.Print(string(r), " ")
			r++
			if r%w == 0 {
				fmt.Println()
				fmt.Printf("%04x ", r)
			}
		}

		v := ""
		fmt.Scanln(&v)
		if v != "" {
			return
		}
	}
}
