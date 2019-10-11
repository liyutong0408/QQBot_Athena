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

	if !strings.Contains(framework.GetRecMsg(), "[IR:at="+os.Getenv("BOT")+"]") {
		ch <- false
	}

	ch <- true
}
