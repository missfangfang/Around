package service

import (
	"fmt"
	"reflect"

	"around/backend"
	"around/constants"
	"around/model" // Contains user struct

	"github.com/olivere/elastic/v7"
)

func CheckUser(username, password string) (bool, error) {
	// Option 1:
	// Read ES based on username
	// Compare the given password with the password returned from ES
	// query := elastic.NewTermQuery("username", user.Username)
	// searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	// if err != nil {
	// 	return false, err
	// }
	// var utype model.User
	// for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
	// 	u := item.(model.User)
	// 	if u.Password == password {
	// 		return true, nil
	// 	}
	// }

	// return false, nil

	// Option 2:
	// Read ES based on both username and password
	// Check if there is a hit
	// query := elastic.NewBoolQuery()
	// query.Must(elastic.NewTermQuery("username", username)) // Must fulfill this condition
	// query.Must(elastic.NewTermQuery("password", password)) // Must also fulfill this condition

	// searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	// if err != nil {
	// 	return false, err
	// }

	// if searchResult.TotalHits() > 0 {
	// 	return true, err
	// }

	// return false, nil

	// Option 3:
	// Read ES based on both username and password
	// Compare the given password with the password returned from ES
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("username", username)) // Must fulfill this condition
	query.Must(elastic.NewTermQuery("password", password)) // Must also fulfill this condition

	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	// When you have 2 query.Must conditions like below, it's not necessary to check the given password with the password returned from ES
	var utype model.User
	for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		u := item.(model.User)
		if u.Password == password {
			fmt.Printf("Login as %s\n", username)
			return true, nil
		}
	}

	return false, nil
}

func AddUser(user *model.User) (bool, error) { // bool: true or false, can the user be added to DB?
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	// If search result is not empty, return false
	if searchResult.TotalHits() > 0 {
		return false, nil
	}

	// If search result is empty, save to DB
	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	fmt.Printf("User is added: %s\n", user.Username)
	return true, nil
}
