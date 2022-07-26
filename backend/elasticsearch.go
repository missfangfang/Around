package backend // Backend matches DAO in OnlineOrder

import (
	"context"
	"fmt"

	"around/constants"

	"github.com/olivere/elastic/v7"
)

/* Create an index and reading data in Elasticsearch: https://github.com/olivere/elastic/blob/release-branch.v7/example_test.go */

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

	// 3. If index does not exist, create a new index
	if !exists {
		// If searching with "id" or "user", it needs to completely match to return a result (keyword: select * from post where id = "123")
		// If searching with "message", search result will return posts that contain the message or part of the message (text: select * from post where message contains/like "%tiffany%")
		// "index": false -> determine if the property needs to be indexed. "id" O(1), "user" O(logn) and "message" is indexed, "url" O(n) and "type" is not. Indexing does not intefere with how search is performed, but the effectiveness of the search.
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
		// createIndex, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background()) returns 2 results (createIndex and err) but we don't need createIndex so we can substitute it with _
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
	fmt.Println("Indexes are created.") // For debugging purposes.

	// New an object
	// & because *ElasticsearchBackend, : means to initialize
	// client (private property/myclient): client (ES connection/esclient)
	// Java equivalent:
	// class ElasticsearchBackend {
	// 	private Client myclient;
	// 	ElasticsearchBackend() {}
	// }
	// ESBackend = new ElasticsearchBackend(esclient)
	ESBackend = &ElasticsearchBackend{client: client}
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	searchResult, err := backend.client.Search().
		Index(index).            // Search in index
		Query(query).            // Specify the query
		Pretty(true).            // Pretty print request and response JSON
		Do(context.Background()) // Execute
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}
