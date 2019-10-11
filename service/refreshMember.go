package service

import (
	"Athena/model"
	"encoding/json"
	"os"
	"strconv"
)

type RefreshMemberService struct {
}

const admittedGroup = "547902826"

func (service RefreshMemberService) AskRefresh(ch chan bool, framework model.Framework) {
	//fmt.Println("entering AskRefresh")
	if framework.GetFrom() != admittedGroup {
		ch <- false
		//fmt.Println("leaving AskRefresh")
		return
	}

	if framework.GetRecMsg() != "收束" {
		ch <- false
		//fmt.Println("leaving AskRefresh")
		return
	}

	if framework.GetOperator() != os.Getenv("MASTER") {
		framework.SetSendMsg("权限不足").DoSendMsg()
		ch <- true
		//fmt.Println("leaving AskRefresh")
		return
	}

	framework.DoGetGroupMember()
	ch <- true
	//fmt.Println("leaving AskRefresh")
	return
}

func (service RefreshMemberService) Refresh(result string) {
	var t gm
	json.Unmarshal([]byte(result), &t)

	for _, item := range t.Mems {
		member := model.Member{
			QQ:         strconv.Itoa(item.Uin),
			NickName:   item.Card,
			Sponsor:    false,
			Role:       item.Role,
			Money:      0,
			GameServer: 0,
			GameName:   "",
		}

		model.DB.Where(model.Member{QQ: member.QQ}).FirstOrCreate(&model.Member{}, member)
	}

	model.NewFramework().SimpleConstruct(2).SetFrom("547902826").SetSendMsg("finish").DoSendMsg()
}

type gm struct {
	Ec        int    `json:"ec"`
	Errcode   int    `json:"errcode"`
	Em        string `json:"em"`
	AdmNum    int    `json:"adm_num"`
	AdmMax    int    `json:"adm_max"`
	Vecsize   int    `json:"vecsize"`
	Levelname struct {
		Num1   string `json:"1"`
		Num2   string `json:"2"`
		Num3   string `json:"3"`
		Num4   string `json:"4"`
		Num5   string `json:"5"`
		Num6   string `json:"6"`
		Num10  string `json:"10"`
		Num11  string `json:"11"`
		Num12  string `json:"12"`
		Num13  string `json:"13"`
		Num14  string `json:"14"`
		Num15  string `json:"15"`
		Num101 string `json:"101"`
		Num102 string `json:"102"`
		Num103 string `json:"103"`
		Num104 string `json:"104"`
		Num105 string `json:"105"`
		Num106 string `json:"106"`
		Num107 string `json:"107"`
		Num108 string `json:"108"`
		Num109 string `json:"109"`
		Num110 string `json:"110"`
		Num111 string `json:"111"`
		Num112 string `json:"112"`
		Num113 string `json:"113"`
		Num114 string `json:"114"`
		Num115 string `json:"115"`
		Num116 string `json:"116"`
		Num117 string `json:"117"`
		Num118 string `json:"118"`
		Num197 string `json:"197"`
		Num198 string `json:"198"`
		Num199 string `json:"199"`
	} `json:"levelname"`
	Mems []struct {
		Uin           int `json:"uin"`
		Role          int `json:"role"`
		Flag          int `json:"flag"`
		G             int `json:"g"`
		JoinTime      int `json:"join_time"`
		LastSpeakTime int `json:"last_speak_time"`
		Lv            struct {
			Point int `json:"point"`
			Level int `json:"level"`
		} `json:"lv"`
		Nick string `json:"nick"`
		Card string `json:"card"`
		Qage int    `json:"qage"`
		Tags string `json:"tags"`
		Rm   int    `json:"rm"`
	} `json:"mems"`
	Count       int `json:"count"`
	SvrTime     int `json:"svr_time"`
	MaxCount    int `json:"max_count"`
	SearchCount int `json:"search_count"`
}
