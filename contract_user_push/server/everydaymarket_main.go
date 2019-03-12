package main

import (
	"bbexgo/bbexutil"
	"bbexgo/common"
	"bbexgo/config"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"time"
)

var (
	r_mq_queue string
	w_mq_queue string

	user_push_msg_chan = make(chan common.PushSruct, 10)
)

func DonePush(cache *common.PushUserCache, c chan common.ServerQueueStruct,
	user_push_msg_chan *chan common.PushSruct) {
	for {
		select {
		case queue_obj := <-c:

			var obj common.QueueEveryDayMarketStruct
			err := json.Unmarshal([]byte(queue_obj.Msg), &obj)
			if err != nil {
				glog.Error("Unmarshal failed:", err)
				continue
			}

			queue_obj.Type = 1
			queue_obj.Exchange = 1
			now := uint64(time.Now().Unix())
			push_msg := common.PushSruct{
				ServerQueueStruct: queue_obj,
				Bpush:             true,
			}

			hour := time.Now().Hour()

			template := "\r\n【追踪币种】%s\r\n【指数价格】%.2f\r\n【24H波动】%+.2f%%\r\n【短线行情】波动超过1%%的行情: %d次\r\n【合约持仓】%.0f %s, 较昨日%+.2f%%\r\n【合约成交】%.0f %s, 较昨日%+.2f%%\r\n【多空分布】%.0f%%做多, %.0f%%做空\r\n\r\n关注更多数据，请前往订阅提醒 >"
			push_msg.Msg = fmt.Sprintf(template, obj.TraceSymbol, obj.IndexPrice, obj.Change24Houer, obj.ShortChange, obj.ContractHold, obj.TraceSymbol, obj.ContractHoldChange, obj.ContractTrade, obj.TraceSymbol, obj.ContractTradeChange, obj.ContractRatioMore, obj.COntractRatioEmpty)

			cache.UpdateCache()

			for _, event_cache := range cache.Cahce {
				user_map := event_cache.EventInfo
				for _, user_info := range user_map {
					if now < user_info.NextPushTime {
						glog.Info("user:", user_info.Id, " openid:", user_info.Openid, " NextPushTime:", user_info.NextPushTime, " not push")
						continue
					}
					push_msg.Bpush = true
					// 不推送 但是 存储
					if user_info.Night_push == "0" && hour < 8 {
						push_msg.Bpush = false
					}

					push_msg.Openid = user_info.Openid
					push_msg.UID = user_info.Id
					bbexutil.DayReportRemind(&push_msg)
					*user_push_msg_chan <- push_msg

					// doing
					glog.Info("send push 2 id:", user_info.Id, " openid:", user_info.Openid, " msg:", push_msg, " channle len:", cap(*user_push_msg_chan))
					user_info.LastPushTime = now
					user_info.NextPushTime = user_info.LastPushTime + 3600
				}
			}
		}
	}
	return
}

func MsgPushHandle(cache *common.PushUserCache) {
	glog.Info("MsgPushHandle start ....")
	done := make(chan struct{})
	send_mq_chan := make(chan struct{})

	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()
	cache.UpdateCache()
	msgchan := make(chan common.ServerQueueStruct, 10)

	event_chan := make(map[int]chan common.ServerQueueStruct)

	go bbexutil.WatchPushMsg(done, msgchan, r_mq_queue)
	go bbexutil.WritePushMsg(send_mq_chan, &user_push_msg_chan, w_mq_queue)
	for {
		select {
		case <-done:
			done = make(chan struct{})
			time.Sleep(1 * time.Second)
			go bbexutil.WatchPushMsg(done, msgchan, r_mq_queue)

		// 写push到Mq异常
		case <-send_mq_chan:
			time.Sleep(1 * time.Second)
			send_mq_chan = make(chan struct{})
			go bbexutil.WritePushMsg(send_mq_chan, &user_push_msg_chan, w_mq_queue)

		case msg := <-msgchan:
			c, ok := event_chan[msg.SubType]
			if !ok {
				msgchan := make(chan common.ServerQueueStruct, 10)
				event_chan[msg.SubType] = msgchan
				glog.Info("create a new chan for type:", msg.SubType)
				go DonePush(cache, msgchan, &user_push_msg_chan)
				msgchan <- msg
			} else {
				c <- msg
			}

		case <-ticker.C:
			glog.Info("MsgPushHandle live .....\n")
			break
		}
	}
}
func Init() {
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "./logs")
	flag.Set("log_backtrace_at", "1")
	flag.Set("max_log_size", "50")
	flag.Set("v", "3")
}

func main() {
	Init()
	flag.Parse()
	w_mq_queue = config.Get("mq.queue_user_push")
	r_mq_queue = config.Get("mq.queue_everydaymarket")
	var cache common.PushUserCache
	err := cache.Init("", -1)
	if err != nil {
		glog.Error("InitDB failed:", err)
		return
	}
	// 更新用户数据缓存
	go MsgPushHandle(&cache)
	for {
		time.Sleep(30 * time.Second)
	}
}
