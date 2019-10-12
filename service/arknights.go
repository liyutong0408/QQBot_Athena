package service

import (
	"Athena/model"
	"encoding/json"
	"fmt"
	combinations "github.com/mxschmitt/golang-combinations"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

var (
	hrData       hrDataT
	tagList      = ""
	ifNew        = false
	arkMap       = make(map[string]bool)
	arkLastGroup string
)

type (
	ArknightsService struct {
	}
	resultT struct {
		comb  []string
		chars []int
	}
	arkRecJSON struct {
		Name           string   `json:"name"`
		Camp           string   `json:"camp"`
		Type           string   `json:"type"`
		Level          int      `json:"level"`
		Sex            string   `json:"sex"`
		Characteristic string   `json:"characteristic"`
		Tags           []string `json:"tags"`
		Hidden         bool     `json:"hidden"`
		NameEn         string   `json:"name-en,omitempty"`
	}
	// 满足sort interface
	arkRecJSONSlice []arkRecJSON
	hrDataT         struct {
		characters []characterT
		data       map[string][]int
		avgCharTag float64
	}
	characterT struct {
		name string
		rare int
	}
)

func (service ArknightsService) Arknights(ch chan bool, framework model.Framework) {
	if framework.GetRecMsg() == "更新数据" && framework.GetOperator() == os.Getenv("MASTER") {
		initHrData()
		framework.SetSendMsg("数据更新完成").DoSendMsg()
		ch <- true
		return
	}
	// start ocr
	if framework.GetRecMsg() == "公招" {
		if _, ok := arkMap[framework.GetFrom()]; !ok {
			arkMap[framework.GetFrom()] = true

		} else {
			arkMap[framework.GetFrom()] = true
		}
		framework.SetSendMsg("请发送公招截图").DoSendMsg()
		ch <- true
		return
	}

	fmt.Println("start upload", arkMap[framework.GetFrom()])
	if arkMap[framework.GetFrom()] {
		if strings.Index(framework.GetRecMsg(), "{") == 0 {
			if strings.Contains(framework.GetRecMsg(), "}") {
				fmt.Println("entering")
				framework.GetPicURL()

				arkLastGroup = framework.GetFrom()
				ch <- true
				arkMap[framework.GetFrom()] = false
				return
			}
		}
	}

	ch <- false
	return
}

func (service ArknightsService) UploadPic(result string) {
	fmt.Println("entering up")
	r := model.BAIUploadPicFromURL(result)
	var tags []string
	for _, item := range r.WordsResult {
		if strings.Contains(tagList, item.Words) {
			tags = append(tags, item.Words)
		}
	}

	combs := getCombination(tags)

	text := "tag:"
	for _, tag := range tags {
		text += tag + " "
	}

	for _, r := range combs {
		text += "\n\n"
		for _, tag := range r.comb {
			text += tag + " "
		}
		text += "\n"
		for _, char := range r.chars {
			text += getChar(char) + " "
		}
	}

	model.NewFramework().SimpleConstruct(2).SetFrom(arkLastGroup).SetSendMsg(text).DoSendMsg()
	fmt.Println("leaving up")
}

func pullData() hrDataT {
	ret := newHrData()
	// 获取最新数据
	resp, err := http.Get("https://graueneko.github.io/akhr.json")
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var rj []arkRecJSON
	_ = json.Unmarshal(body, &rj)

	// 去除rj中12星干员
	var data []arkRecJSON
	for _, item := range rj {
		if item.Level > 2 {
			data = append(data, item)
		}
	}

	// 按稀有度高低排序
	sort.Sort(arkRecJSONSlice(data))
	// 处理数据
	charTagSum := 0
	for _, character := range data {
		if character.Hidden {
			continue
		}

		tags := character.Tags
		tags = append(tags, character.Sex+"性干员")
		tags = append(tags, character.Type+"干员")
		// 将干员存入hrData.character中
		ret.characters = append(ret.characters, characterT{
			name: character.Name,
			rare: character.Level,
		})
		index := len(ret.characters) - 1
		for _, tag := range tags {
			// 该tag是否存在于map中
			_, ok := ret.data[tag]
			if !ok {
				ret.data[tag] = []int{index}
				tagList += tag
			} else {
				ret.data[tag] = append(ret.data[tag], index)
			}
		}
		charTagSum += len(tags)
	}
	ret.avgCharTag = float64(charTagSum) / float64(len(ret.data))
	return ret
}

func getCombination(tags []string) []resultT {
	// 获取所有排序组合，去掉不需要的4，5
	combs := combinations.All(tags)
	var temp [][]string
	for _, val := range combs {
		if len(val) > 3 {
			continue
		}
		temp = append(temp, val)
	}
	combs = temp

	var result []resultT
	for _, comb := range combs {
		if len(comb) == 1 && comb[0] == "女性干员" {
			continue
		}

		// 取出需要tag的index
		var need [][]int
		for _, tag := range comb {
			need = append(need, hrData.data[tag])
		}

		// 取交集
		var chars []int
		switch len(need) {
		case 1:
			chars = need[0]
		case 2:
			chars = intersect(need[0], need[1])
		case 3:
			chars = intersect(need[0], intersect(need[1], need[2]))
		}

		// 删除6星干员
		ifGJZS := false
		for _, val := range comb {
			if val == "高级资深干员" {
				ifGJZS = true
				break
			}
		}
		var charsResult []int
		if !ifGJZS {
			ifPrint := true
			for index := range chars {
				if hrData.characters[index].rare == 3 {
					ifPrint = false
					break
				}
			}

			if !ifPrint {
				for index := range chars {
					if hrData.characters[index].rare != 6 {
						charsResult = append(charsResult, index)
					}
				}
			}
		} else {
			charsResult = chars
		}
		if len(charsResult) == 0 {
			continue
		}

		result = append(result, resultT{
			comb:  comb,
			chars: charsResult,
		})
	}

	return result
}
func intersect(nums1 []int, nums2 []int) []int {
	if nums1 == nil || nums2 == nil {
		return []int{}
	}
	sort.Ints(nums1)
	sort.Ints(nums2)
	var x = 0
	var y = 0
	var result []int
	for {
		if x < len(nums1) && y < len(nums2) {
			if nums1[x] == nums2[y] {
				result = append(result, nums1[x])
				x++
				y++
			} else if nums1[x] > nums2[y] {
				y++
			} else {
				x++
			}
		} else {
			break
		}

	}
	return result
}

func initHrData() {
	hrData = pullData()
}
func newHrData() hrDataT {
	d := make(map[string][]int)
	return hrDataT{
		data: d,
	}
}
func getChar(i int) string {
	return hrData.characters[i].name
}

// 重写method
func (a arkRecJSONSlice) Len() int {
	return len(a)
}
func (a arkRecJSONSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a arkRecJSONSlice) Less(i, j int) bool {
	return a[j].Level < a[i].Level
}
