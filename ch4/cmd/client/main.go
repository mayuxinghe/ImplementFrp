package main

import (
    "fmt"
    "net"
    "os"
)
/*
* 客户端
*/
func main() {
    // 连接到服务器
    // old conn, err := net.Dial("tcp", "localhost:8080")
    // new 通过命名行参数获取服务端地址
    conn, err := net.Dial("tcp", os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }
    // 程序结束关闭连接
    defer conn.Close()
    fmt.Println("连接到服务器：" + os.Args[1])
    // 向服务器发送数据
    _, err = conn.Write([]byte("你好，服务器！"))
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("消息发送成功！")
    
    // 接收服务端的回复消息
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("接收消息失败:", err)
		return
	}

	// 打印服务端回复的消息
	receivedMsg := string(buffer[:n])
	fmt.Println("服务端回复: ", receivedMsg)
}