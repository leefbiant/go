package bitmax_api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"elf/bbexgo/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

func Int64ToString(arg int64) string {
	return strconv.FormatInt(arg, 10)
}

func IntToString(arg int) string {
	return strconv.Itoa(arg)
}

func StringToInt64(arg string) int64 {
	value, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0
	} else {
		return value
	}
}

func StringToInt(arg string) int {
	value, err := strconv.Atoi(arg)
	if err != nil {
		return 0
	} else {
		return value
	}
}

func JsonBytes2Struct(jsonBytes []byte, result interface{}) error {
	err := json.Unmarshal(jsonBytes, result)
	return err
}

func ParseUriParams(params map[string]interface{}, method string) (string, error) {
	if params == nil {
		return "", errors.New("illegal parameter")
	}
	var byte_data []byte

	data_str := ""
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if data_str == "" {
			data_str = fmt.Sprintf("%s?%s=%v", data_str, k, params[k])
		} else {
			data_str = fmt.Sprintf("%s&%s=%v", data_str, k, params[k])
		}
	}
	byte_data = []byte(data_str)
	jsonBody := string(byte_data)
	return jsonBody, nil
}

func ParseBodyParams(params map[string]interface{}) (string, *bytes.Reader, error) {
	if params == nil {
		return "", nil, errors.New("illegal parameter")
	}
	byte_data, err := json.Marshal(params)
	if err != nil {
		return "", nil, errors.New("json convert string error")
	}
	jsonBody := string(byte_data)
	binBody := bytes.NewReader(byte_data)
	return jsonBody, binBody, nil
}

func PreHashString(timestamp string, method string, requestPath string, body string) string {
	return strings.ToUpper(method) + requestPath + timestamp + body
}

func Headers(request *http.Request, config Config, timestamp string, sign string) {
	request.Header.Add(ACCEPT, APPLICATION_JSON)
	request.Header.Add(CONTENT_TYPE, APPLICATION_JSON_UTF8)
	request.Header.Add(API_KEY, config.ApiKey)
	request.Header.Add(API_NONCE, timestamp)
	request.Header.Add(API_SIGNATURE, sign)
}

func BitmaxSigner(message string, secretKey string) string {
	key := []byte(secretKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func printRequest(config Config, request *http.Request, body string, preHash string) {
	if true {

/* 		if config.SecretKey != "" {
			log.Info("  Secret-Key: " + config.SecretKey)
		}
		log.Info("  Request(" + IsoTime() + "):") */

		log.Info("\tUrl: " + request.URL.String())

/* 		log.Info("\tMethod: " + strings.ToUpper(request.Method))
		if len(request.Header) > 0 {
			log.Info("\tHeaders: ")
			for k, v := range request.Header {
				if strings.Contains(k, "Ok-") {
					k = strings.ToUpper(k)
				}
				log.Info("\t\t" + k + ": " + v[0])
			}
		}
		log.Info("\tBody: " + body)
		if preHash != "" {
			log.Info("  PreHash: " + preHash)
		} */
	}
}

func printResponse(status int, message string, body []byte) {
	if status != 200 {
		log.Error("  Response(" + IsoTime() + "):")
		statusString := strconv.Itoa(status)
		message = strings.Replace(message, statusString, "", -1)
		message = strings.Trim(message, " ")
		log.Error("\tStatus: " + statusString)
		log.Error("\tMessage: " + message)
	}
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "\t")
	log.Info("Body:\n" + string(prettyJSON.Bytes()))
}

func Struct2JsonString(structt interface{}) (jsonString string, err error) {
	data, err := json.Marshal(structt)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
