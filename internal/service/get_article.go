package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	godtoken "github.com/hpifu/go-godtoken/api"
)

type ArticleReq struct {
	ID int `uri:"id" json:"id"`
}

func (s *Service) GETArticle(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ArticleReq{}

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

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	res, err := s.godtokenCli.GetToken(ctx, &godtoken.GetTokenReq{Rid: rid})
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("godtoken verify failed. err: [%v]", err)
	}

	accounts, err := s.accountCli.GETAccounts(rid, res.Token, []int{article.AuthorID})
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get accounts failed. err: [%v]", err)
	}

	var avatar string
	author := "unknown"
	if accounts != nil && len(accounts) != 0 {
		avatar = accounts[0].Avatar
		author = strings.Split(accounts[0].Email, "@")[0]
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
