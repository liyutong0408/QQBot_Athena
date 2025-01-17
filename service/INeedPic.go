package service

import (
	"Athena/model"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type (
	INeedPicService struct {
	}
	picSource string
	setu      struct {
		Pid int    `json:"pid"`
		P   int    `json:"p"`
		URL string `json:"url"`
		R18 bool   `json:"r18"`
	}
)

const (
	pixiv picSource = "https://api.pixivic.com/illust"
	pic0  picSource = "https://s0.xinger.ink/acgimg/acgurl.php"
	pic1  picSource = "https://sotama.cool/picture"
	pic2  picSource = "http://www.dmoe.cc/random.php"
	pic3  picSource = "http://laoliapi.cn/king/tupian/2cykj"
	pic4  picSource = "http://acg.bakayun.cn/randbg.php"
	pic5  picSource = "https://acg.toubiec.cn/random"
	pic6  picSource = "http://pic.tsmp4.net/api/erciyuan/img.php"
	pic7  picSource = "https://random.52ecy.cn/randbg.php"
)

var (
	piclist    = [8]picSource{pic0, "pic1,再来一次吧", pic2, pic3, "pic4,再来一次吧", pic5, pic6, "pic7,再来一次吧"}
	setuswitch = false
)

// INeedPic service entry
func (service INeedPicService) INeedPic(ch chan bool, framework model.Framework) {
	if framework.GetFromQQ() == os.Getenv("MASTER") {
		if framework.GetRecMsg() == "瑟图time" {
			setuswitch = true
			framework.SetSendMsg("瑟图time").DoSendMsg()
			ch <- true
			return
		}
		if framework.GetRecMsg() == "瑟图stop" {
			setuswitch = false
			ch <- true
			framework.SetSendMsg("master不让了").DoSendMsg()
			return
		}
	}

	if setuswitch != true {
		ch <- false
		return
	}

	//fmt.Println("entering INeedPic")
	if !strings.HasPrefix(framework.GetRecMsg(), "一张瑟图") {
		ch <- false
		//fmt.Println("leaving INeedPic")
		return
	}

	if framework.GetRecMsg() == "一张瑟图" {
		var resp *http.Response
		var err error
		for {
			resp, err = http.Get("https://api.lolicon.app/setu/")
			if err == nil {
				break
			}
		}
		body, _ := ioutil.ReadAll(resp.Body)
		var st setu
		json.Unmarshal(body, &st)
		framework.SetSendMsg(model.GetConstPic(st.URL)).DoSendMsg()
		/*
			rd := rand.New(rand.NewSource(time.Now().UnixNano()))
			framework.SetSendMsg(string(piclist[rd.Intn(7)])).DoSendMsg()

		*/
		ch <- true
		//fmt.Println("leaving INeedPic")
		return
	}

	arg := framework.GetRecMsg()[13:]
	if arg == "pixiv" {
		framework.SetSendMsg(model.GetConstPic(string(pixiv))).DoSendMsg()
		ch <- true
		//fmt.Println("leaving INeedPic")
		return
	}
	i, err := strconv.Atoi(arg)
	if err != nil {
		framework.SetSendMsg("参数错误").DoSendMsg()
		ch <- true
		//fmt.Println("leaving INeedPic")
		return
	}
	if i > 7 || i < 0 {
		framework.SetSendMsg("超出范围").DoSendMsg()
		ch <- true
		//fmt.Println("leaving INeedPic")
		return
	} else {
		framework.SetSendMsg(model.GetConstPic(string(piclist[i]))).DoSendMsg()
		ch <- true
		//fmt.Println("leaving INeedPic")
		return
	}
}
