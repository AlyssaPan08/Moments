package service

import (
    "mime/multipart"
    "reflect"

    "around/backend"
    "around/constants"
    "around/model"

    "github.com/olivere/elastic/v7"
)

//按user生成query, 并用此query搜索
func SearchPostsByUser(user string) ([]model.Post, error) {
    query := elastic.NewTermQuery("user", user)
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

//按keywords生成query, 并用此query搜索
func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
    query := elastic.NewMatchQuery("message", keywords)
    query.Operator("AND") //"vincent+richard+sean" ES通过“+” 连接keyword
    if keywords == "" {  //corner case: no keyword => return all
        query.ZeroTermsQuery("all")
    }
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
    if err != nil {
        return nil, err
    }
    return getPostFromSearchResult(searchResult), nil
}

//把搜索结果加到result里面，最终return to client
func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
    var ptype model.Post    //a post obj
    var posts []model.Post  //a list of posts to return

    for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {  //reflect.TypeOf(ptype): NoSQL要先判断type是否一致
        p := item.(model.Post)
        posts = append(posts, p)
    }
    return posts
}

//save post metadata to GCS and ES
func SavePost(post *model.Post, file multipart.File) error { //multipart 约等于 http body文件
    medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
    if err != nil {
        return err
    }

    post.Url = medialink
    return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
    //online service vs offline service => 处理GCS和ES不一致的情况
    //先传到GCS, 再传到ES, 只要存ES成功就认为是成功，如果存ES失败user就会看到上传失败
    //因此只要以ES为准，定期（offline)检查GCS和ES的不一致情况并delete多的文件即可
}

//delete post
func DeletePost(id string, user string) error {
    query := elastic.NewBoolQuery()
    query.Must(elastic.NewTermQuery("id", id))
    query.Must(elastic.NewTermQuery("user", user))

    return backend.ESBackend.DeleteFromES(query, constants.POST_INDEX)
}
