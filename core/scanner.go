package core

import (
	"fmt"
	"unicode"
	"unicode/utf8"
	"bytes"
)

type Position struct {
	Offset int
	Line int
	Column int
}

func (p *Position) IsValid() bool {
	return p.Line > 0
}

func (p Position) String() string {
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
	EoP bool // End of paragraph or not
}

func (n *Node) set(s *splitter, eop bool) {
	n.Pos.Offset = s.offset
	n.Pos.Line = s.lineOffset
	n.Pos.Column = s.columnOffset
	n.EoP = eop
	n.Data = ""
}

type splitter struct {

	src *bytes.Reader
	size int

	ch rune
	offset int
	rdOffset int

	lineOffset int
	columnOffset int
	error func(int, string)
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
		s.offset = s.size
		if s.ch == '\n' {
			s.lineOffset++
			s.columnOffset = 0
		}
		s.ch = -1
	}
	return size
}

func (s *splitter) init(src []byte) {
	s.src = bytes.NewReader(src)
	s.size = len(src)
	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.columnOffset = 0
}


func (s *splitter) split() (nodes []*Node) {

	var sentense []rune
	var n Node
	n.set(s, false)
	for {
		ch := s.ch
		s.next()
		switch {
		case ch == -1:
			n.EoP = true
			nodes = append(nodes, &n)
			break

		case unicode.IsSpace(ch):
			continue
		case ch == '\n':
			if s.ch == '\n' {

				n.Data = string(sentense)
				nodes = append(nodes, &n)
				n.set(s, true)
			}
		case isDelimNoSpace(ch):
			sentense = append(sentense,ch)
			n.Data = string(sentense)
			nodes = append(nodes, &n)
			n.set(s, false)
		case isDelimNeedSpace(ch):
			if s.ch == ' ' {
				sentense = append(sentense, ch)
				n.Data = string(sentense)
				nodes = append(nodes, &n)
				n.set(s, false)
			}
		case isDelimNeedWideSpace(ch):
			if s.ch == '　' {
				sentense = append(sentense, ch)
				nodes = append(nodes, &n)
				n.set(s, false)
			}
		case isFollower(ch):
			sentense = append(sentense, ch)
			if s.ch == '\n' || s.ch == ' ' {
				nodes = append(nodes, &n)
			}
		}
		sentense = append(sentense, ch)

	}

	return nodes
}


func isDelimNoSpace(r rune) bool {
	for _, d := range splitDelimNoSpace {
		if r == d {
			return true
		}
	}
	return false
}

func isDelimNeedSpace(r rune) bool {
	for _, d := range splitDelimNeedSpace {
		if r == d {
			return true
		}
	}
	return false
}

func isDelimNeedWideSpace(r rune) bool {
	for _, d := range splitDelimNeedWideSpace {
		if r == d {
			return true
		}
	}
	return false
}


func isFollower(r rune) bool {
	for _, f := range splitFollower {
		if f == r {
			return true
		}
	}
	return false
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