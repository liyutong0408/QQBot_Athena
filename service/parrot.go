package service

import "Athena/model"

type ParrotService struct {
}

var parrotMap = make(map[string]*parrotGroup)

type parrotGroup struct {
	strTemp      [3]string
	flagParroted bool
}

func (service ParrotService) Parrot(framework model.Framework) {
	if _, ok := parrotMap[framework.GetFromGroup()]; !ok {
		parrotMap[framework.GetFromGroup()] = &parrotGroup{
			flagParroted: false,
		}
	}

	parrotMap[framework.GetFromGroup()].strTemp[2] = parrotMap[framework.GetFromGroup()].strTemp[1]
	parrotMap[framework.GetFromGroup()].strTemp[1] = parrotMap[framework.GetFromGroup()].strTemp[0]
	parrotMap[framework.GetFromGroup()].strTemp[0] = framework.GetRecMsg()

	if parrotMap[framework.GetFromGroup()].strTemp[1] == parrotMap[framework.GetFromGroup()].strTemp[0] {
		if parrotMap[framework.GetFromGroup()].strTemp[2] == parrotMap[framework.GetFromGroup()].strTemp[1] {
			if parrotMap[framework.GetFromGroup()].flagParroted {
				return
			}
			framework.SetSendMsg(framework.GetRecMsg()).DoSendMsg()
			parrotMap[framework.GetFromGroup()].flagParroted = true
		}
		return
	} else {
		parrotMap[framework.GetFromGroup()].flagParroted = false
		return
	}
}
