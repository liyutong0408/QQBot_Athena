package service

import (
	"Athena/model"
	"os"
	"strings"
)

type AtService struct {
}

func (service AtService) AtMe(ch chan bool, framework model.Framework) {
	if framework.GetTypeCode() != 2 {
		ch <- false
		return
	}

	if !strings.Contains(framework.GetRecMsg(), "[QQ:at="+os.Getenv("BOT")+"]") {
		ch <- false
	}

	if framework.GetRecMsg() == "[QQ:at="+os.Getenv("BOT")+"]" {
		text := "在呢在呢\n项目地址：https://github.com/Logiase/QQBot_Athena"
		framework.SetSendMsg(text).DoSendMsg()
		ch <- true
		return
	}
}
