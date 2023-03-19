package oss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	aliyunOss *AliyunOSS
	endpoint  = ""
	ak        = ""
	sk        = ""
	bucket    = ""
)

func TestAliyunOSS_Put(t *testing.T) {
	aliyunOss = NewAliyunOss(endpoint, ak, sk)
	err := aliyunOss.Put(bucket, "fable", []byte("ai-framework"))
	assert.Nil(t, err)
}
