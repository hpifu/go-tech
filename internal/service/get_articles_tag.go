package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GETArticlesTagReq struct {
	Tag    string `json:"tag" uri:"tag"`
	Offset int    `json:"offset" form:"offset"`
	Limit  int    `json:"limit" form:"limit"`
}

type GETArticlesTagRes []*Article

func (s *Service) GETArticlesTag(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &GETArticlesTagReq{Limit: 20}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if req.Limit > 50 {
		req.Limit = 50
	}

	articles, err := s.db.SelectArticlesByTag(req.Tag, req.Offset, req.Limit)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select articles failed. err: [%v]", err)
	}

	if articles == nil {
		return req, nil, http.StatusNoContent, nil
	}

	var ids []int
	for _, article := range articles {
		ids = append(ids, article.AuthorID)
	}

	accountMap, err := s.GetAccounts(rid, ids)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get accounts failed. err: [%v]", err)
	}

	likeviewMap, err := s.db.SelectLikeviewsByArticles(articles)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get likeviews failed. err: [%v]", err)
	}

	var res GETArticlesTagRes
	for _, article := range articles {
		var avatar string
		author := "unknown"
		if a, ok := accountMap[article.AuthorID]; ok {
			avatar = a.Avatar
			author = strings.Split(a.Email, "@")[0]
		}
		like := 0
		view := 0
		if lv, ok := likeviewMap[article.ID]; ok {
			like = lv.Like
			view = lv.View
		}

		res = append(res, &Article{
			ID:       article.ID,
			AuthorID: article.AuthorID,
			Author:   author,
			Title:    article.Title,
			Tags:     strings.Split(article.Tags, ","),
			Brief:    article.Brief,
			CTime:    article.CTime.Format(time.RFC3339),
			UTime:    article.UTime.Format(time.RFC3339),
			Avatar:   avatar,
			Like:     like,
			View:     view,
		})
	}

	return req, res, http.StatusOK, nil
}
