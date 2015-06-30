package main

import (
	"fmt"
	"io"
	"os"
)

type EmotionScanner struct {
	positive      *WordDict
	negative      *WordDict
	in            *WordsReader
	positiveMarks int
	negativeMarks int
}

func (self *EmotionScanner) Scan(in io.Reader) bool {
	var wordCount int
	self.in = NewWordsReader(in)
	word, err := self.in.Read()
	for err == nil {
		wordCount++
		if self.positive.MatchWord(word) {
			self.positiveMarks++
		}
		if self.negative.MatchWord(word) {
			self.negativeMarks++
		}
		word, err = self.in.Read()
	}
	if err != io.EOF {
		fmt.Fprintf(os.Stderr, "I/O error: %s\n", err.Error())
		return false
	}
	if self.negativeMarks+self.positiveMarks == 0 {
		fmt.Fprintf(os.Stdout, "Text has %d words without any emotions.\n", wordCount)
	} else {
		percentPos := 100.0 * float32(self.positiveMarks) / float32(wordCount)
		percentNeg := 100.0 * float32(self.negativeMarks) / float32(wordCount)
		fmt.Fprintf(os.Stdout, "Text has %d words, %.2f%% positive and %.2f%% negative\n", wordCount, percentPos, percentNeg)
	}
	return true
}
