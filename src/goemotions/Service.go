package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"runtime"
	"strconv"
)

const (
	SERVICE_POSITIVE_PATH = "positive.dat"
	SERVICE_NEGATIVE_PATH = "negative.dat"
)

const (
	SERVICE_COMMAND_CLOSE = iota
	SERVICE_COMMAND_RELOAD
)

type Service struct {
	positive     *WordDict
	negative     *WordDict
	tcpPort      int
	commandInput chan int
}

func NewServiceOnTcp(tcpPort int) *Service {
	self := new(Service)
	self.commandInput = make(chan int, 1)
	self.tcpPort = tcpPort
	self.ReloadImpl()
	rpc.Register(NewServiceRPC(self))
	return self
}

func GetServiceSocket(tcpPort int) (string, string) {
	return "tcp", ":" + strconv.Itoa(tcpPort)
}

func (self *Service) Loop() {
	listener, err := net.Listen(GetServiceSocket(self.tcpPort))
	if err != nil {
		panic(err)
	}
	connInput := make(chan net.Conn, 2)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Service cannot accept new connections: %v", err)
				self.Close()
			}
			connInput <- conn
		}
	}()
	for {
		select {
		case conn := <-connInput:
			rpc.ServeConn(conn)
			runtime.GC()
		case command := <-self.commandInput:
			if command == SERVICE_COMMAND_CLOSE {
				return
			} else if command == SERVICE_COMMAND_RELOAD {
				self.ReloadImpl()
			}
		}
	}
}

func (self *Service) ReloadImpl() {
	log.Printf("Loading dictionaries '%s' and '%s'\n", SERVICE_POSITIVE_PATH, SERVICE_NEGATIVE_PATH)
	self.positive = NewWordDict(SERVICE_POSITIVE_PATH)
	self.negative = NewWordDict(SERVICE_NEGATIVE_PATH)
}

func (self *Service) Close() {
	self.commandInput <- SERVICE_COMMAND_CLOSE
}

func (self *Service) Reload() {
	self.commandInput <- SERVICE_COMMAND_RELOAD
}

func (self *Service) RunScan(inputPath string) *EmotionalResult {
	log.Printf("Running scan for file '%s'", inputPath)
	var scanner EmotionScanner
	scanner.positive = self.positive
	scanner.negative = self.negative
	if scanner.positive == nil {
		panic(errors.New(fmt.Sprintf("Positive dict is nil")))
	}
	if scanner.negative == nil {
		panic(errors.New(fmt.Sprintf("Negative dict is nil")))
	}
	file, err := os.Open(inputPath)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Fatal error: %s", err.Error())))
	}
	defer file.Close()

	return scanner.Scan(file)
}
