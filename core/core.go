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
}

func Split(data []byte) []*Node {
	s := splitter{}
	s.init(data)
	nodes := s.split()
	return nodes
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