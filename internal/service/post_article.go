package service

import (
	"fmt"
	"github.com/hpifu/go-tech/internal/es"
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
	Brief    string   `json:"brief,omitempty"`
	Content  string   `form:"content" json:"content,omitempty"`
	CTime    string   `json:"ctime,omitempty"`
	UTime    string   `json:"utime,omitempty"`
	Avatar   string   `json:"avatar,omitempty"`
}

type POSTArticleReq Article

type POSTArticleRes struct {
	ID int `json:"id,omitempty"`
}

func (s *Service) POSTArticle(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &POSTArticleReq{
		Token: c.GetHeader("Authorization"),
	}

	// get account
	account, err := s.accountCli.GETAccountToken(rid, req.Token)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("get account failed. err: [%v]", err)
	}
	if account == nil {
		return req, nil, http.StatusForbidden, fmt.Errorf("没有该资源权限")
	}

	// bind request
	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if err := s.validPOSTArticle(req); err != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("valid request failed. err: [%v]", err)
	}

	req.AuthorID = account.ID

	// check if article exists
	dbArticle, err := s.db.SelectArticleByAuthorAndTitle(req.AuthorID, req.Title)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select article failed. err: [%v]", err)
	}
	if dbArticle != nil {
		return req, nil, http.StatusBadRequest, fmt.Errorf("文章已存在")
	}

	// insert article
	reqArticle := &mysql.Article{
		AuthorID: req.AuthorID,
		Tags:     strings.Join(req.Tags, ","),
		Title:    req.Title,
		Content:  req.Content,
		Brief:    runecut(req.Content, 60),
		CTime:    time.Now(),
		UTime:    time.Now(),
	}
	if err := s.db.InsertArticle(reqArticle); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql insert article failed. err: [%v]", err)
	}

	// select article id
	dbArticle, err = s.db.SelectArticleByAuthorAndTitle(req.AuthorID, req.Title)
	if err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select article failed. err: [%v]", err)
	}
	if dbArticle == nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql select article failed. err: [not found]")
	}

	for _, tag := range req.Tags {
		if err := s.db.InsertTag(tag, dbArticle.ID); err != nil {
			return req, nil, http.StatusInternalServerError, fmt.Errorf("mysql insert tag failed. err: [%v]", err)
		}
	}

	esArticle := &es.Article{
		ID:      dbArticle.ID,
		Title:   reqArticle.Title,
		Author:  strings.Split(account.Email, "@")[0],
		Tags:    reqArticle.Tags,
		Content: reqArticle.Content,
	}
	if err := s.es.InsertArticle(esArticle); err != nil {
		return req, nil, http.StatusInternalServerError, fmt.Errorf("es insert article failed. err: [%v]", err)
	}

	return req, &POSTArticleRes{
		ID: dbArticle.ID,
	}, http.StatusCreated, nil
}

func (s *Service) validPOSTArticle(req *POSTArticleReq) error {
	if err := rule.Check([][3]interface{}{
		{"标题", req.Title, []rule.Rule{rule.Required, rule.AtMost(128)}},
	}); err != nil {
		return err
	}

	return nil
}
