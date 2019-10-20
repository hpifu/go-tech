package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SearchReq struct {
	Q      string `form:"q"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

func (s *Service) Search(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &SearchReq{Limit: 20}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if req.Limit > 50 {
		req.Limit = 50
	}

	// search articles ids from es
	ids, err := s.es.SearchArticle(req.Q, req.Offset, req.Limit)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("elasticsearch search article failed. err: [%v]", err)
	}
	if ids == nil {
		return req, nil, http.StatusNoContent, nil
	}

	// select articles from mysql
	articles, err := s.db.SelectArticlesByIDs(ids)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select articles by ids failed. err: [%v]", err)
	}

	// get author info from account
	accountMap, err := s.GetAccounts(rid, ids)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get accounts failed. err: [%v]", err)
	}

	var as GETArticlesRes
	for _, article := range articles {
		var avatar string
		author := "unknown"
		if a, ok := accountMap[article.AuthorID]; ok {
			avatar = a.Avatar
			author = strings.Split(a.Email, "@")[0]
		}

		as = append(as, &Article{
			ID:       article.ID,
			AuthorID: article.AuthorID,
			Author:   author,
			Title:    article.Title,
			Tags:     strings.Split(article.Tags, ","),
			Brief:    article.Brief,
			CTime:    article.CTime.Format(time.RFC3339),
			UTime:    article.UTime.Format(time.RFC3339),
			Avatar:   avatar,
		})
	}

	return req, as, http.StatusOK, nil
}
