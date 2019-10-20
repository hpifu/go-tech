package es

import (
	"context"
	"github.com/olivere/elastic/v7"
	"net/http"
	"strconv"
	"strings"
)

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title,omitempty"`
	Author  string `json:"author,omitempty"`
	Tags    string `json:"tags,omitempty"`
	Brief   string `json:"brief,omitempty"`
	Content string `json:"content,omitempty"`
}

func split(s string) []string {
	f := func(r rune) bool {
		for _, s := range []rune("，。、？！； ,.") {
			if r == s {
				return true
			}
		}
		return false
	}
	return strings.FieldsFunc(s, f)
}

func (e *ES) InsertArticle(article *Article) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	if _, err := e.es.Index().Index(e.index).Id(strconv.Itoa(article.ID)).BodyJson(article).Do(ctx); err != nil {
		return err
	}

	return nil
}

func (e *ES) DeleteArticle(id int) error {
	ctx, cancle := context.WithTimeout(context.Background(), e.timeout)
	defer cancle()
	if _, err := e.es.Delete().Index(e.index).Id(strconv.Itoa(id)).Do(ctx); err != nil && err.(*elastic.Error).Status != http.StatusNotFound {
		return err
	}

	return nil
}

func (e *ES) UpdateArticle(article *Article) error {
	ctx, cancle := context.WithTimeout(context.Background(), e.timeout)
	defer cancle()
	if _, err := e.es.Update().Index(e.index).Id(strconv.Itoa(article.ID)).Doc(article).Do(ctx); err != nil {
		if err.(*elastic.Error).Status == http.StatusNotFound {
			return e.InsertArticle(article)
		}
		return err
	}

	return nil
}

func (e *ES) SearchArticle(value string, offset int, limit int) ([]int, error) {
	query := elastic.NewBoolQuery()
	for _, val := range split(value) {
		if len(val) == 0 {
			continue
		}
		q := elastic.NewBoolQuery()
		q.Should(elastic.NewTermQuery("title", val))
		q.Should(elastic.NewTermQuery("author", val))
		q.Should(elastic.NewTermQuery("tags", val))
		q.Should(elastic.NewTermQuery("content", val))
		query.Must(q)
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	res, err := e.es.Search().
		Index(e.index).NoStoredFields().
		Query(query).
		From(offset).Size(limit).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	var ids []int
	for _, hit := range res.Hits.Hits {
		id, err := strconv.Atoi(hit.Id)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	return ids, err
}
