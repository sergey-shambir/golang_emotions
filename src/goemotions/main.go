package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func ShowApplicationUsage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [-s] [-in=<path>]\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()
}

type Application struct {
	input        string
	positivePath string
	negativePath string
	scanner      EmotionScanner
}

var g_app *Application

func (self *Application) Init() {
	flag.StringVar(&self.input, "in", "", "use self file instead of stdin input")
	flag.StringVar(&self.positivePath, "dict-positive", "positive.dat", "path to dictionary with positive emotion rules")
	flag.StringVar(&self.negativePath, "dict-negative", "negative.dat", "path to dictionary with negative emotion rules")
	flag.Usage = ShowApplicationUsage
	flag.Parse()
}

func (self *Application) Exec() {
	defer func() {
		if err := recover(); err != nil {
			PrintBacktrace(err)
		}
	}()
	self.scanner.positive = NewWordDict(self.positivePath)
	self.scanner.negative = NewWordDict(self.negativePath)

	var in io.Reader
	if len(self.input) != 0 {
		file, err := os.Open(self.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			os.Exit(1)
		}
		defer file.Close()
		in = file
	} else {
		in = os.Stdin
	}

	if self.scanner.Scan(in) {
		os.Exit(0)
	}
	os.Exit(1)
}

func main() {
	g_app = new(Application)
	g_app.Init()
	g_app.Exec()
}
