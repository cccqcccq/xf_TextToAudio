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
	HostUrl   = "wss://tts-api.xfyun.cn/v2/tts" // 文本生成语音api
)

func main() {
	// 创建一个websocket默认客户端
	conn, _, err := websocket.DefaultDialer.Dial(createUrl(), nil)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	// 发送数据
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
	err = conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		panic(err)
	}
	// 创建文件
	file, _ := os.Create(time.Now().Format("2006-01-02-15-04-05") + ".mp3")
	defer file.Close()
	for {
		// 循环读取数据
		_, message, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		// 判断是否是结束
		if createFile(file, message) == 2 {
			break
		}
	}
	fmt.Println("创建成功")
}

func createFile(file *os.File, message []byte) int {
	// 写入文件
	data := make(map[string]any)
	json.Unmarshal(message, &data)
	decodeString, _ := base64.StdEncoding.DecodeString(data["data"].(map[string]any)["audio"].(string))
	status, _ := data["data"].(map[string]any)["status"].(float64)
	file.WriteString(string(decodeString))
	// 返回创建状态
	return int(status)
}

func createUrl() (baseurl string) {
	ul, err := url.Parse(HostUrl)
	if err != nil {
		return
	}
	// 签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	// 参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	// 使用换行符拼接,因为开发文档是这么要求的
	sign := strings.Join(signString, "\n")
	/*
		HmacSha256 计算HmacSha256
		key 是控制台获取的key
		加密内容是sign
	*/
	hash := hmac.New(sha256.New, []byte(APISecret))
	hash.Write([]byte(sign))
	// 加密转为base64字符串,一定要使用StdEncoding,如果使用URLEncoding按照url规则转会导致签名结果偶尔不一致而出现偶尔可以链接成功偶尔失败的原因
	sha := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	// 签名转为base64字符串
	authorization := fmt.Sprintf("api_key=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", ApiKey,
		"hmac-sha256", "host date request-line", sha)
	authorization = base64.StdEncoding.EncodeToString([]byte(authorization))
	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	baseurl = HostUrl + "?" + v.Encode()
	return
}

type requestJson struct {
	Common   map[string]string `json:"common"`
	Business map[string]any    `json:"business"`
	Data     map[string]any    `json:"data"`
}

