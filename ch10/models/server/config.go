package server

import (
	"strconv"

	"github.com/vaughan0/go-ini"
)

var (
	// 供代理客户端连接
	BindAddr string = "0.0.0.0"
	BindPort int64  = 7000
	// 日志参数
	LogFile           string = "console"
	LogWay            string = "console"
	LogLevel          string = "info"
	// 供外部用户连接
	ListenAddr string
	ListenPort int64
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