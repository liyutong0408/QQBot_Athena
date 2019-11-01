package service

import (
	"Athena/model"
)

type (
	GMChangeService struct {
	}
)

func (service GMChangeService) GMDecrease(framework model.Framework) {
	msg := "有人退群了，，，"
	framework.SetType(2).SetSendMsg(msg).DoSendMsg()
	return
}

func (service GMChangeService) GMIncrease(framework model.Framework) {
	msg := "欢迎新大佬"
	if framework.GetOprater() != "" {
		msg = "欢迎" + model.GetConstAT(framework.GetOprater()) + "邀请的新大佬" + "\n"
	}

	msg += model.GetConstAT(framework.GetFromQQ())

	framework.SetType(2).SetSendMsg(msg).DoSendMsg()
	return
}
