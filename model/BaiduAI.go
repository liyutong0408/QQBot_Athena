package model

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var (
	token   string
	iftoken bool = false
)

// GetToken 获取Token
func bAIGetToken() {
	client := http.Client{}
	url := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + os.Getenv("BAAPIKEY") + "&client_secret=" + os.Getenv("BASECRETKEY")
	var tokenj tokenjson
	fmt.Println(url)

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &tokenj)
	if err != nil {
		panic(err)
	}
	token = tokenj.AccessToken
	iftoken = true
	//fmt.Println(tokenj)
}

// UploadPic 上传本地图片
func BAIUploadPic(path string) Recjson {
	if !iftoken {
		bAIGetToken()
	}

	bytesdata, _ := ioutil.ReadFile(path)
	encodestr := base64.StdEncoding.EncodeToString(bytesdata)

	v := url.Values{}
	v.Add("image", encodestr)
	b := v.Encode()
	client := http.Client{}
	req, _ := http.NewRequest("POST", "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token="+token, bytes.NewReader([]byte(b)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)

	var rj Recjson
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &rj)
	return rj
	//fmt.Println(rj)
	//9 10 11 14 15
}

// UploadPicFromURL 通过url上传图片
func BAIUploadPicFromURL(urls string) Recjson {
	if !iftoken {
		bAIGetToken()
	}
	v := url.Values{}
	v.Add("url", urls)
	b := v.Encode()
	client := http.Client{}
	req, _ := http.NewRequest("POST", "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token="+token, bytes.NewReader([]byte(b)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)

	var rj Recjson
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &rj)

	fmt.Println(rj)
	if rj.ErrorCode != 0 {
		bAIGetToken()
		BAIUploadPicFromURL(urls)
	}
	//fmt.Println(rj)
	return rj
}

// Recjson 接收json格式
type (
	Recjson struct {
		LogID          int64 `json:"log_id"`
		WordsResultNum int   `json:"words_result_num"`
		WordsResult    []struct {
			Words string `json:"words"`
		} `json:"words_result"`
		ErrorCode int `json:"error_code"`
	}
	tokenjson struct {
		RefreshToken     string `json:"refresh_token"`
		ExpiresIn        int    `json:"expires_in"`
		Scope            string `json:"scope"`
		SessionKey       string `json:"session_key"`
		AccessToken      string `json:"access_token"`
		SessionSecret    string `json:"session_secret"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
)
