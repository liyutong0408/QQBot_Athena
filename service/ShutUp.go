package service

import (
	"Athena/model"
	"os"
	"strconv"
	"strings"
)

type ShutUpService struct {
}

var userList []string

func (service ShutUpService) ShutUp(ch chan bool, framework model.Framework) {
	//fmt.Println("entering ShutUp")
	recMsg := framework.GetRecMsg()
	operator := framework.GetOperator()

	ifResponse := false

	for _, item := range userList {
		if item == operator {
			ifResponse = true
			break
		}
	}

	// 鉴权
	if framework.GetOperator() != os.Getenv("MASTER") && !ifResponse {
		ch <- false
		//fmt.Println("leaving ShutUp")
		return
	}

	if strings.HasPrefix(recMsg, "禁言") {
		if strings.Contains(recMsg, "[IR:at=") {
			obj := recMsg[strings.Index(recMsg, "[IR:at=")+7 : strings.Index(recMsg, "]")]
			time, err := strconv.Atoi(recMsg[strings.Index(recMsg, "]")+2:])
			if err != nil {
				framework.SetSendMsg("arg error").DoSendMsg()
				ch <- true
				//fmt.Println("leaving ShutUp")
				return
			}
			framework.DoShutUp(obj, time*60)
			framework.SetSendMsg("如您所愿,Master").DoSendMsg()
			ch <- true
			//fmt.Println("leaving ShutUp")
			return
		}
	} else if strings.HasPrefix(recMsg, "解禁") {
		if strings.Contains(recMsg, "[IR:at=") {
			obj := recMsg[strings.Index(recMsg, "[IR:at=")+7 : strings.Index(recMsg, "]")]
			//time, err := strconv.Atoi(recMsg[strings.Index(recMsg, "]")+2:])
			framework.DoShutUp(obj, 0)
			framework.SetSendMsg("如您所愿,Master").DoSendMsg()
			ch <- true
			//fmt.Println("leaving ShutUp")
			return
		}
	}

	ch <- false
	return
}

func (service ShutUpService) Answer(framework model.Framework) {
	framework.SetType(2)
	if framework.GetTrigger() == os.Getenv("MASTER") {
		whatAreYouDOING(framework)
	}

	framework.SetSendMsg("[Face178.gif][Face67.gif]").DoSendMsg()
}

func (service ShutUpService) Refresh(ch chan bool, framework model.Framework) {
	if framework.GetRecMsg() != "refresh" {
		ch <- false
		return
	}

	member, err := model.GetUser(framework.GetOperator())
	if err != nil {
		ch <- true
		framework.SetSendMsg("你好像没有人权呢...").DoSendMsg()
		return
	}
	if member.Sponsor == true || member.Role != 2 {
		framework.SetSendMsg("你就是我的Master吗？").DoSendMsg()
		for _, item := range userList {
			if framework.GetOperator() == item {
				ch <- true
				return
			}
		}

		userList = append(userList, member.QQ)
		ch <- true
		return
	}

	framework.SetSendMsg("你不是老子的Master").DoSendMsg()
	ch <- true
	return
}

func (service ShutUpService) refreshList(operator string) {
	member, err := model.GetUser(operator)
	if err != nil {
		_ = model.UpdateUser(operator, model.Member{
			QQ: operator,
		})
		return
	}

	if member.Sponsor == true {
		userList = append(userList, operator)
	}
}

func whatAreYouDOING(framework model.Framework) {
	framework.DoShutUp(os.Getenv("MASTER"), 0)
	framework.SetSendMsg("You Shall Not Pass !!!").DoSendMsg()
}
