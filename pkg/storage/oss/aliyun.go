package oss

import (
	"bytes"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliyunOSS struct {
	*oss.Client
}

func NewAliyunOss(endpoint, ak, sk string) *AliyunOSS {
	client, err := oss.New(endpoint, ak, sk)
	if err != nil {
		return nil
	}
	return &AliyunOSS{client}
}

func (a *AliyunOSS) Put(bucket, key string, value []byte) error {
	bucketIns, err := a.Bucket(bucket)
	if err != nil {
		return err
	}
	return bucketIns.PutObject(key, bytes.NewReader(value))
}

func (a *AliyunOSS) Get(bucket, key string) ([]byte, error) {
	bucketIns, err := a.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	reader, err := bucketIns.GetObject(key)
	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
