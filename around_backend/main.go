package main

import (
    "fmt"
    "log"
    "net/http" 

    "around/backend"
    "around/handler"  
    "around/util" 
)
func main() {
    fmt.Println("started-service")
    
    config, err := util.LoadApplicationConfig("conf", "deploy.yml")
    if err != nil {
        panic(err)
    }

    backend.InitElasticsearchBackend(config.ElasticsearchConfig) //自动连接数据库/新建一个数据库
    backend.InitGCSBackend(config.GCSConfig) //初始化GCS
    
    log.Fatal(http.ListenAndServe(":8080", handler.InitRouter(config.TokenConfig))) //listen 8080
}