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

type Article struct {
	Token    string   `json:"token,omitempty"`
	ID       int      `json:"id,omitempty" uri:"id"`
	AuthorID int      `json:"authorID,omitempty"`
	Author   string   `json:"author,omitempty"`
	Title    string   `form:"title" json:"title,omitempty"`
	Tags     []string `form:"tags" json:"tags,omitempty"`
	Content  string   `form:"content" json:"content,omitempty"`
	CTime    string   `json:"ctime,omitempty"`
	UTime    string   `json:"utime,omitempty"`
}

type POSTArticleReq Article

type POSTArticleRes struct{}

func (s *Service) POSTArticle(c *gin.Context) (interface{}, interface{}, int, error) {
	req := &POSTArticleReq{
		Token: c.GetHeader("Authorization"),
	}

	account, err := s.getAccount(req.Token)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, fmt.Errorf("get account failed. err: [%v]", err)
	}
	if account == nil {
		return nil, nil, http.StatusForbidden, fmt.Errorf("authorization failed")
	}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validPOSTArticle(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	req.AuthorID = account.ID
	req.Author = strings.Split(account.Email, "@")[0]

	if err := s.db.InsertArticle(&mysql.Article{
		AuthorID: req.AuthorID,
		Author:   req.Author,
		Tags:     strings.Join(req.Tags, ", "),
		Title:    req.Title,
		Content:  req.Content,
		CTime:    time.Now(),
		UTime:    time.Now(),
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
