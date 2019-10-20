package es

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
)

var articleMapping = `{
	"mappings": {
		"properties": {
			"id": {
				"type": "long"
			},
			"title": {
				"type": "text",
				"analyzer": "ik_max_word",
				"search_analyzer": "ik_smart"
			},
			"author": {
				"type": "text",
				"analyzer": "ik_max_word",
				"search_analyzer": "ik_smart"
			},
			"tags": {
				"type": "text",
				"analyzer": "ik_max_word",
				"search_analyzer": "ik_smart"
			},
			"content": {
				"type": "text",
				"analyzer": "ik_max_word",
				"search_analyzer": "ik_smart"
			}
		}
	}
}`

type ES struct {
	es      *elastic.Client
	timeout time.Duration
	index   string
}

func NewES(uri string, index string, timeout time.Duration) (*ES, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(uri),
		elastic.SetSniff(false),
	)

	if err != nil {
		return nil, err
	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		defer cancel()
		_, _, err = client.Ping(uri).Do(ctx)
		if err != nil {
			return nil, err
		}
	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		defer cancel()
		exists, err := client.IndexExists(index).Do(ctx)
		if err != nil {
			return nil, err
		}
		if !exists {
			createIndex, err := client.CreateIndex(index).Body(articleMapping).Do(context.Background())
			if err != nil {
				return nil, err
			}
			if !createIndex.Acknowledged {
				return nil, fmt.Errorf("not acknowledged")
			}
		}
	}

	return &ES{
		es:      client,
		timeout: timeout,
		index:   index,
	}, nil
}
