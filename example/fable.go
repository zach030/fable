package main

import (
	"github.com/zach030/fable/cmd"
	"github.com/zach030/fable/server"
)

const (
	testInput1 = "/Users/zach/code/project/fable/testdata/milvus"
	testInput2 = "/Users/zach/code/project/fable/testdata/milvus2"
)

func main() {
	server.NewFableSrv(cmd.NewFable())
	server.StartServer()
}
