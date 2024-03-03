// Copyright 2016 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"ch13/models/consts"
	"ch13/models/msg"
	"ch13/utils/conn"
	"ch13/utils/log"
)

type ProxyClient struct {
	Name      string
	LocalIp   string
	LocalPort int64
}

/**
* 和目标服务器建立连接
*/
func (p *ProxyClient) GetLocalConn() (c *conn.Conn, err error) {
	c, err = conn.ConnectServer(p.LocalIp, p.LocalPort)
	if err != nil {
		log.Error("ProxyName [%s], connect to local port error, %v", p.Name, err)
	}
	return
}
/**
* 建立工作连接
*/
func (p *ProxyClient) GetRemoteConn(addr string, port int64) (c *conn.Conn, err error) {
	defer func() {
		if err != nil {
			c.Close()
		}
	}()

	c, err = conn.ConnectServer(addr, port)

	req := &msg.ControlReq{
		Type:      consts.NewWorkConn,
		ProxyName: p.Name,
	}

	c.WriteObj(req)

	err = nil
	return
}

func (p *ProxyClient) StartTunnel(serverAddr string, serverPort int64) (err error) {
	localConn, err := p.GetLocalConn()
	remoteConn, err := p.GetRemoteConn(serverAddr, serverPort)
	go conn.Join(localConn, remoteConn)
	return nil
}
