package bbexutil

import (
	"bbexgo/common"
	"bbexgo/config"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ContractAttr struct {
	// symbol       string // 火币名称
	// exchange     string // 交易所
	// contractType string // 合约类型
	Attr map[string]interface{}
}

var (
	SubTypeNameList []string // 订阅类型列表
	ExchangeList    []string // 交易所列表
	// redisClient     = redis.GetInstance()
)

func GetSubTypeName(subType int) string {
	if subType <= 0 {
		return ""
	}
	if len(SubTypeNameList) != 0 && subType > 0 && subType <= len(SubTypeNameList) {
		return SubTypeNameList[subType-1]
	}
	if err := json.Unmarshal([]byte(config.Get("SubType", "elf")), &SubTypeNameList); err != nil || len(SubTypeNameList) == 0 {
		glog.Fatal(err, len(SubTypeNameList))
		return ""
	}
	return SubTypeNameList[subType-1]
}

func GetExchangeName(exchangeID int) string {
	if exchangeID <= 0 {
		glog.Error("err exchangeID:", exchangeID)
		return ""
	}
	if len(ExchangeList) != 0 && exchangeID > 0 && exchangeID <= len(ExchangeList) {
		return ExchangeList[exchangeID-1]
	}
	if err := json.Unmarshal([]byte(config.Get("ExchangeList", "elf")), &ExchangeList); err != nil || len(ExchangeList) == 0 {
		glog.Fatal(err, len(ExchangeList))
		return ""
	}
	if exchangeID > len(ExchangeList) {
		glog.Error("err exchangeID:", exchangeID)
		return ""
	}
	return ExchangeList[exchangeID-1]
}

func MakePushContent(data *common.PushSruct) {
	title := fmt.Sprintf("%s%s\n", strings.ToUpper(data.Symbol), GetSubTypeName(data.SubType))
	if data.Title != "" {
		title = data.Title
	}
	data.Pusg_msg = map[string]map[string]string{
		"first": {
			"value": fmt.Sprintf("%s\n%s\n", title, data.Msg),
			"color": "#f5a100",
		},
		"keyword1": {
			"value": fmt.Sprintf("%s的%s", GetExchangeName(data.Exchange), data.ContractName),
			"color": "#767680",
		},
		"keyword2": {
			"value": GetSubTypeName(data.SubType),
			"color": "#767680",
		},
		"keyword3": {
			"value": time.Unix(data.Time, 0).Format("2006-01-02 15:04"),
			"color": "#767680",
		},
		"remark": {
			"value": "\n本提醒由【合约精灵】服务号提供，点击[提醒设置]定制个性化推送规则>>",
			"color": "#767680",
		},
	}
}

func DayReportRemind(data *common.PushSruct) {
	data.Pusg_msg = map[string]map[string]string{
		"first": {
			"value": "今日份行情总结已生成，请查收。",
			"color": "#ffba45",
		},
		"keyword1": {
			"value": fmt.Sprintf("今日行情总结 (%s合约)", GetExchangeName(data.Exchange)),
			"color": "#767680",
		},
		"keyword2": {
			"value": time.Unix(data.Time, 0).Format("2006年01月02日 15:04"),
			"color": "#767680",
		},
		"remark": {
			"value": data.Msg,
			"color": "#f5a100",
		},
	}
}

func WatchPushMsg(done chan struct{}, msgchan chan common.ServerQueueStruct, queue string) {
	defer close(done)
	mq_addr := config.Get("mq.addr")
	conn, err := amqp.Dial(mq_addr)
	if err != nil {
		glog.Info("amqp failed:%v", err)
		return
	}

	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		glog.Info("Channel failed:%v", err)
		return
	}

	glog.Info("init watch mq addr:", mq_addr, " queue:", queue)
	defer ch.Close()
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
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
	for d := range msgs {
		glog.Info("Received a message:", string(d.Body))
		var msg common.ServerQueueStruct
		err := json.Unmarshal([]byte(d.Body), &msg)
		if err != nil {
			glog.Error("Unmarshal failed:", err)
			return
		}
		msgchan <- msg
	}
}

func WatchMq(done chan struct{}, msgchan chan string, queue string) {
	defer close(done)
	mq_addr := config.Get("mq.addr")
	conn, err := amqp.Dial(mq_addr)
	if err != nil {
		glog.Info("amqp failed:%v", err)
		return
	}

	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		glog.Info("Channel failed:%v", err)
		return
	}

	glog.Info("init watch mq addr:", mq_addr, " queue:", queue)
	defer ch.Close()
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
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
	for d := range msgs {
		msgchan <- string(d.Body)
	}
}

func WritePushMsg(c chan struct{}, user_push_msg_chan *chan common.PushSruct, w_mq_queue string) {
	defer close(c)
	conn, err := amqp.Dial(config.Get("mq.addr"))
	if err != nil {
		glog.Error("amqp.Dial failed:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		glog.Error("conn.Channel failed:", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		w_mq_queue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		glog.Error("QueueDeclare.Dial failed:", err)
		return
	}
	glog.Info("WritePushMsg queue:", w_mq_queue)

	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()

	t_now := time.Now().Unix()
	t_last_publish := t_now

	for {
		select {
		case push_obj := <-*user_push_msg_chan:
			glog.Info("publish obj:", push_obj)
			jsonBytes, _ := json.Marshal(push_obj)

			body := string(jsonBytes)
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			if err != nil {
				glog.Error("err Publish:", err)
				return
			}
			t_last_publish = time.Now().Unix()
		case <-ticker.C:
			if time.Now().Unix()-t_last_publish > 300 {
				glog.Error("mq Publish timeout")
				return
			}
		}
	}
}

func WritePushMsgOnce(queue_name string, data string) {
	glog.Info("start push queue:", queue_name, " data:", data)
	conn, err := amqp.Dial(config.Get("mq.addr"))
	if err != nil {
		glog.Error("amqp.Dial failed:", err)
		return
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		glog.Error("conn.Channel failed:", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"sub_notify", // exchange
		"fanout",     // kind
		false,        // delete when unused
		false,        // exclusive
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		glog.Error("QueueDeclare.Dial failed:", err)
		return
	}

	err = ch.Publish(
		"sub_notify", // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
}

func BBxNotify(queue_name string, data string) {
	mq_addr := config.Get("mq.addr")
	conn, err := amqp.Dial(mq_addr)
	if err != nil {
		glog.Error("amqp.Dial failed:", err)
		return
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		glog.Error("conn.Channel failed:", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"bbx_notify", // name
		"direct",     // durable
		true,         // delete when unused
		false,        // exclusive
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		glog.Error("ExchangeDeclare failed:", err, " mq_addr:", mq_addr, " queue:", queue_name)
		return
	}

	_, err = ch.QueueDeclare(
		queue_name, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		glog.Error("QueueDeclare.Dial failed:", err, " mq_addr:", mq_addr, " queue:", queue_name)
		return
	}
	err = ch.Publish(
		"bbx_notify", // exchange
		queue_name,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	return
}

func String2Float(value string) float64 {
	if s, err := strconv.ParseFloat(value, 64); err == nil {
		return s
	}
	glog.Error("err value:", value)
	return 0.0
}

func String2Int(value string, default_val int64) int64 {
	if s, err := strconv.ParseInt(value, 10, 64); err == nil {
		return s
	}
	glog.Error("err value:", value)
	return default_val
}

func GetContractAttr(obj common.ServerQueueStruct) ContractAttr {
	contract_attr := ContractAttr{}
	contract_attr.Attr = make(map[string]interface{})
	contract_attr.Attr["symbol"] = obj.Symbol
	contract_attr.Attr["exchange"] = GetExchangeName(obj.Exchange)
	contract_attr.Attr["contractType"] = fmt.Sprintf("%d", obj.ContractType)
	return contract_attr
}

func GetContractUrl(obj common.ServerQueueStruct) string {
	contract_attr := GetContractAttr(obj)
	base_url := config.Get("contract.jump_url")
	jump_url := base_url
	for k, v := range contract_attr.Attr {
		jump_url = fmt.Sprintf("%s%s=%s&", jump_url, k, v)
	}
	if jump_url != base_url {
		jump_url = jump_url[0 : len(jump_url)-1]
	}
	// glog.Info("make jump_url:", jump_url)
	return jump_url
}

func SyncBBx(nodify_obj common.BBXNotify) {
	jsonBytes, _ := json.Marshal(nodify_obj)
	queue_name := config.Get("mq.queue_bbx_notify")
	BBxNotify(queue_name, string(jsonBytes))
	glog.Info("write queue:", queue_name, " msg:", string(jsonBytes))
}

func HttpPost(post_data string, url string) {
	conn := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := conn.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(post_data))
	if err != nil {
		glog.Error(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error("ReadAll err:", err)
		return
	}
	glog.Info(string(body))
}
