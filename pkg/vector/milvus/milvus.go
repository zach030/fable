package milvus

import (
	"context"

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

func (m *MilvusClient) Insert(ctx context.Context, collection string, value []string, vector [][]float32) error {
	if len(value) != len(vector) {
		return model.ErrInconsistentLength
	}
	for i := range value {
		val, vec := value[i], vector[i]
		contentColumn := entity.NewColumnVarChar(defaultValueColumn, []string{val})
		vecs := make([][]float32, 0)
		vectorColumn := entity.NewColumnFloatVector(defaultVectorColumn, model.OpenAIEmbeddingDimensions, append(vecs, vec))
		if _, err := m.Client.Insert(ctx, collection, "", contentColumn, vectorColumn); err != nil {
			return err
		}
	}
	return nil
}

func (m *MilvusClient) Search(ctx context.Context, collection, field string, vector []float32) ([]model.SearchResult, error) {
	if err := m.LoadCollection(ctx, collection, false); err != nil {
		return nil, err
	}
	vec := entity.FloatVector(vector[:])
	sp, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		return nil, err
	}
	res, err := m.Client.Search(ctx, collection, []string{}, "", []string{field}, []entity.Vector{vec}, defaultVectorColumn, entity.L2, 10, sp)
	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
	}
	ret := make([]model.SearchResult, 0, len(res))
	for _, result := range res {
		var contentColumn *entity.ColumnVarChar
		for _, f := range result.Fields {
			if f.Name() == field {
				c, ok := f.(*entity.ColumnVarChar)
				if ok {
					contentColumn = c
				}
			}
		}
		for i := 0; i < result.ResultCount; i++ {
			content, err := contentColumn.ValueByIdx(i)
			if err != nil {
				return nil, err
			}
			ret = append(ret, model.SearchResult{
				Payload: content,
				Score:   result.Scores[i],
			})
		}
	}
	return ret, nil
}
