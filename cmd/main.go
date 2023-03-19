package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/zach030/fable/pkg/llm"
	"github.com/zach030/fable/pkg/llm/openai"
	"github.com/zach030/fable/pkg/storage"
	"github.com/zach030/fable/pkg/storage/oss"
	"github.com/zach030/fable/pkg/vector"
	"github.com/zach030/fable/pkg/vector/milvus"
)

type Fable struct {
	llm       llm.Llm
	storageDB storage.Storage
	vectorDB  vector.Vector
}

func NewFable(cfg *FableConfig) *Fable {
	vectorDB, err := milvus.NewMilvusClient(cfg.MilvusAddr, cfg.MilvusCollection)
	if err != nil {
		return nil
	}
	return &Fable{
		llm:       openai.NewOpenAI(cfg.OpenAIApiKey),
		storageDB: oss.NewAliyunOss(cfg.AliyunOssEndpoint, cfg.AliyunOssAccessKey, cfg.AliyunOssSecretKey),
		vectorDB:  vectorDB,
	}
}

func (f *Fable) IngestWithRawData(ctx context.Context, path string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return f.storageDB.Put("", "", buf)
}

func (f *Fable) IngestWithURL(ctx context.Context, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return f.storageDB.Put("", "", bytes)
}

func (f *Fable) IndexData(ctx context.Context, key string) error {
	buf, err := f.storageDB.Get("", key)
	if err != nil {
		return err
	}
	embedding, err := f.llm.Embedding(ctx, []string{string(buf)})
	if err != nil {
		return err
	}
	return f.vectorDB.Insert(ctx, string(buf), embedding)
}

func (f *Fable) Search(ctx context.Context, input string) error {
	embedding, err := f.llm.Embedding(ctx, []string{input})
	if err != nil {
		return err
	}
	f.vectorDB.Search(ctx, "", embedding)

	var res string
	prompt, err := f.llm.PreparePrompt(ctx, []string{input}, []string{res})
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
