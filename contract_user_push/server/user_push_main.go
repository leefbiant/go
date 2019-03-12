package main

import (
	"bbexgo/common"
	"bbexgo/config"
	"bbexgo/redis"
	"bbexgo/wechat"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/streadway/amqp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PushStat struct {
	cnt            int
	last_push_time int64
}

var (
	ExchangeList              []string // 交易所列表
	user_push_msg_chan        = make(chan common.PushData, 10)
	user_push_queue           string
	OneDB                     = initDB()
	TableName_UserPushHistory = "user_push_history"
	push_stat                 = make(map[int]*PushStat)
	mutex                     = &sync.Mutex{}
	curr_date                 = ""
	last_data                 = ""
	statis_times              = 0
	redisClient               = redis.GetInstance()
)

func AddPushStat(subtype int) {
	mutex.Lock()
	defer mutex.Unlock()
	obj, ok := push_stat[subtype]
	if !ok {
		push_stat[subtype] = &PushStat{
			cnt:            1,
			last_push_time: time.Now().Unix(),
		}
		return
	}
	(*obj).cnt += 1
	(*obj).last_push_time = time.Now().Unix()
}

func PushDataStat() {
	glog.Info("PushDataStat")
	mutex.Lock()
	defer mutex.Unlock()
	for k, v := range push_stat {
		glog.Info("type:", k, " last_push_time:", time.Unix((*v).last_push_time, 0).Format("2006-01-02 15:04:05"), " cnt:", (*v).cnt)
	}
	curr_t := time.Now().Unix()
	curr_date = time.Unix(curr_t, 0).Format("2006-01-02")
	if last_data == "" {
		index := 0
		for index < 8 {
			sql := fmt.Sprintf("insert into push_stat(sub_type, push_num, date, updated_at, created_at) values(%d, 0, '%s', now(), now())", index, curr_date)
			index += 1
			_, err := OneDB.Exec(sql)
			if err != nil {
				glog.Error(err)
				glog.Error(sql)
			}
		}
		last_data = curr_date
		return
	}
	statis_times += 1

	// 跨天 先写入新的
	if curr_date != last_data {
		index := 0
		for index < 8 {
			sql := fmt.Sprintf("insert into push_stat(sub_type, push_num, date, updated_at, created_at) values(%d, 0, '%s', now(), now())", index, curr_date)
			index += 1
			_, err := OneDB.Exec(sql)
			if err != nil {
				glog.Error(err)
				glog.Error(sql)
			}
		}
		for k, v := range push_stat {
			value := (*v).cnt
			if value == 0 {
				continue
			}
			(*v).cnt = 0
			sql := fmt.Sprintf("update push_stat set push_num = push_num + %d where sub_type = %d and date = '%s'", value, k, last_data)
			_, err := OneDB.Exec(sql)
			if err != nil {
				glog.Error(err)
				glog.Error(sql)
			}
		}
		last_data = curr_date
	}

	// 每5分钟更新数据
	if statis_times%10 == 0 {
		for k, v := range push_stat {
			value := (*v).cnt
			if value == 0 {
				continue
			}
			(*v).cnt = 0
			sql := fmt.Sprintf("update push_stat set push_num = push_num + %d where sub_type = %d and date = '%s'", value, k, curr_date)
			glog.Info("sql:", sql)
			_, err := OneDB.Exec(sql)
			if err != nil {
				glog.Error(err)
				glog.Error(sql)
			}
		}
	}
}

func PushStatThread() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			PushDataStat()
		}
	}
}

func initDB() *sql.DB {
	constr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", config.Get("mysql.default.username"),
		config.Get("mysql.default.password"), config.Get("mysql.default.host"), config.Get("mysql.default.port"),
		config.Get("mysql.default.database"), config.Get("mysql.default.charset"))
	glog.Info(constr)
	db, err := sql.Open("mysql", constr)
	if err != nil {
		glog.Fatal(err)
	}
	return db
}

func save(sqlStrPrefix string, insertData *[]string) {
	batchInsert := *insertData
	if len(batchInsert) == 0 {
		return
	}
	*insertData = batchInsert[:0:0]
	insertStr := fmt.Sprintf("%s ( %s )", sqlStrPrefix, strings.Join(batchInsert, " ), ("))
	glog.Info("save:", insertStr)
	_, err := OneDB.Exec(insertStr)
	if err != nil {
		glog.Error(err)
		glog.Error(insertStr)
	}
}

func saveMsg(data common.PushSruct) {
	fields := []string{
		"uid", "type", "sub_type",
		"exchange", "symbol", "contract_type", "send_time", "msg",
		"createtime", "wechat_msg_id", "status", "err"}
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", TableName_UserPushHistory, strings.Join(fields, " , "))
	batchInsert := make([]string, 0, 1) // 批量插入

	tmp := make([]string, 0, len(fields))
	for _, filed := range fields {
		var value string
		switch filed {
		case "uid":
			value = data.UID
			break
		case "type":
			value = strconv.Itoa(data.Type)
			break
		case "sub_type":
			value = strconv.Itoa(data.SubType)
			break
		case "exchange":
			value = strconv.Itoa(data.Exchange)
			break
		case "symbol":
			value = strings.ToUpper(data.Symbol)
			break
		case "contract_type":
			value = strconv.Itoa(data.ContractType)
			break
		case "send_time":
			value = strconv.FormatInt(data.Time, 10)
			break
		case "msg":
			value = data.Msg
			break
		case "createtime":
			value = strconv.FormatInt(time.Now().Unix(), 10)
			break
		case "wechat_msg_id":
			value = data.WechatMsgId
			break
		case "status":
			value = strconv.FormatInt(data.PushStatusCode, 10)
			break
		case "err":
			value = data.PushErrorMsg
			break
		default:
			value = ""
			break
		}
		tmp = append(tmp, value)
	}
	batchInsert = append(batchInsert, fmt.Sprintf("'%s'", strings.Join(tmp, "' , '")))
	glog.Info(tmp)
	save(sqlStr, &batchInsert)
}

func DonePushChannle(threadid int, done chan struct{}) {
	defer close(done)
	mq_addr := config.Get("mq.addr")
	conn, err := amqp.Dial(mq_addr)
	if err != nil {
		glog.Info("amqp failed:", err)
		return
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		glog.Info("Channel failed:%v", err)
		return
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		user_push_queue, // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		glog.Error("QueueDeclare failed:%v", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		glog.Error("Consume failed:%v", err)
		return
	}
	bpush := config.Get("user_push.bpush")
	for d := range msgs {
		var data common.PushSruct
		err := json.Unmarshal([]byte(d.Body), &data)
		if err != nil {
			glog.Error("Unmarshal failed:", err)
			return
		}

		if data.Bpush {
			AddPushStat(data.SubType)
			glog.Info("thread id", threadid, " recv a user push openid:", data.Openid, " msg:", data)
			api := wechat.Push{}
			if "0" == bpush {
				glog.Info("bpush == 0 not send nitify and save msg")
				continue
			}
			jump_url := data.Url
			if jump_url == "" {
				jump_url = "http://contract.bbex.io/remind"
			}
			api.SendUserSub(data.Openid, data.Type, data.Pusg_msg, jump_url)

			glog.Info("thread id", threadid, " msg :", api)
			if api.Success {
				data.WechatMsgId = strconv.FormatInt(api.MsgId, 10)
				data.PushStatusCode = api.ErrorCode
				data.PushErrorMsg = api.ErrMsg
			} else if api.ErrorCode == 43004 {
				glog.Error("find a user not sub :", data.Openid)
				sql := fmt.Sprintf("update user_subs set status = 0 where status = 1 and uid in	(select id from users where openid = '%s') limit 100", data.Openid)
				glog.Info(sql)
				_, err := OneDB.Exec(sql)
				if err != nil {
					glog.Error(err)
				}
				continue
			}
		}
		saveMsg(data)
	}
}

func Handlechannle(threadid int) {
	glog.Info("Handlechannle start ....")
	c := make(chan struct{})
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()

	go DonePushChannle(threadid, c)
	for {
		select {
		case <-c:
			time.Sleep(1 * time.Second)
			c = make(chan struct{})
			go DonePushChannle(threadid, c)

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
	user_push_queue = config.Get("mq.queue_user_push")

	// 多协程发送
	thread_num := 5
	for i := 0; i < thread_num; i++ {
		go Handlechannle(i)
	}

	go PushStatThread()
	for {
		time.Sleep(30 * time.Second)
	}
}
