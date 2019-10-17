package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DELETEArticleReq struct {
	Token string `json:"token,omitempty"`
	ID    int    `uri:"id" json:"id"`
}

type DELETEArticleRes struct{}

func (s *Service) DELETEArticle(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &DELETEArticleReq{
		Token: c.GetHeader("Authorization"),
	}

	if req.Token == "" {
		return req, nil, http.StatusBadRequest, fmt.Errorf("验证信息有误")
	}

	// select account
	account, err := s.client.GETAccount(req.Token, rid)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get account failed. err: [%v]", err)
	}
	if account == nil {
		return req, nil, http.StatusForbidden, fmt.Errorf("没有该资源权限")
	}

	// bind req
	if err := c.BindUri(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	// select article
	article, err := s.db.SelectArticleByID(req.ID)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select article failed. err: [%v]", err)
	}
	if article == nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("未找到该资源")
	}

	// check authorization
	if article.AuthorID != account.ID {
		return req, nil, http.StatusForbidden, fmt.Errorf("没有该资源权限")
	}

	// delete article
	if err := s.db.DeleteArticleByID(req.ID); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql update article failed. err: [%v]", err)
	}

	return req, nil, http.StatusAccepted, nil
}
