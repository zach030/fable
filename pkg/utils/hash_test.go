package utils

import (
	"testing"
)

const (
	testInput1 = "/Users/zach/code/project/fable/testdata/milvus"
	testInput2 = "/Users/zach/code/project/fable/testdata/milvus2"
)

func TestReadAndCalcHash(t *testing.T) {
	_, hash1 := ReadAndCalcHash(testInput1)
	_, hash2 := ReadAndCalcHash(testInput2)
	t.Log(hash1)
	t.Log(hash2)
}
