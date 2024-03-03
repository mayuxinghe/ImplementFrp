package main

import (
    // "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
)

func handleRequestAndRedirect(w http.ResponseWriter, req *http.Request) {
    // 创建一个新的HTTP请求，与原始请求具有相同的方法、URL和主体
    targetURL, _ := url.Parse("http://localhost:8000") // 替换为你要代理的目标服务器地址和端口号
	targetURL.Path = req.URL.Path
	targetURL.RawQuery = req.URL.RawQuery
    proxyReq, err := http.NewRequest(req.Method, targetURL.String(), req.Body)
    if err != nil {
        http.Error(w, "创建代理请求失败。", http.StatusInternalServerError)
        return
    }

    // 将原始请求的标头复制到代理请求
    for name, values := range req.Header {
        for _, value := range values {
            proxyReq.Header.Add(name, value)
        }
    }

    // 使用默认传输发送代理请求
    resp, err := http.DefaultTransport.RoundTrip(proxyReq)
    if err != nil {
        http.Error(w, "发送代理请求失败。", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // 将代理响应的标头复制到原始响应
    for name, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(name, value)
        }
    }

    // 将原始响应的状态代码设置为代理响应的状态代码
    w.WriteHeader(resp.StatusCode)

    // 将代理响应的主体复制到原始响应
    io.Copy(w, resp.Body)
}

func main() {
    // 创建一个新的HTTP服务器，使用handleRequestAndRedirect函数作为处理程序
    server := http.Server{
        Addr:    ":8080",
        Handler: http.HandlerFunc(handleRequestAndRedirect),
    }

    // 启动服务器并记录任何错误
    log.Println("启动代理服务器在 :8080")
    err := server.ListenAndServe()
    if err != nil {
        log.Fatal("启动代理服务器时发生错误: ", err)
    }
}