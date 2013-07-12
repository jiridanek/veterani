package main

import (
	"flag"
	"fmt"
	"veterani2013/input"
)

func main() {
	flag.Parse()
	for _, fname := range flag.Args() {
		rs := input.ReadCsos(fname)
		for _, j := range rs {
			fmt.Println(j)
		}
		fmt.Println()
	}
}
