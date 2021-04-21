package main

import (
	"encoding/json"
	"fmt"
	"github.com/reujab/wallpaper"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var site = "http://service.aibizhi.adesk.com/v1/wallpaper/category"
var ua = "(picasso,170,windows)"
var catName = "girl"
var count int
var ch chan string
var paperList chan []interface{}

func loadInf(r interface{}, key string) (ret interface{}) {
	m, isOk := r.(map[string]interface{})
	if isOk {
		for k, v := range m {
			if k == key {
				ret = v
			}
		}
	}
	return ret
}

func mapValue(r interface{}, key string) (value interface{}) {
	m, isOk := r.(map[string]interface{})
	if isOk {
		for k, v := range m {
			if k == key {
				value = v
				break
			}
		}
	}
	return value
}

func reqSite(url string) {
	text := ""
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err)
		return
	}
	req.Header.Add("User-Agent", ua)
	resp, err := httpClient.Do(req)
	if err != nil {
		println(err)
		return
	}
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err)
		return
	}

	var r interface{}
	_ = json.Unmarshal(s, &r)

	res := loadInf(r, "res")
	if res != nil {
		category := loadInf(res, "category")
		if category != nil {
			cats, isOk := category.([]interface{})
			if isOk {
				for i := range cats {
					ename := loadInf(cats[i], "ename")
					if ename == nil {
						continue
					}
					name := mapValue(cats[i], "ename")
					if name == catName {
						id := mapValue(cats[i], "id")
						t, _ := mapValue(cats[i], "count").(float64)
						count = int(t)
						fmt.Printf("find %d papers\n", count)
						text = id.(string)
						break
					}
				}
			}

		}
	}
	ch <- text
}

func getPageURL(id string, page int) {
	_ = ioutil.WriteFile("log.txt", []byte(strconv.Itoa(page)), 0755)
	url := fmt.Sprintf("%s/%s/wallpaper?skip=%d", site, id, page*20)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", ua)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	var r interface{}
	_ = json.Unmarshal(body, &r)

	res := loadInf(r, "res")
	if res != nil {
		wallpapers := loadInf(res, "wallpaper")
		if wallpapers != nil {
			papers, _ := wallpapers.([]interface{})
			paperList <- papers
			return
		}
	}
}
func setPaper(papers []interface{}) {
	for i := range papers {
		img := loadInf(papers[i], "img")
		if img == nil {
			return
		}
		imgURL := mapValue(papers[i], "img").(string)
		fmt.Printf("set paper with %s\n", imgURL)
		err := wallpaper.SetFromURL(imgURL)
		if err != nil {
			println(err)
			return
		}
		time.Sleep(15 * time.Minute)
	}
}

func onStart() (page int) {
	s, err := ioutil.ReadFile("log.txt")
	if err != nil {
		_ = ioutil.WriteFile("log.txt", []byte("0"), 0755)
		return 0
	}
	i, _ := strconv.Atoi(string(s))
	return i
}

func main() {
	ch = make(chan string)
	paperList = make(chan []interface{})
	defer close(ch)
	page := onStart()
getSiteLoop:
	for {
		go reqSite(site)
		var id string
		select {
		case id = <-ch:
			if id == "" {
				time.Sleep(5 * time.Minute)
				continue getSiteLoop
			}
		getPageLoop:
			for ; page <= count/20; page++ {
				go getPageURL(id, page)
				select {
				case paperURL := <-paperList:
					setPaper(paperURL)
				case <-time.After(20 * time.Second):
					break getPageLoop
				}
			}
		case <-time.After(20 * time.Second):
			time.Sleep(5 * time.Minute)
		}
		page = 0
		_ = ioutil.WriteFile("log.txt", []byte("0"), 0755)
	}
}
