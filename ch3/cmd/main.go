package main // 包名

// 引入包
import (
    "io"
    "log"
    "net/http"
)
// 入口函数
func main() {
    // Hello world, the web server
    helloHandler := func(w http.ResponseWriter, req *http.Request) {
        io.WriteString(w, "Hello, world!\n")
    }

    // 设置路由方法
    http.HandleFunc("/hello", helloHandler)


    log.Println("Listening for requests at http://localhost:8000/hello")
    // 启动 web 服务监听在 8000端口，并打印返回信息
    log.Fatal(http.ListenAndServe(":8000", nil))
}