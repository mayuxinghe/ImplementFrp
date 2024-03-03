package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 创建一个新的反向代理
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:8000", // 替换为你要代理的目标服务器地址和端口号
	})
	
	
	// 创建一个新的反向代理
	/*
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			targetURL, _ := url.Parse("http://目标服务器地址:端口号") // 替换为你要代理的目标服务器地址和端口号
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.Host = targetURL.Host
		},
	}
	*/
	// 更改请求头中的Host字段，确保目标服务器能够正确处理请求
	r.Host = "目标服务器地址" // 替换为你要代理的目标服务器地址

	// 调用ServeHTTP方法处理请求转发
	proxy.ServeHTTP(w, r)
}

func main() {
	// 设置代理路由和处理函数
	http.HandleFunc("/", handler)

	// 启动代理服务器
	fmt.Println("代理服务器已启动，监听端口：8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}