package client

import (
	"ch11/utils/log"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/vaughan0/go-ini"
)

var (
	// 代理服务端地址信息
	ServerAddr string = "0.0.0.0"
	ServerPort int64  = 7000
	// 日志参数
	LogFile  string = "console"
	LogWay   string = "console"
	LogLevel string = "info"
	// 应用服务端地址信息
	LocalIp   string
	LocalPort int64
	// 心跳常量
	HeartBeatInterval int64 = 20
	HeartBeatTimeout  int64 = 90
)

func LoadConf(confFile string) (err error) {
	var tmpStr string
	var ok bool

	conf, err := ini.LoadFile(confFile)
	if err != nil {
		return err
	}

	// common
	tmpStr, ok = conf.Get("common", "server_addr")
	if ok {
		ServerAddr = tmpStr
	}

	tmpStr, ok = conf.Get("common", "server_port")
	if ok {
		ServerPort, _ = strconv.ParseInt(tmpStr, 10, 64)
	}

	// common
	tmpStr, ok = conf.Get("app", "local_ip")
	if ok {
		LocalIp = tmpStr
	}

	tmpStr, ok = conf.Get("app", "local_port")
	if ok {
		LocalPort, _ = strconv.ParseInt(tmpStr, 10, 64)
	}

	// 日志相关
	tmpStr, ok = conf.Get("common", "log_file")
	if ok {
		LogFile = tmpStr
		if LogFile == "console" {
			LogWay = "console"
		} else {
			LogWay = "file"
		}
	}

	tmpStr, ok = conf.Get("common", "log_level")
	if ok {
		LogLevel = tmpStr
	}
	return nil
}

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

func DoInit() {
	args, err := docopt.Parse(usage, nil, true, "V frpc_day6", false)

	if args["-c"] != nil {
		configFile = args["-c"].(string)
	}
	err = LoadConf(configFile)
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
		ServerAddr = addr[0]
		ServerPort = serverPort
	}

	if args["-L"] != nil {
		if args["-L"].(string) == "console" {
			LogWay = "console"
		} else {
			LogWay = "file"
			LogFile = args["-L"].(string)
		}
	}

	if args["--log-level"] != nil {
		LogLevel = args["--log-level"].(string)
	}
	// 0.1 初始化
	log.InitLog(LogWay, LogFile, LogLevel)
}