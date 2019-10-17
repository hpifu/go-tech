package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticlesReq struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

func (s *Service) GETArticles(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ArticlesReq{Limit: 20}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if req.Limit > 50 {
		req.Limit = 50
	}

	articles, err := s.db.SelectArticles(req.Offset, req.Limit)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select articles failed. err: [%v]", err)
	}

	if articles == nil {
		return req, nil, http.StatusNoContent, nil
	}

	var as []*Article
	for _, article := range articles {
		as = append(as, &Article{
			ID:       article.ID,
			AuthorID: article.AuthorID,
			Author:   article.Author,
			Title:    article.Title,
			Tags:     strings.Split(article.Tags, ","),
			Content:  article.Content,
			CTime:    article.CTime.Format(time.RFC3339),
			UTime:    article.UTime.Format(time.RFC3339),
		})
	}

	return req, as, http.StatusOK, nil
}
