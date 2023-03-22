package qdrant

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/zach030/fable/pkg/vector/model"
)

type Qdrant struct {
	client *resty.Client
	apiKey string
}

func NewQdrant(baseURL, apiKey string) *Qdrant {
	c := resty.New()
	c.BaseURL = baseURL
	return &Qdrant{client: c, apiKey: apiKey}
}

func (q *Qdrant) Insert(ctx context.Context, collection string, value []string, vector [][]float32) error {
	if len(value) != len(vector) {
		return model.ErrInconsistentLength
	}
	collectionResp := &CommonResponse{}
	_, err := q.client.R().
		SetContext(ctx).
		SetHeader("api-key", q.apiKey).
		SetResult(collectionResp).
		Get(fmt.Sprintf("/collections/%s", collection))
	if err != nil {
		return err
	}
	if collectionResp.Status != "ok" {
		return errors.New("unhealthy collection")
	}
	newPointReq := &UpsertPointRequest{Points: make([]Point, 0)}
	for i := range value {
		val, vec := value[i], vector[i]
		newPointReq.Points = append(newPointReq.Points, Point{
			ID:      uuid.New().String(),
			Payload: val,
			Vector:  vec,
		})
	}
	newPointResp := &CommonResponse{}
	_, err = q.client.R().
		SetContext(ctx).
		SetHeader("api-key", q.apiKey).
		SetBody(newPointReq).
		SetResult(newPointResp).
		Put(fmt.Sprintf("/collection/%s/points?wait=true", collection))
	if err != nil {
		return err
	}
	if newPointResp.Status != "ok" {
		return errors.New("new point failed")
	}
	return nil
}

func (q *Qdrant) Search(ctx context.Context, collection, field string, vector []float32) ([]model.SearchResult, error) {
	req := &PointSearchRequest{
		Params: map[string]interface{}{
			"exact":   false,
			"hnsw_ef": 128,
		},
		Vector:      vector,
		Limit:       3,
		WithPayload: true,
		WithVector:  true,
	}
	resp := &CommonResponse{}
	_, err := q.client.R().
		SetContext(ctx).
		SetHeader("api-key", q.apiKey).
		SetBody(req).
		SetResult(resp).
		Post(fmt.Sprintf("/collections/%s/points/search", collection))
	if err != nil {
		return nil, err
	}
	if resp.Result == nil || resp.Status != "ok" {
		return nil, errors.New("search point failed")
	}
	var result []model.SearchResult
	for _, v := range resp.Result.([]interface{}) {
		res := SearchResult{}
		err = mapstructure.Decode(v, &res)
		if err != nil {
			return nil, err
		}
		result = append(result, model.SearchResult{
			Payload: fmt.Sprintf("%v", res.Payload),
			Score:   float32(res.Score),
		})
		fmt.Printf("search id=%v val=%v score=%v", res.ID, res.Payload, res.Score)
	}
	return result, nil
}
