package service

import (
	"Athena/model"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type PerformanceService struct {
}

func (service PerformanceService) SeeAll(ch chan bool, framework model.Framework) {
	if framework.GetRecMsg() != "性能状态" {
		ch <- false
		return
	}

	m, _ := mem.VirtualMemory()
	message := "内存\n" + fmt.Sprint(m.Used) + "/" + fmt.Sprint(m.Total) + "   " + fmt.Sprint(m.UsedPercent) + "%\n"

	c, _ := cpu.Percent(time.Second, false)
	message += "CPU\n" + fmt.Sprint(c[0]) + "%"

	framework.SetSendMsg(message).DoSendMsg()
	ch <- true
	return
}
