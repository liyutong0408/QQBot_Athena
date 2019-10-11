package tasks

import "Athena/model"

func WarnAbyss() error {
	framework := model.NewFramework().SimpleConstruct(2).SetFrom("547902826")

	framework.SetSendMsg("Athena 提醒:\n该打深渊啦！！！").DoSendMsg()
	return nil
}
