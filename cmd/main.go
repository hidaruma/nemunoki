package main

import (
	"github.com/hidaruma/nemunoki/core"
	"fmt"
	"flag"
)

var filename string


func main() {
	flag.StringVar(&filename, "file", "", "filepath")

	flag.Parse()

	text := core.LoadTextFromPath(filename)

	ns := core.Split(text, 0,1,0)
	for i, n := range ns {
		fmt.Printf("%d-%d:%d %s\n", i,n.Pos.Line, n.Pos.Column, n.Data)
	}


}