package client

import (
	"strconv"

	"github.com/vaughan0/go-ini"
)

var (
	// 代理服务端地址信息
	ServerAddr string = "0.0.0.0"
	ServerPort int64  = 7000
	// 目标服务端地址信息
	LocalIp   string
	LocalPort int64
)

func LoadConf(confFile string) (err error) {
	var tmpStr string
	var ok bool

	conf, err := ini.LoadFile(confFile)
	if err != nil {
		return err
	}

	// frps 代理服务器地址
	tmpStr, ok = conf.Get("common", "server_addr")
	if ok {
		ServerAddr = tmpStr
	}

	tmpStr, ok = conf.Get("common", "server_port")
	if ok {
		ServerPort, _ = strconv.ParseInt(tmpStr, 10, 64)
	}

	// 目标服务器地址
	tmpStr, ok = conf.Get("app", "local_ip")
	if ok {
		LocalIp = tmpStr
	}

	tmpStr, ok = conf.Get("app", "local_port")
	if ok {
		LocalPort, _ = strconv.ParseInt(tmpStr, 10, 64)
	}

	return nil
}
