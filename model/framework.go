package model

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

type (
	// Framework 框架模型
	Framework struct {
		client *http.Client
		ctx    msgContext
	}
	msgContext struct {
		TypeCode  int    // 事件类型
		ReceiveQQ string // 接收qq
		From      string // 消息来源
		Operator  string // 主动对象
		//Triggee   string // 被动对象
		RecMsg  string // 接收消息
		SendMsg string // 将要发送的消息
	}
	// ReceiveJson 接收json
	ReceiveJson struct {
		Result     string    `json:"Result"`
		TypeCode   string    `json:"TypeCode"`
		Message    string    `json:"Message"`
		Type       int       `json:"Type"`
		Fromgroup  string    `json:"Fromgroup"`
		Fromqq     string    `json:"Fromqq"`
		MessageID  string    `json:"MessageId"`
		Platform   int       `json:"Platform"`
		CreateTime time.Time `json:"CreateTime"`
	}
)

// NewFramework 创建Framework
func NewFramework() *Framework {
	c := &http.Client{}
	return &Framework{client: c}
}

func (framework Framework) ConstructContext(rj ReceiveJson) *Framework {
	framework.ctx = msgContext{
		TypeCode:  rj.Type,
		ReceiveQQ: os.Getenv("BOT"),
		From:      rj.Fromgroup,
		Operator:  rj.Fromqq,
		RecMsg:    rj.Message,
		SendMsg:   "",
	}
	return &framework
}
func (framework Framework) SimpleConstruct(typeCode int) *Framework {
	framework.ctx.ReceiveQQ = os.Getenv("BOT")
	framework.ctx.TypeCode = typeCode
	return &framework
}

func (framework Framework) SetSendMsg(msg string) *Framework {
	framework.ctx.SendMsg = msg
	return &framework
}
func (framework Framework) SetFrom(from string) *Framework {
	framework.ctx.From = from
	return &framework
}
func (framework Framework) SetType(t int) *Framework {
	framework.ctx.TypeCode = t
	return &framework
}

func (framework Framework) GetRecMsg() string {
	return framework.ctx.RecMsg
}
func (framework Framework) GetFrom() string {
	return framework.ctx.From
}
func (framework Framework) GetOperator() string {
	return framework.ctx.Operator
}
func (framework Framework) GetTypeCode() int {
	return framework.ctx.TypeCode
}
func (framework Framework) GetPicURL() {

}

func (framework Framework) DoSendMsg() {
	// 构建数据
	sendJSON := make(map[string]interface{})
	switch framework.ctx.TypeCode {
	case 1:
		sendJSON["类型"] = 1
		sendJSON["群组"] = framework.ctx.From
		sendJSON["qQ号"] = framework.ctx.From
		sendJSON["内容"] = framework.ctx.SendMsg
	case 2:
		sendJSON["类型"] = 2
		sendJSON["群组"] = framework.ctx.From
		sendJSON["qQ号"] = framework.ctx.Operator
		sendJSON["内容"] = framework.ctx.SendMsg
	default:
		return
	}

	// 序列化
	bytesData, _ := json.Marshal(sendJSON)

	// 发送请求
	url := "http://localhost:36524/api/v1/QQLight/Api_SendMsg"
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	framework.client.Do(req)
}
func (framework Framework) DoShutUp(obj string, t int) {
	sendJSON := make(map[string]interface{})
	sendJSON["群号"] = framework.ctx.From
	sendJSON["qq"] = obj
	sendJSON["禁言时长"] = t

	bytesData, _ := json.Marshal(sendJSON)
	url := "http://localhost:36524/api/v1/QQLight/Api_Ban"

	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	framework.client.Do(req)
}
func (framework Framework) DoJSONMusic(name string) {
	sendJSON := make(map[string]interface{})
	sendJSON["歌曲ID"] = name

	bytesData, _ := json.Marshal(sendJSON)
	url := "http://localhost:36524/api/v1/QQLight/Api_JsonMusic"

	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	framework.client.Do(req)
}
func (framework Framework) Do163Music(id string) {
	sendJSON := make(map[string]interface{})
	sendJSON["歌曲ID"] = id

	bytesData, _ := json.Marshal(sendJSON)
	url := "http://localhost:36524/api/v1/QQLight/Api_163Music"

	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	framework.client.Do(req)
}

// stop using
func (framework Framework) DoSendObjectMsg() {
	sendJSON := make(map[string]interface{})
	sendJSON["响应的QQ"] = framework.ctx.ReceiveQQ
	sendJSON["收信对象类型"] = framework.ctx.TypeCode
	sendJSON["收信对象所属群_讨论组"] = framework.ctx.From
	sendJSON["收信对象QQ"] = framework.ctx.Operator
	sendJSON["objectMsg"] = framework.ctx.SendMsg
	sendJSON["结构子类型"] = 02

	bytesData, _ := json.Marshal(sendJSON)

	url := "http://localhost:36524/api/v1/Mpq/Api_SendObjectMsg"
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	framework.client.Do(req)
}
func (framework Framework) DoGetGroupMember() {
	sendJSON := make(map[string]interface{})
	sendJSON["群号"] = framework.ctx.From

	bytesData, _ := json.Marshal(sendJSON)
	url := "http://localhost:36524/api/v1/QQLight/Api_GetGroupMemberList"

	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	framework.client.Do(req)
}

func (framework Framework) IsRecMsgPic() bool {
	if strings.Contains(framework.ctx.RecMsg, "[QQ:pic=") {
		if strings.Contains(framework.ctx.RecMsg, "]") {
			return true
		}
	}
	return false
}
func (framework Framework) IsRecMsgContainAT() (bool, []string) {
	obj := []string{}
	str := framework.ctx.RecMsg
	flag := false

	for {
		if !strings.Contains(str, "[QQ:at=") {
			break
		}

		if strings.Contains(str, "]") {
			obj = append(obj, str[strings.Index(str, "[QQ:at=")+7:strings.Index(str, "]")])
			str = str[strings.Index(str, "]")+1:]
			flag = true
		}
	}

	return flag, obj
}

func GetConstPic(picURL string) string {
	return "[QQ:pic=" + picURL + "]"
}
func GetConstAT(obj string) string {
	return "[QQ:at=" + obj + "]"
}
