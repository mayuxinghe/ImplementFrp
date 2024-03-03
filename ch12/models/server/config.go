package server

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
	// 供代理客户端连接
	BindAddr string = "0.0.0.0"
	BindPort int64  = 7000
	// 日志参数
	LogFile  string = "console"
	LogWay   string = "console"
	LogLevel string = "info"
	// 供外部用户连接
	ListenAddr string
	ListenPort int64
	// 心跳超时常量
	HeartBeatTimeout int64 = 90
	UserConnTimeout  int64 = 10
)

/**
* 加载指定的配置文件
 */
func LoadConf(confFile string) (err error) {
	var tmpStr string
	var ok bool

	// 解析配置文件
	conf, err := ini.LoadFile(confFile)
	if err != nil {
		return err
	}

	// common secction
	tmpStr, ok = conf.Get("common", "bind_addr")
	if ok {
		BindAddr = tmpStr
	}

	tmpStr, ok = conf.Get("common", "bind_port")
	if ok {
		BindPort, _ = strconv.ParseInt(tmpStr, 10, 64)
	}

	// app section
	tmpStr, ok = conf.Get("app", "listen_addr")
	if ok {
		ListenAddr = tmpStr
	}

	tmpStr, ok = conf.Get("app", "listen_port")
	if ok {
		ListenPort, _ = strconv.ParseInt(tmpStr, 10, 64)
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

var configFile = "./frps.ini"
var usage string = `代理服务端程序的使用方法

Usage: 
	frps [-c config_file] [-L log_file] [--log-level=<log_level>] [--addr=<bind_addr>]
	frps -h | --help | --version

Options:
	-c config_file            设置配置文件路径
	--addr=<bind_addr>        代理监听地址，形式如: 0.0.0.0:7000
	-L log_file               指定日志输出文件, including console
    --log-level=<log_level>   指定日志等级: debug, info, warn, error
	-h --help                 显示本帮助信息
	--version                 显示版本号
`

func DoInit() {
	// 加载参数
	// 解析命令行参数
	args, err := docopt.Parse(usage, nil, true, "v frps_day6", false)
	// 根据 -c 选项获取命令行指定配置文件
	if args["-c"] != nil {
		configFile = args["-c"].(string)
	}
	err = LoadConf(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
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
		BindAddr = addr[0]
		BindPort = bindPort
	}
	// 0.1 初始化
	log.InitLog(LogWay, LogFile, LogLevel)
}
