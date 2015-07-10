package main

import (
	"errors"
	"fmt"
	"io"
)

type EmotionalResult struct {
	PercentPositive float32
	PercentNegative float32
	WordCount       int
}

type EmotionScanner struct {
	positive      *WordDict
	negative      *WordDict
	in            *WordsReader
	positiveMarks int
	negativeMarks int
}

func (self *EmotionScanner) Scan(in io.Reader) *EmotionalResult {
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
		panic(errors.New(fmt.Sprintf("I/O error: %s\n", err.Error())))
	}
	return &EmotionalResult{
		WordCount:       wordCount,
		PercentPositive: 100.0 * float32(self.positiveMarks) / float32(wordCount),
		PercentNegative: 100.0 * float32(self.negativeMarks) / float32(wordCount),
	}
}
