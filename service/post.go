package service

import (
	"mime/multipart"
	"reflect"

	"around/backend"
	"around/constants"
	"around/model"

	"github.com/olivere/elastic/v7"
)

// Support user-based search and keyword-based search

func SearchPostsByUser(user string) ([]model.Post, error) { // Return result is an array
	query := elastic.NewTermQuery("user", user)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX) // Need to add global variable ESBackend
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
	query := elastic.NewMatchQuery("message", keywords)
	query.Operator("AND") // Returns posts that contain all keywords
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
	var ptype model.Post   // A post object
	var posts []model.Post // Return a list of posts

	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) { // Need to use ptype, not model.Post
		p := item.(model.Post)
		posts = append(posts, p)
	}
	return posts
}

// Save uploaded data to Elasticsearch and GCS
func SavePost(post *model.Post, file multipart.File) error { // 2 inputs: post and image file
	// SaveToGCS
	medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
	if err != nil {
		return err
	}
	post.Url = medialink

	// SaveToES
	return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
}

// Delete a post based on its ID and user
func DeletePost(id string, user string) error {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", id))
	query.Must(elastic.NewTermQuery("user", user))

	return backend.ESBackend.DeleteFromES(query, constants.POST_INDEX)
}
