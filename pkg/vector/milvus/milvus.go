package milvus

import (
	"context"
	"encoding/json"

	"github.com/zach030/fable/pkg/vector/model"

	"log"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusClient struct {
	client.Client
}

func NewMilvusClient(addr, user, pwd string) (*MilvusClient, error) {
	ctx, cf := context.WithTimeout(context.Background(), 3*time.Second)
	defer cf()
	milvusCli, err := client.NewDefaultGrpcClientWithURI(ctx, addr, user, pwd)
	if err != nil {
		return nil, err
	}
	return &MilvusClient{milvusCli}, nil
}

func (m *MilvusClient) Insert(ctx context.Context, req *model.InsertRequest) error {
	collection := req.Collection
	contentKey := req.ContentKey
	metadata := model.ChunkMetadata{Author: "zach"}
	for i, content := range req.Content {
		embed := req.Embeddings[i]
		var (
			contentKeyColumn = newContentKey(contentKey)
			contentColumn    = newContent(content)
			vectorColumn     = newContentVector(embed)
			metadataColumn   = newMetadata(metadata)
		)
		if _, err := m.Client.Insert(ctx, collection, "", contentKeyColumn, contentColumn, vectorColumn, metadataColumn); err != nil {
			return err
		}
	}
	return nil
}

func (m *MilvusClient) Search(ctx context.Context, req *model.SearchRequest) ([]model.SearchResult, error) {
	collection := req.Collection
	vector := req.Vector
	if err := m.LoadCollection(ctx, collection, false); err != nil {
		return nil, err
	}
	vec := entity.FloatVector(vector[:])
	sp, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		return nil, err
	}
	res, err := m.Client.Search(ctx, collection, []string{}, "", []string{filedContentKey, filedContent, filedMetadata}, []entity.Vector{vec}, filedContentVector, entity.L2, 10, sp)
	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
	}
	ret := make([]model.SearchResult, 0, len(res))
	for _, result := range res {
		var (
			contentColumn    *entity.ColumnVarChar
			contentKeyColumn *entity.ColumnVarChar
			metadataColumn   *entity.ColumnVarChar
		)
		for _, f := range result.Fields {
			if f.Name() == filedContent {
				c, ok := f.(*entity.ColumnVarChar)
				if ok {
					contentColumn = c
				}
			}
			if f.Name() == filedContentKey {
				c, ok := f.(*entity.ColumnVarChar)
				if ok {
					contentKeyColumn = c
				}
			}
			if f.Name() == filedMetadata {
				c, ok := f.(*entity.ColumnVarChar)
				if ok {
					metadataColumn = c
				}
			}
		}
		for i := 0; i < result.ResultCount; i++ {
			content, err := contentColumn.ValueByIdx(i)
			if err != nil {
				return nil, err
			}
			key, err := contentKeyColumn.ValueByIdx(i)
			if err != nil {
				return nil, err
			}
			metadataBuf, err := metadataColumn.ValueByIdx(i)
			if err != nil {
				return nil, err
			}
			metadata := model.ChunkMetadata{}
			_ = json.Unmarshal([]byte(metadataBuf), &metadata)
			ret = append(ret, model.SearchResult{
				Payload:  content,
				Key:      key,
				Metadata: metadata,
				Score:    result.Scores[i],
			})
		}
	}
	return ret, nil
}
