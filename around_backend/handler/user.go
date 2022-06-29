package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "regexp"
    "time"

    "around/model"
    "around/service"

    jwt "github.com/form3tech-oss/jwt-go" //create token; jwt: import as "jwt"
)

//initialize it in another place
//var mySigningKey = []byte("secret") //安全密钥 (我们使用对称加密，即加密和解密的安全密钥是同一个 => input string can be any string)

//read user info from request, return token (login status) if success
func signinHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signin request")
    w.Header().Set("Content-Type", "text/plain")

    if r.Method == "OPTIONS" {
        return
    }

    //step 1: Get User information from client
    decoder := json.NewDecoder(r.Body) //JSON -> user
    var user model.User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

	//step 2: check user
    exists, err := service.CheckUser(user.Username, user.Password)
    if err != nil {
        http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
        return
    }

    if !exists {
        http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized) //401
        fmt.Printf("User doesn't exists or wrong password\n")
        return
    }

	//step 3: create token by jwt
	//Two types of claim: StandardClaim: claim name 严格限制； CustomClaim: 自定义 claim name 
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //claim: user's data/payload（紫色部分）
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),	//Unix(): Unix epoch time, standard timestamp
    })
	//为什么不放password: header & payload容易被解码 => unsecure

    tokenString, err := token.SignedString(mySigningKey) //valid signature(蓝色部分) 通过安全密钥加密
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        fmt.Printf("Failed to generate token %v\n", err)
        return
    }

    w.Write([]byte(tokenString))
}

//read user info from request body (JSON), turn to go obj (user), save to ES
func signupHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signup request")
    w.Header().Set("Content-Type", "text/plain")

	//step 1: decode JSON to user
    decoder := json.NewDecoder(r.Body) 
    var user model.User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

	//step 2: validate input string (这个功能在前端check更好)
    if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }

	//step 3: add user to ES
    success, err := service.AddUser(&user)
    if err != nil {
        http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
        return
    }

    if !success {
        http.Error(w, "User already exists", http.StatusBadRequest)
        fmt.Println("User already exists")
        return
    }
    fmt.Printf("User added successfully: %s.\n", user.Username)
}