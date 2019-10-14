package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
	"github.com/hpifu/go-tech/internal/mysql"
)

type PUTArticleReq Article

type PUTArticleRes struct{}

func (s *Service) PUTArticle(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &PUTArticleReq{
		Token: c.GetHeader("Authorization"),
	}

	// select account
	account, err := s.getAccount(req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get account failed. err: [%v]", err)
	}
	if account == nil {
		return req, nil, http.StatusForbidden, fmt.Errorf("access deny")
	}

	// bind req
	if err := c.Bind(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}
	if err := c.BindUri(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("bind uri failed. err: [%v]", err)
	}

	if err := s.validPUTArticle(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	// select article
	article, err := s.db.SelectArticleByID(req.ID)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select article failed. err: [%v]", err)
	}
	if article == nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("no such article")
	}

	if article.AuthorID != account.ID {
		return req, nil, http.StatusForbidden, fmt.Errorf("access deny")
	}

	req.Author = strings.Split(account.Email, "@")[0]

	if err := s.db.UpdateArticle(&mysql.Article{
		ID:      req.ID,
		Author:  req.Author,
		Tags:    strings.Join(req.Tags, ", "),
		Title:   req.Title,
		Content: req.Content,
		CTime:   time.Now(),
		UTime:   time.Now(),
	}); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql update article failed. err: [%v]", err)
	}

	return req, nil, http.StatusAccepted, nil
}

func (s *Service) validPUTArticle(req *PUTArticleReq) error {
	if req.Title != "" {
		if err := rule.Check([][3]interface{}{
			{"标题", req.Title, []rule.Rule{rule.Required, rule.AtMost(128)}},
		}); err != nil {
			return err
		}
	}

	return nil
}
