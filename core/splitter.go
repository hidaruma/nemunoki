package core

import (
	"fmt"
	"unicode/utf8"
	"bytes"
	"os"
)

type Position struct {
	Offset int
	Line int
	Column int
}

func (p *Position) IsValid() bool {
	return p.Line > 0
}

func (p *Position) String() string {
	var s string
	if p.IsValid() {
		s += fmt.Sprintf("%d", p.Line)

		s += fmt.Sprintf(":")
		s += fmt.Sprintf("%d", p.Column)
	}
	return s
}

type Pos uint
const NoPos Pos = 0
func (p Pos) IsValid() bool {
	return p != NoPos
}

type Node struct {
	Pos Position
	Data string
}

func (n *Node) set(o int,l int, c int) {
	n.Pos.Offset = o
	n.Pos.Line = l
	n.Pos.Column = c
	n.Data = ""
}

type splitter struct {

	src *bytes.Reader
	size int
	totalLines int

	ch rune
	offset int
	rdOffset int

	lineOffset int
	columnOffset int
	err errorHandler
}

type errorHandler interface {
	error(int, string)
}

func (s *splitter) error(o int, msg string) {
	fmt.Printf("%d: %s", o, msg)
	os.Exit(1)
}

func (s *splitter) currentOffset() int {
	return s.offset
}

func (s *splitter) currentLine() int {
	return s.lineOffset
}

func (s *splitter) currentColumn() int {
	return s.columnOffset
}

const bom = 0xFEFF

func (s *splitter) next() (size int){

	if s.rdOffset < s.size {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.lineOffset++
			s.columnOffset = 0
		}
		r, size, err := s.src.ReadRune()
		if err != nil {
			s.error(s.offset, "ReadRune Error")
		}
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NUL")
		case r >= utf8.RuneSelf:
			if r == utf8.RuneError && size == 1 {
				s.error(s.offset, "illegal UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "illegal byte order mark")
			}
		}
		s.rdOffset += size
		s.columnOffset++
		s.ch = r
	} else {
		//s.rdOffset = s.size
		if s.ch == '\n' {
			s.lineOffset++
			s.columnOffset = 1
		}
		s.ch = -1
	}
	return size
}

func (s *splitter) init(src []byte) {
	s.src = bytes.NewReader(src)
	s.size = len(src)
	s.totalLines = bytes.Count(src, []byte{'\n'})

	s.ch = 0
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 1
	s.columnOffset = 1
}


func (s *splitter) split() (*[]Node) {
	nodes := make([]Node, s.totalLines, s.size)
	var sentence []rune

	max := s.size
	eos := true //ch is End of Sentence or not.
	for i := 0; i < max; {
		var o, l, c int
		ch := s.ch
		s.next()
		if eos {
			o = s.currentOffset()
			l = s.currentLine()
			c = s.currentColumn()
			nodes[i].set(o, l, c)
			eos =false
		}

		switch {
		case ch == -1:
			nodes[i].Data = string(sentence)
			sentence = []rune{}
			i = max
			break
		case ch == '\n':
			if s.ch == '\n' {

				nodes[i].Data = string(sentence)
				sentence = []rune{}
				eos = true
				i++

			} else {
				continue
			}
		case isDelimNoSpace(ch):
			sentence = append(sentence,ch)
			nodes[i].Data = string(sentence)
			sentence = []rune{}
			eos = true
			i++

		case isDelimNeedSpace(ch):
			if s.ch == ' ' || s.ch == '\n' {
				sentence = append(sentence, ch)
				nodes[i].Data = string(sentence)
				sentence = []rune{}
				eos = true
				i++
			}
			fallthrough

		case isDelimNeedWideSpace(ch):
			if s.ch == '　' || s.ch == '\n' || isFollower(s.ch) {
				sentence = append(sentence, ch)
				nodes[i].Data = string(sentence)
				sentence = []rune{}
				eos = true
				i++
				continue
			}
			fallthrough

		case isFollower(ch):

/*			if s.ch == '\n' || unicode.IsSpace(s.ch) {
				sentence = append(sentence, ch)
				nodes[i].Data = string(sentence)
				sentence = []rune{}
				eos = true
				i++
			}
			fallthrough
*/
		default:
			if ch != 0 || ch != '\n' || ch != '　' {
				sentence = append(sentence, ch)
			}

		}


	}

	return &nodes
}


func runeCheck(ch rune, rs []rune) bool {
	for _, d := range rs {
		if ch == d {
			return true
		}
	}
	return false

}

func isDelimNoSpace(r rune) bool {
	return runeCheck(r, splitDelimNoSpace)
}

func isDelimNeedSpace(r rune) bool {
	return runeCheck(r, splitDelimNeedSpace)
}

func isDelimNeedWideSpace(r rune) bool {
	return runeCheck(r, splitDelimNeedWideSpace)
}


func isFollower(r rune) bool {
	return runeCheck(r, splitFollower)
}

var splitDelimNoSpace = []rune{
	'。',
	'．',
}

var splitDelimNeedSpace = []rune{
	'.',
	'!',
	'?',
}

var splitDelimNeedWideSpace = []rune{
	'！',
	'？',
}

var splitFollower = []rune{
	']',
	'}',
	')',
	'」',
	'』',
	'〕',
	'】',
	'’',
	'》',
	'”',
	'〟',
	'〉',

}

