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
	if _, ok := parrotMap[framework.GetFrom()]; !ok {
		parrotMap[framework.GetFrom()] = &parrotGroup{
			flagParroted: false,
		}
	}

	parrotMap[framework.GetFrom()].strTemp[2] = parrotMap[framework.GetFrom()].strTemp[1]
	parrotMap[framework.GetFrom()].strTemp[1] = parrotMap[framework.GetFrom()].strTemp[0]
	parrotMap[framework.GetFrom()].strTemp[0] = framework.GetRecMsg()

	if parrotMap[framework.GetFrom()].strTemp[1] == parrotMap[framework.GetFrom()].strTemp[0] {
		if parrotMap[framework.GetFrom()].strTemp[2] == parrotMap[framework.GetFrom()].strTemp[1] {
			if parrotMap[framework.GetFrom()].flagParroted {
				return
			}
			framework.SetSendMsg(framework.GetRecMsg()).DoSendMsg()
			parrotMap[framework.GetFrom()].flagParroted = true
		}
		return
	} else {
		parrotMap[framework.GetFrom()].flagParroted = false
		return
	}
}
