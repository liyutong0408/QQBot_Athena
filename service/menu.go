package service

import (
	"Athena/model"
	"os"
	"strings"
)

type MenuService struct {
}

func (service MenuService) Menu(ch chan bool, framework model.Framework) {
	if framework.GetTypeCode() != 2 {
		ch <- false
		return
	}

	if !strings.Contains(framework.GetRecMsg(), "[@"+os.Getenv("BOT")+"]") {
		ch <- false
	}

	text := "项目地址：https://github.com/Logiase/QQBot_Athena"

	framework.SetSendMsg(text).DoSendMsg()

	ch <- true
}
