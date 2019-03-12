package main

import (
	"bbexgo/bbexutil"
	"bbexgo/common"
	"bbexgo/config"
	"bbexgo/help"
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
	push_interval := uint64(bbexutil.String2Int(config.Get("push_interval.premium"), 43200))
	for {
		select {
		case queue_obj := <-c:

			var obj common.QueuePremiumStruct
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
			user_map := cache.GetCache(keys)
			if user_map == nil {
				// log.Error("not find sub user for key:", keys)
				continue
			}
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
			side := fmt.Sprintf("正溢价%.2f%%", math.Abs(obj.Premium))
			if obj.Premium < 0 {
				side = fmt.Sprintf("负溢价%.2f%%", math.Abs(obj.Premium))
			}
			template := "【溢价情况】%s\r\n【合约价格】%.2f\r\n【现货指数】%.2f"

			push_msg.Msg = fmt.Sprintf(template, side, obj.FuturePrice, obj.IndexPrice)

			for _, user_info := range user_map.EventInfo {
				if now < user_info.NextPushTime {
					// glog.Info("user:", user_info.Id, " openid:", user_info.Openid, " NextPushTime:", user_info.NextPushTime, " not push")
					continue
				}
				if math.Abs(obj.Premium) > user_info.Val1 {

					// 不推送 但是 存储
					if user_info.Night_push == "0" && hour < 8 || IsDelay {
						push_msg.Bpush = false
					}

					push_msg.Openid = user_info.Openid
					push_msg.UID = user_info.Id
					push_msg.Url = bbexutil.GetContractUrl(queue_obj)
					bbexutil.MakePushContent(&push_msg)
					*user_push_msg_chan <- push_msg

					// doing
					glog.Info("send push 2 id:", user_info.Id, " openid:", user_info.Openid, " msg:", push_msg, " channle len:", cap(*user_push_msg_chan))
					user_info.LastPushTime = now
					user_info.NextPushTime = user_info.LastPushTime + push_interval

					cache.UpdateUserNextNotify(uint64(help.String2Int(user_info.Id)),
						queue_obj.SubType, Exchange, queue_obj.Symbol, queue_obj.ContractType, user_info.NextPushTime)
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
	r_mq_queue = config.Get("mq.queue_premium")
	var cache common.PushUserCache
	err := cache.Init("sub_notify_premium", 3)
	if err != nil {
		glog.Error("InitDB failed:%v", err)
		return
	}
	// 更新用户数据缓存
	go common.UpdateCacheHandle(&cache)
	go MsgPushHandle(&cache)
	for {
		time.Sleep(30 * time.Second)
	}
}
