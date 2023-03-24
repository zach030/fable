package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/zach030/fable/cmd"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

const (
	limitSize = 1024 * 200
)

var (
	ErrFileExceedLimit = errors.New("file too large")
	ErrEmptyInput      = errors.New("empty input")

	engine   *gin.Engine
	fableSrv *FableService
)

func NewFableSrv(f *cmd.Fable) {
	fableSrv = &FableService{fable: f}
}

func StartServer() {
	engine = gin.Default()
	initRouter()
	engine.Run(":8888")
}

func initRouter() {
	api := engine.Group("/api")
	{
		api.POST("/upload", uploadHandler)
		api.GET("/search", searchHandler)
	}
}

func uploadHandler(c *gin.Context) {
	key, _ := c.GetPostForm("name")
	fh, err := c.FormFile("file")
	if err != nil {
		c.Error(err)
		return
	}
	if fh.Size > limitSize {
		c.Error(ErrFileExceedLimit)
		return
	}
	name := fmt.Sprintf("./tmp/%s", uuid.New().String())
	if err := c.SaveUploadedFile(fh, name); err != nil {
		c.Error(err)
		return
	}
	defer os.Remove(name)
	key, err = fableSrv.Upload(c, name, key)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": key,
	})
	return
}

func searchHandler(c *gin.Context) {
	q := c.Query("question")
	key := c.Query("key")
	if q == "" || key == "" {
		c.Error(ErrEmptyInput)
		return
	}
	answer, err := fableSrv.Search(c, q, key)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": answer,
	})
	return
}
