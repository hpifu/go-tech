package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/rule"
	"github.com/hpifu/go-tech/internal/mysql"
)

type Article struct {
	AuthorID int    `form:"authorID" json:"authorID,omitempty"`
	Author   string `form:"author" json:"author,omitempty"`
	Title    string `form:"title" json:"title,omitempty"`
	Content  string `form:"content" json:"content,omitempty"`
}

type POSTArticleReq Article

type POSTArticleRes struct{}

func (s *Service) POSTArticle(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &POSTArticleReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validPOSTArticle(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	if err := s.db.InsertArticle(&mysql.Article{
		AuthorID: req.AuthorID,
		Author:   req.Author,
		Title:    req.Title,
		Content:  req.Content,
	}); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql insert article failed. err: [%v]", err)
	}

	return req, nil, http.StatusCreated, nil
}

func (s *Service) validPOSTArticle(req *POSTArticleReq) error {
	if err := rule.Check([][3]interface{}{
		{"标题", req.Title, []rule.Rule{rule.Required, rule.AtMost(128)}},
	}); err != nil {
		return err
	}

	return nil
}
