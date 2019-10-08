package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticlesReq struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

func (s *Service) GETArticles(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ArticlesReq{Limit: 20}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if req.Limit > 50 {
		req.Limit = 50
	}

	articles, err := s.db.SelectArticles(req.Offset, req.Limit)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select ancients failed. err: [%v]", err)
	}

	if articles == nil {
		return req, nil, http.StatusNoContent, nil
	}

	return req, articles, http.StatusOK, nil
}