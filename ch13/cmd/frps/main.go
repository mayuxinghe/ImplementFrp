package main

import (
	"ch13/models/consts"
	"ch13/models/msg"
	"ch13/models/server"
	"ch13/utils/conn"
	"ch13/utils/log"
	"time"
)

func main() {
	// 初始化
	server.DoInit()
	// 供代理客户端使用的监听
	// proxyLn, err := net.Listen("tcp", ":7000")
	proxyLn, _ := conn.Listen(server.BindAddr, server.BindPort) //net.Listen("tcp", fmt.Sprintf("%s:%d", server.BindAddr, server.BindPort))
	defer proxyLn.Close()

	log.Info("服务器正在监听 %d 端口...\n", server.BindPort)

	// 接受传入的连接
	for {
		newConn, _ := proxyLn.GetConn()
		go func() {
			// 读取请求信息
			cliReq := &msg.ControlReq{}
			newConn.ReadObj(cliReq)
			log.Debug("input request %s", cliReq.Type)
			s, _ := server.ProxyServers[cliReq.ProxyName]
			// 检查请求类型
			switch cliReq.Type {
			case consts.NewCtlConn:
				log.Debug("New Control connection.")
				// 如果是工作连接，则返回响应码 0，表示连接成功
				res := &msg.ControlRes{
					Code: 0,
				}
				newConn.WriteObj(res)
				// 启动心跳功能
				go func() {
					timer := time.AfterFunc(time.Duration(server.HeartBeatTimeout)*time.Second, func() {
						log.Error("client heartbeat timeout")
						newConn.Close()
					})
					defer timer.Stop()
					cliReq := &msg.ControlReq{}
					newConn.ReadObj(cliReq)
					if cliReq.Type == consts.HeartbeatReq {
						log.Debug("get heartbeat")
						timer.Reset(time.Duration(server.HeartBeatTimeout) * time.Second)
						heartbeatRes := &msg.ControlRes{
							Type: consts.HeartbeatRes,
						}
						newConn.WriteObj(heartbeatRes)
					}
				}()
				notice := &msg.ControlRes{
					Type: consts.NoticeUserConn,
				}
				for {
					closeFlag := s.WaitUserConn()
					if closeFlag {
						log.Debug("ProxyName [%s], goroutine for dealing user conn is closed", s.Name)
						break
					}
					newConn.WriteObj(notice)
				}
			case consts.NewWorkConn:
				s.GetNewCliConn(newConn) //go func() { workConnChan <- ctrlConn }()
			default:
				log.Warn(" unsupport msgType [%d]", cliReq.Type)
			}


		}()
	}

}
