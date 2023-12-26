package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	AppId     = "控制台获取的AppId"
	ApiKey    = "控制台获取的ApiKey"
	APISecret = "控制台获取的APISecret"
)

func main() {
	// 签名时间,延迟200毫秒,不然总是出现bad handshake
	date := time.Now().UTC().Add(200 * time.Millisecond).Format(time.RFC1123)
	// 参与签名的字段 host ,date, request-line
	signString := []string{"host: " + "tts-api.xfyun.cn", "date: " + date, "GET /v2/tts HTTP/1.1"}
	// 使用换行符拼接,因为开发文档是这么要求的
	sign := strings.Join(signString, "\n")
	/*
		HmacSha256 计算HmacSha256
		key 是控制台获取的key
		加密内容是sign
	*/
	hash := hmac.New(sha256.New, []byte(APISecret))
	hash.Write([]byte(sign))
	// 加密转为base64字符串
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	// 签名转为base64字符串
	authorization := `api_key="` + ApiKey + `", algorithm="hmac-sha256", headers="host date request-line", signature="` + sha + `"`
	authorization = base64.URLEncoding.EncodeToString([]byte(authorization))
	baseUrl := "wss://tts-api.xfyun.cn/v2/tts?authorization=" + authorization + "&date=" + url.QueryEscape(date) + "&host=tts-api.xfyun.cn"
	// 创建一个websocket默认客户端
	conn, _, err := websocket.DefaultDialer.Dial(baseUrl, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	reqJson := requestJson{
		Common: map[string]string{"app_id": AppId},
		Business: map[string]any{
			"aue": "lame",
			"sfl": 1,
			"vcn": "xiaoyan",
			"tte": "UTF8",
		},
		Data: map[string]any{
			"text":   base64.StdEncoding.EncodeToString([]byte("要转换为语音的文本")),
			"status": 2,
		},
	}
	bytes, _ := json.Marshal(reqJson)
	// 发送数据
	err = conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		panic(err)
	}
	// 接收数据
	_, message, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
  // 创建文件
	file, _ := os.Create(time.Now().Format("2006-01-02-15-04-05") + ".mp3")
	data := make(map[string]any)
	json.Unmarshal(message, &data)
	decodeString, err := base64.StdEncoding.DecodeString(data["data"].(map[string]any)["audio"].(string))
	if err != nil {
		return
	}
	file.WriteString(string(decodeString))
}

type requestJson struct {
	Common   map[string]string `json:"common"`
	Business map[string]any    `json:"business"`
	Data     map[string]any    `json:"data"`
}
