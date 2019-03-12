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
	"time"
)

var (
	r_mq_queue string
	w_mq_queue string

	user_push_msg_chan = make(chan common.PushSruct, 10)
)

func DonePush(cache *common.PushUserCache, c chan common.ServerQueueStruct,
	user_push_msg_chan *chan common.PushSruct) {
	push_interval := uint64(bbexutil.String2Int(config.Get("push_interval.sport"), 300))
	for {
		select {
		case queue_obj := <-c:

			var obj common.QueueSpotStruct
			err := json.Unmarshal([]byte(queue_obj.Msg), &obj)
			if err != nil {
				glog.Error("Unmarshal failed:", err)
				continue
			}

			Exchange := bbexutil.GetExchangeName(queue_obj.Exchange)
			if Exchange == "" {
				// log.Error("not find GetExchangeName for id:", queue_obj.Exchange)
				continue
			}
			keys := fmt.Sprintf("%d|%s|%s|%d", queue_obj.SubType, Exchange, queue_obj.Symbol, queue_obj.ContractType)
			user_map := cache.GetCache(keys)
			if user_map == nil {
				// log.Error("not find sub user for key:", keys, " queue_obj:", queue_obj)
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
			template := "【当前报价】%0.2f\r\n【提醒条件】%s%0.2f提醒"

			for _, user_info := range user_map.EventInfo {
				if now < user_info.NextPushTime {
					// log.Info("user:", user_info.Id, " openid:", user_info.Openid, " NextPushTime:", user_info.NextPushTime, " not push")
					continue
				}
				conditon := (1 == user_info.Val1 && obj.Price > user_info.Val2) || (2 == user_info.Val1 && obj.Price < user_info.Val2)
				if conditon {

					desc := "高于"
					if obj.Price < user_info.Val2 {
						desc = "低于"
					}
					push_msg.Msg = fmt.Sprintf(template, obj.Price, desc, user_info.Val2)

					push_msg.Bpush = true
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
					glog.Info("send push 2 id:", help.String2Int(user_info.Id), " openid:", user_info.Openid,
						" msg:", push_msg, " channle len:", cap(*user_push_msg_chan))
					user_info.LastPushTime = now
					user_info.NextPushTime = user_info.LastPushTime + push_interval

					err = cache.DisEnableSub(uint64(help.String2Int(user_info.Id)), queue_obj.SubType, Exchange, queue_obj.Symbol, queue_obj.ContractType)
					if err == nil {
						bbexutil.WritePushMsgOnce("sub_notify", "update")
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
	r_mq_queue = config.Get("mq.queue_spot")
	var cache common.PushUserCache
	err := cache.Init("sub_notify_spot", 1)
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
