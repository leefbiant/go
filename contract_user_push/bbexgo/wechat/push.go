package wechat

import (
	"bbexgo/config"
	"bbexgo/log"
	"bbexgo/redis"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Push struct {
	MsgId     int64
	ErrorCode int64
	ErrMsg    string
	Success   bool
}

// push请求返回数据结构
type pushResult struct {
	ErrorCode int64  `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	MsgId     int64  `json:"msgid"`
}

type pushData struct {
	Touser     string                       `json:"touser"`
	TemplateID string                       `json:"template_id"`
	URL        string                       `json:"url"`
	Data       map[string]map[string]string `json:"data"`
}

type accessTokenStruct struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Errcode     int64  `json:"errcode"`
	Errmsg      string `json:"errcode"`
}

const (
	TableNameUsers = "users" // 用户表名
	AccessTokenKey = "contract_elf_wechat_access_token"
	Token          = "kNxRCaFbKAZVYZMbD7G7bJUop9Hs4UVu"
	EncodingAESKey = "u4bxwtOQWaI8kWCdfG2Q5OIIqlVdiogVKVKYRDrVRpY"
	RequestHost    = "api.weixin.qq.com" // api请求域名
)

func GetAccessToken(force bool) (string, error) {
	if force {
		AppID := config.Get("Wechat.appId")
		AppSecret := config.Get("Wechat.appSecret")
		url := fmt.Sprintf("https://%s/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", RequestHost, AppID, AppSecret)
		resp, err := http.Get(url)
		if err != nil {
			log.Error(err)
			return "", err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			return "", err
		}
		var res accessTokenStruct
		if err := json.Unmarshal(body, &res); err != nil {
			log.Error(err)
			return "", err
		}
		if res.Errcode == 0 {
			redisClient := redis.GetInstance()
			redisClient.Set(AccessTokenKey, res.AccessToken, time.Hour*2)
		}
		fmt.Println(res)
		return res.AccessToken, nil
	} else {
		redisClient := redis.GetInstance()
		accessToken, err := redisClient.Get(AccessTokenKey).Result()
		if accessToken == "" || err != nil {
			log.Error(err)
			for i := 0; i < 5; i++ {
				accessToken, e := GetAccessToken(true)
				if accessToken != "" && e == nil {
					return accessToken, nil
				}
				log.Error(e)
			}
			return "", err
		}
		return accessToken, nil
	}
}

var PushTemplateIDList map[int]string

// SendUserSub 发送push
func (p *Push) SendUserSub(userOpenID string, pushType int, msg map[string]map[string]string, jumpURL string) {
	accessToken, err := GetAccessToken(false)
	if err != nil {
		log.Error(err)
		return
	}
	apiurl := fmt.Sprintf("https://%s/cgi-bin/message/template/send?access_token=%s", RequestHost, accessToken)
	pdata, _ := json.Marshal(pushData{
		Touser:     userOpenID,
		TemplateID: p.GetPushTemplateID(pushType),
		URL:        jumpURL,
		Data:       msg,
	})
	resp, err := http.Post(apiurl, "application/x-www-form-urlencoded", strings.NewReader(string(pdata)))
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	var res pushResult
	if err := json.Unmarshal(body, &res); err != nil {
		log.Error(err)
		return
	}
	if res.ErrorCode != 0 {
		log.Error(res.ErrorCode, res.ErrMsg)
		log.Error(string(pdata))
		p.ErrorCode = res.ErrorCode
		p.ErrMsg = res.ErrMsg
		return
	}
	p.Success = true
	p.MsgId = res.MsgId
	p.ErrorCode = res.ErrorCode
	p.ErrMsg = res.ErrMsg
}

// GetPushTemplateID 根据消息类型ID获取微信消息模板ID
func (p *Push) GetPushTemplateID(typeID int) string {
	if len(PushTemplateIDList) == 0 {
		PushTemplateIDList = make(map[int]string)
	}
	if _, ok := PushTemplateIDList[typeID]; ok {
		return PushTemplateIDList[typeID]
	}
	PushTemplateIDList[typeID] = config.Get(fmt.Sprintf("Wechat.pushTemplateID.%d", typeID))
	return PushTemplateIDList[typeID]
}
