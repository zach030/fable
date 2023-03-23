package main

import (
	"log"

	"context"

	"github.com/zach030/fable/cmd"
)

const (
	testInput1 = "/Users/zach/code/project/fable/testdata/milvus"
	testInput2 = "/Users/zach/code/project/fable/testdata/milvus2"
)

func main() {
	fable := cmd.NewFable()
	ctx := context.Background()
	err := fable.Ingest(ctx, testInput1, "")
	if err != nil {
		log.Fatal(err)
	}
	err = fable.Ingest(ctx, testInput2, "")
	if err != nil {
		log.Fatal(err)
	}
	err = fable.Search(ctx, "如何选择合适的索引类型和参数")
	if err != nil {
		log.Fatal(err)
	}
}
