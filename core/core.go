package core

import (
	_"github.com/ikawaha/kagome/tokenizer"
	//libSvm "github.com/ewalker544/libsvm-go"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	_"unicode/utf8"
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

func RuleEvenSymbol(s Sentense)  {
	fwDash := '—'
	tpLeader := '…'
	symbolPositions := []Position{}
}
func RuleMaxSymbol(s Sentense, symbol rune, max int) (over bool, count int) {


	return over, count
}

func RuleLongKanjiChain(s Sentense, max int) (over bool) {

	over = false

	return over
}



