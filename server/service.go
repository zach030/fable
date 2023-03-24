package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zach030/fable/cmd"
)

type FableService struct {
	fable *cmd.Fable
}

func (s *FableService) Upload(c *gin.Context, path, key string) (string, error) {
	name, err := s.fable.Ingest(c, path, key)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (s *FableService) Search(c *gin.Context, question, key string) ([]string, error) {
	return s.fable.Search(c, question, key)
}
