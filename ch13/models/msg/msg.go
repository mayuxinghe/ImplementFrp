package msg

type GeneralRes struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

// messages between control connection of frpc and frps
type ControlReq struct {
	Type          int64  `json:"type"`
	ProxyName     string `json:"proxy_name,omitempty"`
	AuthKey       string `json:"auth_key, omitempty"`
	UseEncryption bool   `json:"use_encryption, omitempty"`
	Timestamp     int64  `json:"timestamp, omitempty"`
}

type ControlRes struct {
	Type int64  `json:"type"`
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}