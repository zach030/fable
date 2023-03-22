package cmd

import (
	"context"
	"fmt"

	"github.com/zach030/fable/pkg/utils"

	"github.com/zach030/fable/pkg/llm"
	"github.com/zach030/fable/pkg/llm/openai"
	"github.com/zach030/fable/pkg/storage"
	"github.com/zach030/fable/pkg/storage/oss"
	"github.com/zach030/fable/pkg/vector"
	"github.com/zach030/fable/pkg/vector/milvus"
)

func init() {
	cfg = InitCfg()
}

var (
	cfg *FableConfig
)

type Fable struct {
	llm       llm.Llm
	storageDB storage.Storage
	vectorDB  vector.Vector
}

func NewFable() *Fable {
	vectorDB, err := milvus.NewMilvusClient(cfg.Milvus.Addr, cfg.Milvus.User, cfg.Milvus.Password)
	if err != nil {
		return nil
	}
	aliyunOss, err := oss.NewAliyunOss(cfg.AliyunOSS.AliyunOssEndpoint, cfg.AliyunOSS.AliyunOssAccessKey, cfg.AliyunOSS.AliyunOssSecretKey)
	if err != nil {
		return nil
	}
	return &Fable{
		llm:       openai.NewOpenAI(cfg.OpenAI.APIKey, cfg.OpenAI.ProxyURL),
		storageDB: aliyunOss,
		vectorDB:  vectorDB,
	}
}

func (f *Fable) Ingest(ctx context.Context, path, key string) error {
	buf, hashKey := utils.ReadAndCalcHash(path)
	var ossKey = hashKey
	if key != "" {
		ossKey = key
	}
	exist, err := f.storageDB.Exist(cfg.AliyunOSS.Bucket, ossKey)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	if err = f.storageDB.Put(cfg.AliyunOSS.Bucket, ossKey, buf); err != nil {
		return err
	}
	embeddings, err := f.llm.Embedding(ctx, []string{string(buf)})
	if err != nil {
		return err
	}
	vecs := make([][]float32, 0)
	return f.vectorDB.Insert(ctx, cfg.Milvus.Collection, []string{string(buf)}, append(vecs, embeddings))
}

func (f *Fable) Search(ctx context.Context, input string) error {
	embedding, err := f.llm.Embedding(ctx, []string{input})
	if err != nil {
		return err
	}
	result, err := f.vectorDB.Search(ctx, cfg.Milvus.Collection, "", embedding)
	if err != nil {
		return err
	}
	ossKey := make([]string, 0, len(result))
	for _, searchResult := range result {
		ossKey = append(ossKey, searchResult.Payload)
	}
	ossVal := make([]string, 0, len(ossKey))
	for _, key := range ossKey {
		buf, err := f.storageDB.Get(cfg.AliyunOSS.Bucket, key)
		if err != nil {
			continue
		}
		ossVal = append(ossVal, string(buf))
	}
	var res string
	prompt, err := f.llm.PreparePrompt(ctx, ossVal, []string{res})
	if err != nil {
		return err
	}
	answer, err := f.llm.Completion(ctx, prompt)
	if err != nil {
		return err
	}
	fmt.Println(answer)
	return nil
}
