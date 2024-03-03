package consts

// server status
const (
	Idle = iota
	Working
	Closed
)

// msg type
const (
	NewCtlConn = iota
	NewWorkConn
	NoticeUserConn
	NewCtlConnRes
	HeartbeatReq
	HeartbeatRes
)
