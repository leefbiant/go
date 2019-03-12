package main

import (
	"bbexgo/bbexutil"
	"bbexgo/common"
	"bbexgo/wechat"
	"flag"
)

var (
	Openid = flag.String("Openid", "", "help message for flagname")
)

func main() {
	flag.Parse()
	data := common.PushSruct{}
	data.Symbol = "ETC"
	data.SubType = 2
	data.Exchange = 1
	data.ContractName = "ETC1228"
	data.Title = "ETC急涨提醒"
	if *Openid != "" {
		data.Msg = "【涨跌情况】上涨\r\n【当前价格】3421"
		bbexutil.MakePushContent(&data)
		api := wechat.Push{}
		api.SendUserSub(*Openid, 0, data.Pusg_msg, "http://contract.bbex.io/remind")
	}
}
