package main

import (
	"ch9/models/server"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	docopt "github.com/docopt/docopt-go"
)
var configFile = "./frps.ini"
var usage string = `代理服务端程序的使用方法

Usage: 
	frps [-c config_file] [--addr=<bind_addr>]
	frps -h | --help | --version

Options:
	-c config_file            设置配置文件路径
	--addr=<bind_addr>        代理监听地址，形式如: 0.0.0.0:7000
	-h --help                 显示本帮助信息
	--version                 显示版本号
`

func main() {
	// 解析命令行参数
	args, err := docopt.Parse(usage, nil, true, "v frps_day6", false)
	// 根据 -c 选项获取命令行指定配置文件
	if args["-c"] != nil {
		configFile = args["-c"].(string)
	}
	err = server.LoadConf(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if args["--addr"] != nil {
		addr := strings.Split(args["--addr"].(string), ":")
		if len(addr) != 2 {
			fmt.Println("--addr format error: example 0.0.0.0:7000")
			os.Exit(1)
		}
		bindPort, err := strconv.ParseInt(addr[1], 10, 64)
		if err != nil {
			fmt.Println("--addr format error, example 0.0.0.0:7000")
			os.Exit(1)
		}
		server.BindAddr = addr[0]
		server.BindPort = bindPort
	}
	
	// 供代理客户端使用的监听
	// proxyLn, err := net.Listen("tcp", ":7000")
	proxyLn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.BindAddr, server.BindPort))
	
	if err != nil {
		fmt.Println(err)
		return
	}
	defer proxyLn.Close()

	fmt.Printf("服务器正在监听 %d 端口...\n", server.BindPort)

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

	fmt.Printf("服务器正在监听 %d 端口.../n", server.ListenPort)

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
