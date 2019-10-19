package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-tech/internal/mysql"
	"net/http"
)

type GETTagCloudReq struct{}

type GETTagCloudRes []*mysql.TagCountPair

func (s *Service) GETTagCloud(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	res, err := s.db.CountTag()
	if err != nil {
		return nil, nil, http.StatusInternalServerError, fmt.Errorf("mysql count tag failed. err: [%v]", err)
	}
	return nil, res, http.StatusOK, err
}
