package core

import (
	_"github.com/ikawaha/kagome/tokenizer"
	//libSvm "github.com/ewalker544/libsvm-go"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	_"unicode/utf8"
	"strings"
	"unicode/utf8"
	)

type Sentense struct {
	LineStart int
	ColStart int
	LineEnd int
	ColEnd int
	Sentense []byte
}


func LoadTextFromPath(p string) (text []byte) {
	pathExact := filepath.ToSlash(p)
	text, err := ioutil.ReadFile(pathExact)
	if err != nil {
		fmt.Println("can't open the file")
		os.Exit(1)
	}
	return text
}

func Split(data []byte, o int, l int, c int) []Node {
	s := splitter{}
	s.init(data, o, l, c)
	fmt.Println(s.totalLines, s.size)
	nodes := s.split()
	cnt := 0
	res := *nodes
	for _, n := range res {
		if n.Pos.Line == 0 && n.Pos.Column == 0 {
			break
		}
		cnt++
	}
	r := res[:cnt]

	return r
}

func CountCharacters(s *Node) int {
	c := utf8.RuneCount([]byte(s.Data))
	return c
}

func SkipSentences(ns []*Node, cs string) []*Node {
	var nns []*Node
	var inSkip bool

	inSkip = false

	for _, n := range ns {
		if strings.HasPrefix(n.Data, cs) {
			if !inSkip {
				if strings.Contains(n.Data, "STARTSKIP") {
					inSkip = true
				} else {
					nns = append(nns, n)
				}

			} else {
				if strings.Contains(n.Data, "ENDSKIP") {
					inSkip = false
				}
			}
		} else {
			nns = append(nns, n)
		}
	}
	return nns
}