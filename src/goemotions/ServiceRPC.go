package main

import (
	"net/rpc"
)

type ServiceRPC struct {
	service *Service
}

func NewServiceRPC(service *Service) *ServiceRPC {
	self := new(ServiceRPC)
	self.service = service
	return self
}

type ServiceArgsEmpty struct {
}

type ServiceReplyEmpty struct {
}

type ServiceArgsScan struct {
	Path string
}

// RPC for close server

func (self *ServiceRPC) CloseService(args *ServiceArgsEmpty, reply *ServiceReplyEmpty) error {
	self.service.Close()
	return nil
}

func ClientCloseService(client *rpc.Client) {
	args := &ServiceArgsEmpty{}
	var reply ServiceReplyEmpty
	err := client.Call("ServiceRPC.CloseService", args, &reply)
	if err != nil {
		panic(err)
	}
}

// RPC for reload

func (self *ServiceRPC) ReloadService(args *ServiceArgsEmpty, reply *ServiceReplyEmpty) error {
	self.service.Reload()
	return nil
}

func ClientReloadService(client *rpc.Client) {
	args := &ServiceArgsEmpty{}
	var reply ServiceReplyEmpty
	err := client.Call("ServiceRPC.ReloadService", args, &reply)
	if err != nil {
		panic(err)
	}
}

// RPC for scan

func (self *ServiceRPC) ScanEmotions(args *ServiceArgsScan, reply *EmotionalResult) error {
	*reply = *self.service.RunScan(args.Path)
	return nil
}

func ClientScanEmotions(client *rpc.Client, inputPath string) EmotionalResult {
	args := &ServiceArgsScan{Path: inputPath}
	var reply EmotionalResult
	err := client.Call("ServiceRPC.ScanEmotions", args, &reply)
	if err != nil {
		panic(err)
	}
	return reply
}

// RPC for
