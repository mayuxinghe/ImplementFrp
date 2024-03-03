package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

func main() {
	// 1. 建立与二级代理之前的连接 ====================
	// 供二级代理连接的监听
	proxyLn, err := net.Listen("tcp", ":7000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer proxyLn.Close()

	fmt.Println("服务器正在监听代理端口...")

	// 接受二级代理连接
	proxyConn, err := proxyLn.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer proxyConn.Close()

	// 2. 建立与用户之间的连接 =========================
	// 监听端口8000上的供用户连接
	userLn, err := net.Listen("tcp", ":8900")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer userLn.Close()

	fmt.Println("服务器正在监听用户端口...")

	// 接受用户的连接
	userConn, err := userLn.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 3. 建立二级代理与用户连接之前的隧道通讯 ===========
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
			fmt.Println("转发数据发生错误:", err)
		}
	}

	wait.Add(2)
	go pipe(proxyConn, userConn) // 转发用户请求数据到二级代理服务器
	go pipe(userConn, proxyConn) // 转发二级代理数据到用户连接
	wait.Wait()
}
