package model

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
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
		Triggee   string // 被动对象
		RecMsg    string // 接收消息
		SendMsg   string // 将要发送的消息
	}
	// ReceiveJson 接收json
	ReceiveJson struct {
		Result            string `json:"Result"`
		CreateTime        string `json:"CreateTime"`
		EventAdditionType int    `json:"EventAdditionType"`
		EventOperator     string `json:"EventOperator"`
		EventType         int    `json:"EventType"`
		FromNum           string `json:"FromNum"`
		JSON              string `json:"Json"`
		Message           string `json:"Message"`
		MessageID         string `json:"MessageId"`
		MessageNum        string `json:"MessageNum"`
		Platform          int    `json:"Platform"`
		RawMessage        string `json:"RawMessage"`
		ReceiverQq        string `json:"ReceiverQq"`
		Triggee           string `json:"Triggee"`
		TypeCode          string `json:"TypeCode"`
	}
)

// NewFramework 创建Framework
func NewFramework() *Framework {
	c := &http.Client{}
	return &Framework{client: c}
}

func (framework Framework) ConstructContext(rj ReceiveJson) *Framework {
	framework.ctx = msgContext{
		TypeCode:  rj.EventType,
		ReceiveQQ: rj.ReceiverQq,
		From:      rj.FromNum,
		Operator:  rj.EventOperator,
		Triggee:   rj.Triggee,
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
func (framework Framework) GetTrigger() string {
	return framework.ctx.Triggee
}
func (framework Framework) GetPicURL() {
	sendJSON := make(map[string]interface{})
	sendJSON["图片GUID"] = framework.ctx.RecMsg

	bytesData, _ := json.Marshal(sendJSON)

	// 发送请求
	url := "http://localhost:36524/api/v1/Mpq/Api_GuidGetPicLink"
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	framework.client.Do(req)
}

func (framework Framework) DoSendMsg() {
	// 构建数据
	sendJSON := make(map[string]interface{})
	switch framework.ctx.TypeCode {
	case 1:
		sendJSON["响应的QQ"] = framework.ctx.ReceiveQQ
		sendJSON["信息类型"] = 1
		sendJSON["参考子类型"] = 0
		sendJSON["收信群_讨论组"] = ""
		sendJSON["收信对象"] = framework.ctx.From
		sendJSON["内容"] = framework.ctx.SendMsg
	case 2:
		sendJSON["响应的QQ"] = framework.ctx.ReceiveQQ
		sendJSON["信息类型"] = 2
		sendJSON["参考子类型"] = 0
		sendJSON["收信群_讨论组"] = framework.ctx.From
		sendJSON["收信对象"] = ""
		sendJSON["内容"] = framework.ctx.SendMsg
	default:
		return
	}

	// 序列化
	bytesData, _ := json.Marshal(sendJSON)

	// 发送请求
	url := "http://localhost:36524/api/v1/Mpq/Api_SendMsg"
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	framework.client.Do(req)
}
func (framework Framework) DoShutUp(obj string, t int) {
	sendJSON := make(map[string]interface{})
	sendJSON["响应的QQ"] = framework.ctx.ReceiveQQ
	sendJSON["群号"] = framework.ctx.From
	sendJSON["qq"] = obj
	sendJSON["时长"] = t

	bytesData, _ := json.Marshal(sendJSON)
	url := "http://localhost:36524/api/v1/Mpq/Api_Shutup"

	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	framework.client.Do(req)
}
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
	sendJSON["响应的QQ"] = framework.ctx.ReceiveQQ
	sendJSON["群号"] = framework.ctx.From

	bytesData, _ := json.Marshal(sendJSON)
	url := "http://localhost:36524/api/v1/Mpq/Api_GetGroupMemberA"

	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	framework.client.Do(req)
}
