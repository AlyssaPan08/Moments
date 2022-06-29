package backend
//Dao
import (
    "context"
    "fmt"

    "around/util"
    "around/constants"
    
    "github.com/olivere/elastic/v7"
)

//create a gloabl variable to make client obj singleton
var (
    ESBackend *ElasticsearchBackend
)

//new a class
type ElasticsearchBackend struct { //相当于Dao class
    client *elastic.Client //connect to database, 相当于SessionFactory
}

//new ElasticsearchBackend obj
func InitElasticsearchBackend(config *util.ElasticsearchInfo) { 
    //step 1: connect to database (& handle error)
	client, err := elastic.NewClient(
        elastic.SetURL(config.Address), //选择连接数据库的ip地址：ES_URL = "http://10.138.0.2:9200/"
        elastic.SetBasicAuth(config.Username, config.Password))
	
    if err != nil {
        panic(err)
    }

	//step 2: 判断数据库是否存在，不存在就创建一个新的(& handle error)
	//2.1 post index
    exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
    if err != nil {
        panic(err)
    }

    //not exist => create a new schema
    //select * from post where id = "123" => keyword 完全匹配
	//select * from post where message contains xxx => text：模糊匹配
	//"index" : false => 是否要对url/type做搜索上的优化（用户按照id/uesr/message搜索TC = O(1), 按照url/type搜索TC = O(n)）
    if !exists {
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

	//2.2 user index
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
    fmt.Println("Indexes are created.") //help debug

	//Step 3: encapsulate client in ESBackend
    ESBackend = &ElasticsearchBackend{client: client}
}
	
//read 
func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
    searchResult, err := backend.client.Search().
        Index(index).    //search in index
        Query(query).
        Pretty(true).    //pretty print request and response JSON
        Do(context.Background())  //execute
    if err != nil {
        return nil, err
    }

    return searchResult, nil
}

//save
func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
//i interface{}: 相当于java object， 可以支持各种形式的data, 适用于user & post
    _, err := backend.client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())
    return err
}

//delete
func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
    _, err := backend.client.DeleteByQuery().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())

    return err
}