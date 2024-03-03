package main // 声明包名，main 是可执行包
// 引入外部模块
import (
    "fmt"
    "net"
)
/*
* 服务端
*/
func main() {
    // 监听端口8000上的传入连接
    ln, err := net.Listen("tcp", ":8000")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer ln.Close()

    fmt.Println("服务器正在监听8000端口...")

    for {
        // 接受传入的连接
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        // 在新的goroutine中处理连接
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    // 从连接中读取数据
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println(err)
        return
    }
    // 打印接收到的数据
    fmt.Printf("收到客户端 %s 消息: %s", conn.RemoteAddr(), buf)
    // 向客户端发送数据
    receivedMsg := string(buf[:n])
	conn.Write([]byte("你好客户端，收到消息: " + receivedMsg))
}