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
	"math"
	"time"
)

var (
	r_mq_queue string
	w_mq_queue string

	user_push_msg_chan = make(chan common.PushSruct, 10)
)

func DonePush(cache *common.PushUserCache, c chan common.ServerQueueStruct,
	user_push_msg_chan *chan common.PushSruct) {
	push_interval := uint64(bbexutil.String2Int(config.Get("push_interval.blastingorder"), 300))
	for {
		select {
		case queue_obj := <-c:

			var obj common.QueueBlastingOrderStruct
			err := json.Unmarshal([]byte(queue_obj.Msg), &obj)
			if err != nil {
				glog.Error("Unmarshal failed:", err)
				continue
			}

			Exchange := bbexutil.GetExchangeName(queue_obj.Exchange)
			if Exchange == "" {
				glog.Error("not find GetExchangeName for id:", queue_obj.Exchange)
				continue
			}
			keys := fmt.Sprintf("%d|%s|%s|%d", queue_obj.SubType, Exchange, queue_obj.Symbol, queue_obj.ContractType)
			now := uint64(time.Now().Unix())
			push_msg := common.PushSruct{
				ServerQueueStruct: queue_obj,
				Bpush:             true,
			}

			IsDelay := false
			if time.Now().Unix()-queue_obj.Time > 60*30 {
				IsDelay = true
			}

			hour := time.Now().Hour()
			trade_type := fmt.Sprintf("%s", "空单强平")
			if "buy" == obj.Type {
				trade_type = fmt.Sprintf("%s", "多单强平")
			}

			template := "【交易类型】%s\r\n【委托价格】%.2f\r\n【爆仓数量】%d张\r\n【折算价值】%.2f %s"
			push_msg.Msg = fmt.Sprintf(template, trade_type, obj.Price, obj.Size, obj.AmountSymbol, queue_obj.Symbol)

			// sync bbx if bbx notify
			if bbexutil.GetExchangeName(queue_obj.Exchange) == "BBX" {
				nodify_obj := common.BBXNotify{
					Symbol:       queue_obj.Symbol,
					ContractName: queue_obj.ContractName,
					Type:         queue_obj.SubType,
					Time:         queue_obj.Time,
					Msg:          push_msg.Msg,
				}
				bbexutil.SyncBBx(nodify_obj)
			}

			user_map := cache.GetCache(keys)
			if user_map == nil {
				continue
			}

			for _, user_info := range user_map.EventInfo {
				if now < user_info.NextPushTime {
					glog.Info("user:", user_info.Id, " openid:", user_info.Openid, " NextPushTime:", user_info.NextPushTime, " not push")
					continue
				}
				if math.Abs(obj.AmountUSD) > user_info.Val1*10000 {
					push_msg.Bpush = true
					// 不推送 但是 存储
					if user_info.Night_push == "0" && hour < 8 || IsDelay {
						push_msg.Bpush = false
					}

					push_msg.Openid = user_info.Openid
					push_msg.UID = user_info.Id
					bbexutil.MakePushContent(&push_msg)
					*user_push_msg_chan <- push_msg

					// doing
					glog.Info("send push 2 id:", user_info.Id, " openid:", user_info.Openid, " msg:", push_msg, " channle len:", cap(*user_push_msg_chan))
					if push_msg.Bpush {
						user_info.LastPushTime = now
						user_info.NextPushTime = user_info.LastPushTime + push_interval
						glog.Info("user:", user_info.Openid, " key:", keys, " send push next time:", user_info.NextPushTime)
					}
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
			cache.UpdateCache()
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
	r_mq_queue = config.Get("mq.queue_blastingorder")
	var cache common.PushUserCache
	err := cache.Init("sub_notify_blastingorder", 6)
	if err != nil {
		glog.Error("InitDB failed:", err)
		return
	}
	// 更新用户数据缓存
	go common.UpdateCacheHandle(&cache)
	go MsgPushHandle(&cache)
	for {
		time.Sleep(30 * time.Second)
	}
}
