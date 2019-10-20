package es

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title,omitempty"`
	Author  string `json:"author,omitempty"`
	Tags    string `json:"tags,omitempty"`
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

func (e *ES) SearchArticle(value string, offset int, limit int) ([]*Article, error) {
	query := elastic.NewBoolQuery()
	for _, val := range split(value) {
		if len(val) == 0 {
			continue
		}
		q := elastic.NewBoolQuery()
		q.Should(elastic.NewTermQuery("title", val))
		q.Should(elastic.NewTermQuery("author", val))
		q.Should(elastic.NewTermQuery("dynasty", val))
		q.Should(elastic.NewTermQuery("content", val))
		query.Must(q)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	defer cancel()
	res, err := e.es.Search().
		Index("article").
		Query(query).
		From(offset).Size(limit).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	var ancient Article
	var ancients []*Article
	for _, item := range res.Each(reflect.TypeOf(ancient)) {
		if t, ok := item.(Article); ok {
			ancients = append(ancients, &t)
		}
	}

	return ancients, err
}
