package backend

import (
	"context"
	"fmt"

	"around/constants"

	"github.com/olivere/elastic/v7"
)

/* Create an index in Elasticsearch: https://github.com/olivere/elastic/blob/release-branch.v7/example_test.go */

/* ES mapping: https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping.html */

var ( // global variable
	ESBackend *ElasticsearchBackend
)

type ElasticsearchBackend struct { // DAO class
	client *elastic.Client // * pointer, refer to ES library codebase
}

func InitElasticsearchBackend() {
	// 1. New a ElasticsearchBackend object, create a connection
	client, err := elastic.NewClient(
		elastic.SetURL(constants.ES_URL), // URL address that connects to DB
		elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_PASSWORD))
	if err != nil {
		panic(err)
	}

	// 2. Determine if the DB index exist (more complicated than OnlineOrder MySQL)
	exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	// If index does not exist, create a new index
	if !exists {
		// If searching with "id" or "user", it needs to completely match to return a result (select * from post where id = "123")
		// If searching with "message", search result will return posts that contain the message or part of the message (select * from post where message contains/like "%tiffany%")
		// "index": determine if the property needs to be indexed
		mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },  
                    "user":     { "type": "keyword" },
                    "message":  { "type": "text" },
                    "url":      { "type": "keyword", "index": false },
                    "type":     { "type": "keyword", "index": false }
                }
            }
        }`
		_, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
                                        "password": {"type": "keyword"},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                }
                        }
                }`
		_, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Indexes are created.")

	ESBackend = &ElasticsearchBackend{client: client}
}
