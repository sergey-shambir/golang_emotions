package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"os"
	"path/filepath"
)

const APP_USAGE = `Usage:
  %s <command> [<filename>]
  %s -service -port=8080
Commands:
  scan <filename> - scans filename for emotios
  close - stops background service
  reload - reloads settings`

func ShowApplicationUsage() {
	appName := os.Args[0]
	fmt.Fprintf(os.Stderr, APP_USAGE, appName, appName)
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()
}

type Application struct {
	isService bool
	tcpPort   int
}

func (self *Application) Init() {
	flag.BoolVar(&self.isService, "service", false, "run as service on TCP socket")
	flag.IntVar(&self.tcpPort, "port", 8080, "TCP port for service socket")
	flag.Usage = ShowApplicationUsage
	flag.Parse()
}

func (self *Application) Exec() (result bool) {
	defer func() {
		if err := recover(); err != nil {
			PrintBacktrace(err)
			result = false
		}
	}()

	if self.isService {
		service := NewServiceOnTcp(self.tcpPort)
		service.Loop()
		return true
	}
	return self.ExecClient()
}

func (self *Application) ExecClient() (result bool) {
	if flag.NArg() == 0 {
		ShowApplicationUsage()
		return false
	}

	rpcClient, err := rpc.Dial(GetServiceSocket(self.tcpPort))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open RPC connection: %v\n", err)
		return false
	}

	switch flag.Arg(0) {
	case "close":
		ClientCloseService(rpcClient)
	case "reload":
		ClientReloadService(rpcClient)
	case "scan":
		if flag.NArg() != 2 {
			ShowApplicationUsage()
			return false
		}
		inputPath, err := filepath.Abs(flag.Arg(1))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot get full path from '%s'\n", flag.Arg(1))
			return false
		}
		result := ClientScanEmotions(rpcClient, inputPath)
		fmt.Fprintf(os.Stdout, "Text has %d words, %.2f%% positive and %.2f%% negative\n", result.WordCount, result.PercentPositive, result.PercentNegative)
	default:
		ShowApplicationUsage()
		return false
	}
	return true
}

func main() {
	app := new(Application)
	app.Init()
	exitCode := 1
	if app.Exec() {
		exitCode = 0
	}
	os.Exit(exitCode)
}
