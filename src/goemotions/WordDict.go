package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

type WordDict struct {
	exprs       []*regexp.Regexp
	runeMapping map[rune]rune
	lineNo      int
}

func (self *WordDict) Panic(message string) {
	panic(errors.New(fmt.Sprintf("Dict load failed on line #%d: %s", self.lineNo, message)))
}

func (self *WordDict) AddReplaceRule(rule []byte) {
	if utf8.RuneCount(rule) != 2 {
		self.Panic("rule format differs from '=xX'")
	}
	from, fromSize := utf8.DecodeRune(rule)
	to, _ := utf8.DecodeRune(rule[fromSize:])
	self.runeMapping[from] = to
}

func (self *WordDict) AddRegexRule(rule string) {
	expr, err := regexp.Compile(rule)
	if err != nil {
		self.Panic(err.Error())
	}
	self.exprs = append(self.exprs, expr)
}

func (self *WordDict) MatchWord(word []byte) bool {
	var mappedWord []byte
	for i := 0; i < len(word); {
		ch, size := utf8.DecodeRune(word[i:])
		if mapped, ok := self.runeMapping[ch]; ok {
			mappedWord = array_append_rune(mappedWord, mapped)
		} else {
			mappedWord = append(mappedWord, word[i:i+size]...)
		}
		i += size
	}
	for _, expr := range self.exprs {
		if expr.Match(mappedWord) {
			return true
		}
	}
	return false
}

func (self *WordDict) Load(reader io.Reader) {
	self.runeMapping = make(map[rune]rune)
	in := bufio.NewReader(reader)
	for self.lineNo = 1; ; self.lineNo++ {
		lineBytes, isPrefix, err := in.ReadLine()
		if isPrefix {
			self.Panic("too big line")
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			self.Panic("Dict load failed: " + err.Error())
		}
		line := string(lineBytes)
		if strings.HasPrefix(line, "=") {
			self.AddReplaceRule(lineBytes[1:])
		} else if strings.HasPrefix(line, "^") {
			self.AddRegexRule(line[1:])
		} else {
			self.Panic("unknown rule format on line, use '^REGEXP' to match regex or '=xX' to replace letters")
		}
	}
}

func NewWordDict(filePath string) *WordDict {
	file, err := os.Open(filePath)
	if err != nil {
		file, err = os.Open(filepath.Join(GetExecutableDir(), filePath))
		if err != nil {
			panic(err)
		}
	}
	defer file.Close()
	dict := new(WordDict)
	dict.Load(file)
	return dict
}
