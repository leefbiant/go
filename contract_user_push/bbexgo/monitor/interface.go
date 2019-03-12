package monitor

import (
	"bbexgo/config"
	"encoding/json"
)

// 发送push消息结构
type PushData struct {
	UID          string `json:"uid"`          // 用户id
	Type         int    `json:"type"`         // 消息类型id
	SubType      int    `json:"subType"`      // 订阅类型
	Exchange     int    `json:"exchange"`     // 交易所ID
	Symbol       string `json:"symbol"`       // 合约币种
	ContractType int    `json:"contractType"` // 合约类型
	Time         int64  `json:"time"`         // 时间戳
	Msg          string `json:"msg"`          // 信息内容
}
