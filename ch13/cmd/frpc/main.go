package main

import (
	"ch13/models/client"
	"ch13/models/consts"
	"ch13/models/msg"
	"ch13/utils/conn"
	"ch13/utils/log"
	"time"
)

func main() {
	// 0. 加载参数
	client.DoInit()
	// 连接外网代理服务器作为控制连接
	for _, client := range client.ProxyClients {
		loginToServer(client)
	}
	

	log.Warn("All proxy exit!")
}

func loginToServer(cli *client.ProxyClient) {
	// 控制连接，每个客户端代理一个
	ctrlConn, _ := conn.ConnectServer(client.ServerAddr, client.ServerPort)
	// 请求为控制连接
	req := &msg.ControlReq{
		Type: consts.NewCtlConn,
		ProxyName: cli.Name, // 使用名称表示对应的服务
	}
	ctrlConn.WriteObj(req)

	// 响应代码为0，则连接成功
	ctlRes := &msg.ControlRes{}
	ctrlConn.ReadObj(ctlRes)

	if ctlRes.Code != 0 {
		log.Error("start proxy error, %s", ctlRes.Msg)
	}

	log.Debug("connect to server [%s:%d] success!", client.ServerAddr, client.ServerPort)

	// 发送心跳
	go func() {
		for {
			time.Sleep(time.Duration(client.HeartBeatInterval) * time.Second)
			if ctrlConn != nil {
				log.Debug("Send heartbeat to server")
				heartbeatReq := &msg.ControlReq{
					Type: consts.HeartbeatReq,
				}
				ctrlConn.WriteObj(heartbeatReq)

			} else {
				break
			}
		}
	}()
	timer := time.AfterFunc(time.Duration(client.HeartBeatTimeout)*time.Second, func() {
		log.Error("heartbeatRes from frps timeout")
		ctrlConn.Close()
	})
	defer timer.Stop()
	for {
		ctlRes := &msg.ControlRes{}
		ctrlConn.ReadObj(ctlRes)
		if ctlRes.Type == consts.HeartbeatRes {
			log.Debug("receive heartbeat response")
			timer.Reset(time.Duration(client.HeartBeatTimeout) * time.Second)
		} else if ctlRes.Type == consts.NoticeUserConn {
			cli.StartTunnel(client.ServerAddr, client.ServerPort)
		}
	}
}
