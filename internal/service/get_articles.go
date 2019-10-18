package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hpifu/go-account/pkg/account"
	godtoken "github.com/hpifu/go-godtoken/api"

	"github.com/gin-gonic/gin"
)

type ArticlesReq struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

type ArticlesRes []*Article

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
		return req, nil, http.StatusUnauthorized, nil
	}

	var ids []int
	idMap := map[int]struct{}{}
	for _, article := range articles {
		idMap[article.AuthorID] = struct{}{}
	}
	for k := range idMap {
		ids = append(ids, k)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	res, err := s.godtokenCli.GetToken(ctx, &godtoken.GetTokenReq{Rid: rid})
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("godtoken verify failed. err: [%v]", err)
	}

	accounts, err := s.accountCli.GETAccounts(rid, res.Token, ids)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get accounts failed. err: [%v]", err)
	}
	accountMap := map[int]*account.Account{}
	if accounts != nil {
		for _, a := range accounts {
			accountMap[a.ID] = a
		}
	}

	var as ArticlesRes
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
			Content:  article.Content,
			CTime:    article.CTime.Format(time.RFC3339),
			UTime:    article.UTime.Format(time.RFC3339),
			Avatar:   avatar,
		})
	}

	return req, as, http.StatusOK, nil
}
