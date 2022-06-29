package handler

import (
    "net/http"

    "around/util"

    jwtmiddleware "github.com/auth0/go-jwt-middleware" 
    jwt "github.com/form3tech-oss/jwt-go"
    "github.com/gorilla/mux"  
    "github.com/gorilla/handlers"  
)

var mySigningKey []byte

//dispatcher servlet
func InitRouter(config *util.TokenInfo) http.Handler {  //*mux.Router => http.Handler 支持跨域访问
    mySigningKey = []byte(config.Secret)

    //middleware: between dispatcher servlet and service，validate token
    jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
        ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
            return []byte(mySigningKey), nil
        },
        SigningMethod: jwt.SigningMethodHS256,
    })

    router := mux.NewRouter()

    router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST") //need to validate token
    router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")
    router.Handle("/post/{id}", jwtMiddleware.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE")

    router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST") // do not need to validate token
    router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

    //前端发送跨域访问请求的时候（eg, AWS -> Google cloud)
    //default: 屏蔽所有跨域请求
    originsOk := handlers.AllowedOrigins([]string{"*"}) //只有后端声明的ip才支持跨域请求(*表示都支持)
    headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}) //只有包含这几种header的request支持跨域请求(Authorization是包含token的header;Content-Type = json)
    methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"}) //只有后端声明的methods才支持跨域请求

    return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}