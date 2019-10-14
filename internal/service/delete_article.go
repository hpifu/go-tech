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

func (s *Service) DELETEArticle(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &DELETEArticleReq{
		Token: c.GetHeader("Authorization"),
	}

	if req.Token == "" {
		return req, "验证信息有误", http.StatusBadRequest, nil
	}

	// select account
	account, err := s.getAccount(req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get account failed. err: [%v]", err)
	}
	if account == nil {
		return req, "没有该资源权限", http.StatusForbidden, nil
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
		return req, "未找到该资源", http.StatusBadRequest, nil
	}

	// check authorization
	if article.AuthorID != account.ID {
		return req, "没有该资源权限", http.StatusForbidden, nil
	}

	// delete article
	if err := s.db.DeleteArticleByID(req.ID); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql update article failed. err: [%v]", err)
	}

	return req, nil, http.StatusAccepted, nil
}
