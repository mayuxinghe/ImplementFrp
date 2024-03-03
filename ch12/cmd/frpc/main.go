package main

import (
	"ch11/models/client"
	"ch11/models/consts"
	"ch11/models/msg"
	"ch11/utils/conn"
	"ch11/utils/log"
	"time"
)

func main() {
	// 0. 加载参数
	client.DoInit()
	// 连接外网代理服务器作为控制连接
	ctlConn, _ := conn.ConnectServer(client.ServerAddr, client.ServerPort)
	loginToServer(ctlConn)

	log.Warn("All proxy exit!")
}

func loginToServer(ctrlConn *conn.Conn) {

	// 请求为控制连接
	req := &msg.ControlReq{
		Type: consts.NewCtlConn,
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
			// 连接到一层代理服务器作为工作连接
			workConn, _ := conn.ConnectServer(client.LocalIp, client.LocalPort)
			// 请求为工作连接
			req := &msg.ControlReq{
				Type: consts.NewWorkConn,
			}
			workConn.WriteObj(req)
			// 连接内网应用服务器
			targetCoon, _ := conn.ConnectServer(client.LocalIp, client.LocalPort)
			// 3. 建立一层代理和目标服务器之间的隧道通讯 =================
			conn.Join(workConn, targetCoon)
		}
	}
}
