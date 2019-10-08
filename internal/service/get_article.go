package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ArticleReq struct {
	ID int `uri:"id"`
}

func (s *Service) GETArticle(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ArticleReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	article, err := s.db.SelectArticleByID(req.ID)

	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select ancient failed. err: [%v]", err)
	}

	if article == nil {
		return req, nil, http.StatusNoContent, nil
	}

	return req, article, http.StatusOK, nil
}
