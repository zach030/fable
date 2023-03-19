package milvus

import (
	"context"
	"fmt"

	"log"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusClient struct {
	client.Client
	collection string
}

func NewMilvusClient(addr, collection string) (*MilvusClient, error) {
	ctx, cf := context.WithTimeout(context.Background(), 3*time.Second)
	defer cf()
	milvusCli, err := client.NewGrpcClient(ctx, addr)
	if err != nil {
		return nil, err
	}
	has, err := milvusCli.HasCollection(ctx, collection)
	if err != nil {
		return nil, err
	}
	if has {
		return &MilvusClient{milvusCli, collection}, nil
	}
	log.Println("new collection fable now")
	if err = milvusCli.CreateCollection(ctx, fableDefaultCollection, 1); err != nil {
		return nil, err
	}
	return &MilvusClient{milvusCli, collection}, nil
}

func (m *MilvusClient) CreateIndex(ctx context.Context, field string) error {
	idx, err := entity.NewIndexIvfFlat( // NewIndex func
		entity.L2, // metricType
		1024,      // ConstructParams
	)
	if err != nil {
		return err
	}
	return m.Client.CreateIndex(ctx, m.collection, field, idx, false)
}

func (m *MilvusClient) Insert(ctx context.Context, value string, vector []float32) error {
	contentColumn := entity.NewColumnVarChar(defaultValueColumn, []string{value})
	vectors := make([][]float32, 0)
	vectorColumn := entity.NewColumnFloatVector(defaultVectorColumn, 768, append(vectors, vector))
	if _, err := m.Client.Insert(ctx, m.collection, "", contentColumn, vectorColumn); err != nil {
		return err
	}
	return nil
}

func (m *MilvusClient) Search(ctx context.Context, field string, vector []float32) error {
	if err := m.LoadCollection(ctx, m.collection, false); err != nil {
		log.Fatal("failed to check whether collection exists:", err.Error())
		return err
	}
	log.Println("success load collection:" + m.collection)
	vec := entity.FloatVector(vector[:])
	sp, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		return err
	}
	res, err := m.Client.Search(ctx, m.collection, []string{}, "", []string{field, "ID"}, []entity.Vector{vec}, defaultVectorColumn, entity.L2, 10, sp)
	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
	}
	for _, result := range res {
		var (
			idColumn      *entity.ColumnInt64
			contentColumn *entity.ColumnVarChar
		)
		for _, field := range result.Fields {
			if field.Name() == "ID" {
				c, ok := field.(*entity.ColumnInt64)
				if ok {
					idColumn = c
				}
			}
			if field.Name() == "content" {
				c, ok := field.(*entity.ColumnVarChar)
				if ok {
					contentColumn = c
				}
			}
		}
		for i := 0; i < result.ResultCount; i++ {
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				log.Fatal(err.Error())
			}
			content, err := contentColumn.ValueByIdx(i)
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Printf("file id: %d content: %s  scores: %f\n", id, content, result.Scores[i])
		}
	}
	return nil
}
