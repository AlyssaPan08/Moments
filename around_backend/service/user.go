package service

import (
    "fmt"

    "around/backend"
    "around/constants"
    "around/model"

    "github.com/olivere/elastic/v7"
)

//Check username & password
//Method 1: Read ES based both on username and password to check if there's a hit
//Method 2: Read ES based on username, compare the given password with the password returned by ES
func CheckUser(username string, password string) (bool, error) {
//bool: false if username/password wrong
    //method 1
    query := elastic.NewBoolQuery()
    query.Must(elastic.NewTermQuery("username", username))
    query.Must(elastic.NewTermQuery("password", password))
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

    if searchResult.TotalHits() > 0 {
        return true, nil
    }

    return false, nil

    //method 2:
    // query := elastic.NewTermQuery("username", user.Username)
    // searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    // if err != nil {
    //     return false, err
    // }

    // //借鉴post.go func getPostFromSearchResult
    // var utype model.User
    // for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
    //     u := item.(model.User)
    //     if u.Password == password {
    //         return true, nil
    //     }
    // }
    // return false, nil
}

//add user to ElasticSearch
//if add same username, ElasticSearch would update user rather than decline
func AddUser(user *model.User) (bool, error) {
    query := elastic.NewTermQuery("username", user.Username)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

    if searchResult.TotalHits() > 0 { //按照当前username search， 如果存在说明冲突了，return false
        return false, nil
    }

	//不存在 => insert
    err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
    if err != nil {
        return false, err
    }
    fmt.Printf("User is added: %s\n", user.Username)
    return true, nil
}