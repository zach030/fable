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
	aliyunOss, _ = NewAliyunOss(endpoint, ak, sk)
	err := aliyunOss.Put(bucket, "aa", []byte(("ai-framework")))
	assert.Nil(t, err)
}

func TestAliyun(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		aliyunOss, _ = NewAliyunOss(endpoint, ak, sk)
		_, err := aliyunOss.Get(bucket, "aa")
		t.Log(err)
	})
}
