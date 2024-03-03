package main

import (
	"ch11/models/consts"
	"ch11/models/msg"
	"ch11/models/server"
	"ch11/utils/conn"
	"ch11/utils/log"
	"time"
)



func main() {
	// 初始化
	server.DoInit()
	// 供代理客户端使用的监听
	// proxyLn, err := net.Listen("tcp", ":7000")
	proxyLn, _ := conn.Listen(server.BindAddr, server.BindPort)//net.Listen("tcp", fmt.Sprintf("%s:%d", server.BindAddr, server.BindPort))
	defer proxyLn.Close()

	log.Info("服务器正在监听 %d 端口...\n", server.BindPort)

	workConnChan := make(chan *conn.Conn)
	// 接受传入的连接
	go func() {
		for {
			ctrlConn, _ := proxyLn.GetConn() 
			// 读取请求信息
			cliReq := &msg.ControlReq{}
			ctrlConn.ReadObj(cliReq)
			log.Debug("input request %s", cliReq.Type)
			// 检查请求类型
			switch cliReq.Type {
			case consts.NewCtlConn:
				log.Debug("New Control connection.")
				// 如果是工作连接，则返回响应码 0，表示连接成功
				res := &msg.ControlRes{
					Code: 0,
				}
				ctrlConn.WriteObj(res)
				// 启动心跳功能
				go func() {
					timer := time.AfterFunc(time.Duration(server.HeartBeatTimeout) * time.Second, func() {
						log.Error("client heartbeat timeout")
						ctrlConn.Close()
					})
					defer timer.Stop()
					cliReq := &msg.ControlReq{}
					ctrlConn.ReadObj(cliReq)
					if cliReq.Type == consts.HeartbeatReq {
						log.Debug("get heartbeat")
						timer.Reset(time.Duration(server.HeartBeatTimeout) * time.Second)
						heartbeatRes := &msg.ControlRes{
							Type: consts.HeartbeatRes,
						}
						ctrlConn.WriteObj(heartbeatRes)
					}
				}()
			case consts.NewWorkConn:
				go func() { workConnChan <- ctrlConn }()
			default:
				log.Warn(" unsupport msgType [%d]", cliReq.Type)
			}
		}
	}()


	// 监听供用户的传入连接
	ln, _ := conn.Listen(server.ListenAddr, server.ListenPort) //net.Listen("tcp", fmt.Sprintf("%s:%d", server.ListenAddr, server.ListenPort))
	defer ln.Close()

	log.Info("服务器正在监听 %s:%d ...\n", server.ListenAddr,server.ListenPort)

	// 接受传入的连接
	userConn, _ := ln.GetConn()
	workConn, _ := <- workConnChan

	conn.Join(userConn, workConn)

	log.Warn("All proxy exit!")
}
