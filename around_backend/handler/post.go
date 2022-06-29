package handler
//controller
import (
    "encoding/json"
    "fmt"
    "net/http"
    "path/filepath"

    "around/service"
    "around/model"

    "github.com/pborman/uuid"
    "github.com/gorilla/mux"
    jwt "github.com/form3tech-oss/jwt-go"
)

//map media type
var (
    mediaTypes = map[string]string{
        ".jpeg": "image",
        ".jpg":  "image",
        ".gif":  "image",
        ".png":  "image",
        ".mov":  "video",
        ".mp4":  "video",
        ".avi":  "video",
        ".flv":  "video",
        ".wmv":  "video",
    }
)

// func uploadHandler(w http.ResponseWriter, r *http.Request) {
//     // Parse from body of request to get a json object.
//     fmt.Println("Received one upload request") //for debug
//     decoder := json.NewDecoder(r.Body) //read json request body to decoder
//     var p model.Post
//     if err := decoder.Decode(&p); err != nil { //use Decode method of decoder to convert JSON to Post类型的对象
//         panic(err) //try catch
//     }

//     fmt.Fprintf(w, "Post received: %s\n", p.Message) //Fprintf: 把第二个参数打印到第一个参数里面
// }

//重写：前端body返回的不是json,而是form格式的data, 因此通过request.FormValue来读而不是request.body
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one upload request")

    user := r.Context().Value("user") //get token
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

    //step 1: 从request body里面读数据 (new a Post object with struct properties)
    p := model.Post{
        Id:      uuid.New(),  //uuid lib: 随机生成unique id
        User:    username.(string),  //decode from token claim/payload
        Message: r.FormValue("message"), //从body里面读取form data
        //Url: from uploaded file; type: 根据文件类型推断 => 省略
    }

    //step 2: read file
    file, header, err := r.FormFile("media_file") //header: metadata of file
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest) //StatusBadRequest：client端的责任（4开头）
        fmt.Printf("Media file is not available %v\n", err)
        return
    }

    //step 3: 通过上传文件的后缀 map mediaType
    suffix := filepath.Ext(header.Filename)
    if t, ok := mediaTypes[suffix]; ok { //read hashmap; t: value, ok: exist
        p.Type = t
    } else {
        p.Type = "unknown"
    }

    //step 4: save post 
    err = service.SavePost(&p, file)
    if err != nil {
        http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
        fmt.Printf("Failed to save post to backend %v\n", err)
        return
    }

    fmt.Println("Post is saved successfully.")
}

//search
func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for search")
    w.Header().Set("Content-Type", "application/json") //explicitly tell frontend the return type 

    user := r.URL.Query().Get("user") //query(): URL包含参数的部分 eg., user=vincent
    keywords := r.URL.Query().Get("keywords")

    var posts []model.Post
    var err error
    if user != "" {
        posts, err = service.SearchPostsByUser(user)
    } else {
        posts, err = service.SearchPostsByKeywords(keywords)
    }

    //handler不能甩锅，必须handle error
    if err != nil {
        http.Error(w, "Failed to read post from backend", http.StatusInternalServerError) //StatusInternalServerError = 500
        fmt.Printf("Failed to read post from backend %v.\n", err)
        return
    }

    //没有error的话，用Go自带的lib Marshal把golang转成JSON string
    js, err := json.Marshal(posts)
    //handle parse过程中出现的error
    if err != nil {
        http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
        fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
        return
    }
    w.Write(js)
}

//delete
func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for delete")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

    if err := service.DeletePost(id, username); err != nil {
        http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete post from backend %v\n", err)
        return
    }
    fmt.Println("Post is deleted successfully")
}