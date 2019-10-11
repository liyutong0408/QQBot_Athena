package api

import (
	"Athena/model"
	"Athena/service"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// Mahua 实现框架api接口
func Mahua(c *gin.Context) {
	// 处理接收数据
	var rj model.ReceiveJson
	if err := c.ShouldBindJSON(&rj); err != nil {
		return
	}
	//t := model.NewFramework().ConstructContext(rj)
	t := model.NewFramework().ConstructContext(rj)

	// 判断事件类型
	switch rj.TypeCode {
	case "Api_GetGroupMemberAApiOut":
		service.RefreshMemberService{}.Refresh(rj.Result)
	case "Api_GuidGetPicLinkApiOut":
		service.ArknightsService{}.UploadPic(rj.Result)
	case "EventFun":
		// debug 模式下打印接收数据
		if os.Getenv("GIN_MODE") == "debug" {
			fmt.Println(t.GetFrom() + "\t" + t.GetOperator() + "\t" + t.GetRecMsg())
		}
		switch rj.EventType {
		case 2:
			// n service数量
			n, ifResponse := 6, false
			ch := make(chan bool, n)
			// 注册service
			go service.ShutUpService{}.ShutUp(ch, *t)
			go service.INeedPicService{}.INeedPic(ch, *t)
			go service.RefreshMemberService{}.AskRefresh(ch, *t)
			go service.PerformanceService{}.SeeAll(ch, *t)
			go service.ShutUpService{}.Refresh(ch, *t)
			go service.ArknightsService{}.Arknights(ch, *t)
			// 阻塞
			for i := 0; i < n; i++ {
				temp := <-ch
				if temp == true {
					ifResponse = true
				}
			}
			//fmt.Println("解除阻塞")
			// 如果无服务做出回应
			if !ifResponse {
				service.ParrotService{}.Parrot(*t)
			}
		case 2014:
			service.ShutUpService{}.Answer(*t)
		default:
			return
		}
	}
}
