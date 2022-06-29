package backend

import (
    "context"
    "fmt"
    "io"

    "around/util"

    "cloud.google.com/go/storage"
)

var (
    GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
    client *storage.Client
    bucket string
}

func InitGCSBackend(config *util.GCSInfo) {
    client, err := storage.NewClient(context.Background())
    if err != nil {
        panic(err)
    }

    GCSBackend = &GoogleCloudStorageBackend{
        client: client,
        bucket: config.Bucket,
    }
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) { 
//r: 从controller里面读出来的文件；objectName:写到哪个文件里；返回string: url,存到ElasticSearch
    ctx := context.Background()
    object := backend.client.Bucket(backend.bucket).Object(objectName) //在gcs里面创建一个文件夹(bucket)，在里面创建一个新的文件(object)
    //step 1: upload file
    wc := object.NewWriter(ctx) 
    if _, err := io.Copy(wc, r); err != nil {  //io.copy: upload an object; wc: target; r: source
        return "", err
    }

    if err := wc.Close(); err != nil {
        return "", err
    }

    //step 2: 修改权限：private -> alluser can read
    if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil { //ACL():access control list 权限控制
        return "", err
    }

    //step 1: get and return url
    attrs, err := object.Attrs(ctx) //Attrs(): 返回object属性(包括url)
    if err != nil {
        return "", err
    }

    fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
    return attrs.MediaLink, nil
}