package cmd

import (
	"context"
	"errors"
	"log"

	"github.com/zach030/fable/pkg/utils"

	"github.com/zach030/fable/pkg/llm"
	"github.com/zach030/fable/pkg/llm/openai"
	"github.com/zach030/fable/pkg/storage"
	"github.com/zach030/fable/pkg/storage/oss"
	"github.com/zach030/fable/pkg/vector"
	"github.com/zach030/fable/pkg/vector/milvus"
	"github.com/zach030/fable/pkg/vector/model"
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
		log.Println("file exist, no need to ingest")
		return nil
	}
	if err = f.storageDB.Put(cfg.AliyunOSS.Bucket, ossKey, buf); err != nil {
		return err
	}
	log.Println("success put to oss, key=", ossKey)
	split, err := utils.TextSplit(string(buf))
	if err != nil {
		return err
	}
	var (
		contents   = make([]string, 0, len(split))
		embeddings = make([][]float32, 0, len(split))
	)
	for _, s := range split {
		contents = append(contents, s.Content)
		embedding, err := f.llm.Embedding(ctx, []string{s.Content})
		if err != nil {
			return err
		}
		embeddings = append(embeddings, embedding)
	}
	return f.vectorDB.Insert(ctx, model.NewInsertRequest(cfg.Milvus.Collection, ossKey, contents, embeddings))
}

func (f *Fable) Search(ctx context.Context, input string) ([]string, error) {
	if utils.TokensNum(input) > utils.OpenAITokenLimit {
		return nil, errors.New("too many input")
	}
	embedding, err := f.llm.Embedding(ctx, []string{input})
	if err != nil {
		return nil, err
	}
	log.Println("success call llm embedding with input")
	result, err := f.vectorDB.Search(ctx, model.NewSearchRequest(cfg.Milvus.Collection, embedding))
	if err != nil {
		return nil, err
	}
	log.Println("success call vectorDB result=", result)
	indexAnswer := make([]string, 0, len(result))
	for _, searchResult := range result {
		indexAnswer = append(indexAnswer, searchResult.Payload)
	}
	prompt := f.llm.PreparePrompt(1, indexAnswer, input)
	if err != nil {
		return nil, err
	}
	answers, err := f.llm.Completion(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return answers, nil
}
