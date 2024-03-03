package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func handleConnection(client net.Conn, target string) {
    // 连接到目标服务器
    server, err := net.Dial("tcp", target)
    if err != nil {
        log.Println("无法连接到目标服务器:", err)
        client.Close()
        return
    }

    var wait sync.WaitGroup
	pipe := func(to net.Conn, from net.Conn) {
        // 资源清理
		defer to.Close()
		defer from.Close()
		defer wait.Done()

        // 转发数据
		var err error
		_, err = io.Copy(to, from)
		if err != nil {
			log.Fatal("join connections error, %v", err)
		}
	}

	wait.Add(2)
	go pipe(server, client)
	go pipe(client, server)
	wait.Wait()
}

func main() {
    // 代理服务监听端口
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal("无法启动代理服务器:", err)
    }
    defer listener.Close()

    log.Println("启动代理服务器在 :8080")

    // 接受客户端连接并处理
    for {
        client, err := listener.Accept()
        if err != nil {
            log.Println("无法接受客户端连接:", err)
            continue
        }
        go handleConnection(client, "localhost:8000") // 替换为你想要转发的目标地址
    }
}