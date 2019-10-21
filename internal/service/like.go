package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LikeReq struct {
	ID int `json:"id,omitempty" uri:"id"`
}

type LikeRes struct{}

func (s *Service) Like(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &LikeReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	err := s.db.Like(req.ID)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, fmt.Errorf("mysql count tag failed. err: [%v]", err)
	}

	return nil, nil, http.StatusOK, err
}
