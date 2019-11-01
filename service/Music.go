package service

import (
	"Athena/model"
	"strings"
)

type (
	MusicService struct {
	}
)

var (
	musicLastGroup string
)

func (service MusicService) Music(ch chan bool, framework model.Framework) {
	if !strings.HasPrefix(framework.GetRecMsg(), "点歌") {
		ch <- false
		return
	}

	songName := strings.TrimLeft(framework.GetRecMsg(), "点歌")
	songName = strings.TrimLeft(songName, " ")
	framework.DoJSONMusic(songName)
	musicLastGroup = framework.GetFromGroup()

	ch <- true
	return
}

func (service MusicService) SendMusic(result string) {
	model.NewFramework().SimpleConstruct(2).SetFromGroup(musicLastGroup).SetSendMsg(result).DoSendMsg()
}
