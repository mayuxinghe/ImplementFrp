package conn

import (
	"bufio"
	"ch11/utils/log"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
)
type Listener struct {
	l         *net.TCPListener
	conns     chan *Conn
}

func Listen(bindAddr string, bindPort int64) (l *Listener, err error) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", bindAddr, bindPort))
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return l, err
	}

	l = &Listener{
		l:         listener,
		conns:     make(chan *Conn),
	}

	go func() {
		for {
			conn, err := l.l.AcceptTCP()
			if err != nil {
				continue
			}

			c := &Conn{
				TcpConn:   conn,
			}
			l.conns <- c
		}
	}()
	return l, err
}

// wait util get one new connection or listener is closed
// if listener is closed, err returned
func (l *Listener) GetConn() (conn *Conn, err error) {
	var ok bool
	conn, ok = <-l.conns
	if !ok {
		return conn, fmt.Errorf("channel close")
	}
	return conn, nil
}

func (l *Listener) Close()  {
	l.l.Close()
}

type Conn struct {
	TcpConn *net.TCPConn
}

func ConnectServer(host string, port int64) (c *Conn, err error) {
	c = &Conn{}
	servertAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return
	}
	conn, err := net.DialTCP("tcp", nil, servertAddr)
	if err != nil {
		return
	}
	c.TcpConn = conn
	return c, nil
}

func (c *Conn) Close() {
	c.TcpConn.Close()
}

func (c * Conn) ReadStr() (str string, err error) {
	str, err = bufio.NewReader(c.TcpConn).ReadString('\n');
	return str, err
}

func (c *Conn) WriteStr(str string) (i int, err error) {
	i, err = c.TcpConn.Write([]byte(str + "\n"))
	return i, err
}

func (c *Conn) WriteObj(v any)(i int, err error) {
	buf, err := json.Marshal(v)
	if (err != nil) {
		log.Error("对象序列化失败")
		return 0, err
	}
	i, err = c.WriteStr(string(buf))
	return i, err
}

func (c *Conn) ReadObj(v any)(err error) {
	res, err := c.ReadStr()
	if (err != nil) {
		return err
	}
	err = json.Unmarshal([]byte(res), &v)
	return err
}

func Join(c1 *Conn, c2 *Conn) {
	var wait sync.WaitGroup
	pipe := func(to *Conn, from *Conn) {
		defer to.Close()
		defer from.Close()
		defer wait.Done()

		var err error
		_, err = io.Copy(to.TcpConn, from.TcpConn)
		if err != nil {
			log.Warn("join conns error, %v", err)
		}
	}

	wait.Add(2)
	go pipe(c1, c2)
	go pipe(c2, c1)
	wait.Wait()
}