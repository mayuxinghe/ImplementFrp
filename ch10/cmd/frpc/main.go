package main

import (
	"ch10/models/client"
	"ch10/utils/log"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/docopt/docopt-go"
)
var configFile = "./frpc.ini"
var usage string = `frpc 客户端用法

Usage: 
    frpc [-c config_file] [-L log_file] [--log-level=<log_level>] [--server-addr=<server_addr>]
    frpc -h | --help | --version

Options:
    -c config_file              指定配置文件
	-L log_file                 指定日志输出文件, including console
    --log-level=<log_level>     指定日志等级: debug, info, warn, error
    --server-addr=<server_addr> 服务端监听地址，形式为: 0.0.0.0:7000
    -h --help                   显示本帮助信息
    --version                   显示版本号
`

func main() {
	// 0. 加载参数
	args, err := docopt.Parse(usage, nil, true, "V frpc_day6", false)

	if args["-c"] != nil {
		configFile = args["-c"].(string)
	}
	err = client.LoadConf(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if args["--server-addr"] != nil {
		addr := strings.Split(args["--server-addr"].(string), ":")
		if len(addr) != 2 {
			fmt.Println("--server-addr format error: example 0.0.0.0:7000")
			os.Exit(1)
		}
		serverPort, err := strconv.ParseInt(addr[1], 10, 64)
		if err != nil {
			fmt.Println("--server-addr format error, example 0.0.0.0:7000")
			os.Exit(1)
		}
		client.ServerAddr = addr[0]
		client.ServerPort = serverPort
	}

	if args["-L"] != nil {
		if args["-L"].(string) == "console" {
			client.LogWay = "console"
		} else {
			client.LogWay = "file"
			client.LogFile = args["-L"].(string)
		}
	}

	if args["--log-level"] != nil {
		client.LogLevel = args["--log-level"].(string)
	}
	// 0.1 初始化
	log.InitLog(client.LogWay, client.LogFile, client.LogLevel)

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
		log.Error("Error connecting:", err)
		return
	}

	// 2. 连接内网目标服务器 =================
	targetCoon, err := net.Dial("tcp", fmt.Sprintf("%s:%d", targetSvrHost, targetSvrPort))
	if err != nil {
		log.Error("Error connecting:", err)
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
            log.Info("从目标服务器到客户端的数据传输错误:", err)
        }
	}

	wait.Add(2)
	go pipe(targetCoon, pubCoon) // 转发一层代理数据到目标服务器
	go pipe(pubCoon, targetCoon) // 转发目标服务器数据到一层代理
	log.Warn("启动 frpc 客户端成功!")
	wait.Wait()
	log.Warn("All proxy exit!")
}
