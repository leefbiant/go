package bitmax_api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	// "strconv"
	"elf/bbexgo/log"
	"strings"
	"time"
	"sync"
)
var (
	mutex sync.Mutex
)
type Client struct {
	Config     Config
	HttpClient *http.Client
}

type ApiMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(config Config, proxyUrl string) *Client {
	var client Client
	client.Config = config
	timeout := 10
	if proxyUrl != "" {
		proxy, _ := url.Parse(proxyUrl)
		client.HttpClient = &http.Client{
			Timeout:   time.Duration(timeout) * time.Second,
			Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
		}
		log.Info("use proxyUrl:", proxyUrl)
	} else {
		client.HttpClient = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}
	return &client
}

func (client *Client) Close() {

}

func (client *Client) Request(method string, requestPath string,
	uri_params map[string]interface{}, body_params map[string]interface{},
	result interface{}) (response *http.Response, err error) {

	config := client.Config
	// uri
	endpoint := config.Endpoint
	if strings.HasSuffix(config.Endpoint, "/") {
		endpoint = config.Endpoint[0 : len(config.Endpoint)-1]
	}

	// get json and bin styles request body
	var uri_data string
	var binBody = bytes.NewReader(make([]byte, 0))
	if len(uri_params) > 0 {
		uri_data, err = ParseUriParams(uri_params, method)
		if err != nil {
			log.Error("ParseRequestParams err:", err)
			return response, err
		}
		/* fmt.Println("uri_params uri_data:", uri_data) */
		requestPath = requestPath + uri_data
	}

	var jsonBody string
	if len(body_params) > 0 {
		jsonBody, binBody, _ = ParseBodyParams(body_params)
		/* fmt.Println("body_params body:", jsonBody) */
	}

	url := endpoint + requestPath

	// get a http request
	request, err := http.NewRequest(method, url, binBody)
	if err != nil {
		log.Error("NewRequest err:", err)
		return response, err
	}
	mutex.Lock()
	
	// Sign and set request headers
	timestamp := fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond)) 
	preHash := PreHashString(timestamp, method, requestPath, jsonBody)
	sign := BitmaxSigner(preHash, config.SecretKey)

	Headers(request, config, timestamp, sign)
	if config.IsPrint {
		printRequest(config, request, jsonBody, preHash)
	}

	// send a request to remote server, and get a response
	response, err = client.HttpClient.Do(request)
	mutex.Unlock()
	if err != nil {
		log.Error("httpClient err:", err, " request:", request)
		return response, err
	}
	defer response.Body.Close()

	// get a response results and parse
	status := response.StatusCode
	message := response.Status
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("ReadAll err:", err)
		return response, err
	}

	if config.IsPrint {
		printResponse(status, message, body)
	}

	responseBodyString := string(body)
	response.Header.Add(ResultDataJsonString, responseBodyString)

	limit := response.Header.Get("x-ratelimit-limit")
	if limit != "" {
		var page PageResult
		page.RatelimitRemaining = StringToInt(limit)
		limit = response.Header.Get("x-ratelimit-reset")
		if limit != "" {
			page.RatelimitReset = StringToInt(limit)
		}
		limit = response.Header.Get("x-ratelimit-limit")
		if limit != "" {
			page.Ratelimitlimit = StringToInt(limit)
		}
		pageJsonString, err := Struct2JsonString(page)
		if err == nil {
			response.Header.Add(ResultPageJsonString, pageJsonString)
		}
	}

	if status >= 200 && status < 300 {
		if body != nil && result != nil {
			err := JsonBytes2Struct(body, result)
			if err != nil {
				log.Error("JsonBytes2Struct err:", err, " body:", body)
				return response, err
			}
		}
		return response, nil
	} else if status >= 400 || status <= 500 {
		log.Error("Http error(400~500) result: status=" + IntToString(status) + ", message=" + message + ", body=" + responseBodyString)
		if body != nil {
			log.Error("http status err:", status, " body:", body)
			return response, errors.New(string(body))
		}
	} else {
		log.Error("Http error result: status=" + IntToString(status) + ", message=" + message + ", body=" + responseBodyString)
		return response, errors.New(message)
	}
	log.Info("Request sucess")
	return response, nil
}
