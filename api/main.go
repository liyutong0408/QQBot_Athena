package api

import (
	"Athena/model"
	"Athena/service"
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
	case "GetNewMsg":
		// n service数量
		n, ifResponse := 9, false
		ch := make(chan bool, n)
		// 注册service
		go service.ShutUpService{}.ShutUp(ch, *t)
		go service.ShutUpService{}.Answer(ch, *t)
		go service.INeedPicService{}.INeedPic(ch, *t)
		go service.RefreshMemberService{}.AskRefresh(ch, *t)
		go service.PerformanceService{}.SeeAll(ch, *t)
		go service.ShutUpService{}.Refresh(ch, *t)
		go service.ArknightsService{}.Arknights(ch, *t)
		go service.MusicService{}.Music(ch, *t)
		go service.AtService{}.AtMe(ch, *t)
		// 阻塞
		for i := 0; i < n; i++ {
			temp := <-ch
			if temp == true {
				ifResponse = true
			}
		}

		// 如果无服务做出回应
		if !ifResponse {
			service.ParrotService{}.Parrot(*t)
		}

		return
	case "Api_JsonMusicApiOut":
		service.MusicService{}.SendMusic(rj.Result)
		return
	default:
		return
	}
}
