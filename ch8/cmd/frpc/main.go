package main

import (
	"ch8/models/client"
	"fmt"
	"io"
	"net"
	"sync"
)

func main() {
	// new  0. 加载配置 ======================
	client.LoadConf("./frpc.ini")
	// 外网代理服务器
	publicSvrHost := client.ServerAddr
	publicSvrPort := client.ServerPort
	// 目标应用服务器
	targetSvrHost := client.LocalIp
	targetSvrPort := client.LocalPort
	// 1. 连接到一层代理（frps） =============
	pubSvrAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", publicSvrHost, publicSvrPort))
	if err != nil {
		return
	}
	// 连接一层代理服务器
	pubCoon, err := net.DialTCP("tcp4", nil, pubSvrAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	// 2. 连接内网目标服务器 =================
	targetCoon, err := net.Dial("tcp", fmt.Sprintf("%s:%d", targetSvrHost, targetSvrPort))
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	// 3. 建立一层代理和目标服务器之间的隧道通讯 =================
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
	go pipe(targetCoon, pubCoon) // 转发一层代理数据到目标服务器
	go pipe(pubCoon, targetCoon) // 转发目标服务器数据到一层代理
	wait.Wait()
}
