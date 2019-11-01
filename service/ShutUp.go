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
	operator := framework.GetFromQQ()

	ifResponse := false

	for _, item := range userList {
		if item == operator {
			ifResponse = true
			break
		}
	}

	// 鉴权
	if framework.GetFromQQ() != os.Getenv("MASTER") && !ifResponse {
		ch <- false
		//fmt.Println("leaving ShutUp")
		return
	}

	if strings.HasPrefix(recMsg, "禁言") {
		if ok, objL := framework.IsRecMsgContainAT(); ok {
			if strings.LastIndex(recMsg, "]")+2 > len(recMsg) {
				framework.SetSendMsg("arg error").DoSendMsg()
				ch <- true
				return
			}
			time, err := strconv.Atoi(recMsg[strings.LastIndex(recMsg, "]")+2:])
			if err != nil {
				framework.SetSendMsg("arg error").DoSendMsg()
				ch <- true
				//fmt.Println("leaving ShutUp")
				return
			}
			for _, obj := range objL {
				framework.DoShutUp(obj, time*60)
			}
			framework.SetSendMsg("如您所愿,Master").DoSendMsg()
			ch <- true
			//fmt.Println("leaving ShutUp")
			return
		}
	} else if strings.HasPrefix(recMsg, "解禁") {
		if ok, objL := framework.IsRecMsgContainAT(); ok {
			for _, obj := range objL {
				framework.DoShutUp(obj, 0)
			}

			framework.SetSendMsg("如您所愿,Master").DoSendMsg()
			ch <- true
			//fmt.Println("leaving ShutUp")
			return
		}
	}

	ch <- false
	return
}

func (service ShutUpService) Answer(ch chan bool, framework model.Framework) {
	if !strings.HasPrefix(framework.GetRecMsg(), "[QQ:禁言QQ=") {
		ch <- false
		return
	}

	pos1 := strings.Index(framework.GetRecMsg(), "=")
	pos2 := strings.Index(framework.GetRecMsg(), ",")

	if framework.GetRecMsg()[pos1+1:pos2] == os.Getenv("MASTER") {
		whatAreYouDOING(framework)
	}

	//framework.SetSendMsg("[QQ:face=178][QQ:face=67]").DoSendMsg()
	ch <- true
	return
}

func (service ShutUpService) Refresh(ch chan bool, framework model.Framework) {
	if framework.GetRecMsg() != "refresh" {
		ch <- false
		return
	}

	member, err := model.GetUser(framework.GetFromQQ())
	if err != nil {
		ch <- true
		framework.SetSendMsg("你好像没有人权呢...").DoSendMsg()
		return
	}
	if member.Sponsor == true || member.Role != 2 {
		framework.SetSendMsg("你就是我的Master吗？").DoSendMsg()
		for _, item := range userList {
			if framework.GetFromQQ() == item {
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
