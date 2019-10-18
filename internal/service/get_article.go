package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GETArticleReq struct {
	ID int `uri:"id" json:"id"`
}

type GETArticleRes Article

func (s *Service) GETArticle(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &GETArticleReq{}

	if err := c.BindUri(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	article, err := s.db.SelectArticleByID(req.ID)

	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select article failed. err: [%v]", err)
	}

	if article == nil {
		return req, nil, http.StatusNoContent, nil
	}

	accountMap, err := s.GetAccounts(rid, []int{article.AuthorID})
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get accounts failed. err: [%v]", err)
	}

	var avatar string
	author := "unknown"
	if accountMap != nil {
		if a, ok := accountMap[article.AuthorID]; ok {
			avatar = a.Avatar
			author = strings.Split(a.Email, "@")[0]
		}
	}

	return req, &Article{
		ID:       article.ID,
		AuthorID: article.AuthorID,
		Author:   author,
		Title:    article.Title,
		Tags:     strings.Split(article.Tags, ","),
		Content:  article.Content,
		CTime:    article.CTime.Format(time.RFC3339),
		UTime:    article.UTime.Format(time.RFC3339),
		Avatar:   avatar,
	}, http.StatusOK, nil
}
