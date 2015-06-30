package main

import (
	"bufio"
	"io"
	"unicode"
)

const (
	WORD_BUF_SIZE = 1024
)

type WordsReader struct {
	in *bufio.Reader
}

func NewWordsReader(in io.Reader) *WordsReader {
	ret := new(WordsReader)
	ret.in = bufio.NewReader(in)
	return ret
}

func (self *WordsReader) Read() ([]byte, error) {
	word := make([]byte, 0, WORD_BUF_SIZE)
	var didStart bool
	var err error
	var ch rune
	for {
		ch, _, err = self.in.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsLetter(ch) {
			didStart = true
			word = array_append_rune(word, ch)
		} else if didStart {
			break
		}
	}
	if err == nil && len(word) == 0 {
		err = io.EOF
	}
	return word, err
}
