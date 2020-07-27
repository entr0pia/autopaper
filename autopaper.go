package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"github.com/reujab/wallpaper"
)

var site = "http://service.aibizhi.adesk.com/v1/wallpaper/category"
var ua = "(picasso,170,windows)"
var catName = "girl"
var count int

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

func reqSite(url string) (text string, err error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", ua)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r interface{}
	json.Unmarshal(s, &r)

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
						t,_ := mapValue(cats[i], "count").(float64)
						count =int(t)
						fmt.Printf("find %d papers\n", count)
						text = id.(string)
						break
					}
				}
			}

		}
	}
	return text, nil
}

func setPaper(id string, page int) {
	url := fmt.Sprintf("%s/%s/wallpaper?skip=%d", site, id, page*20)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", ua)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	var r interface{}
	json.Unmarshal(body, &r)

	res := loadInf(r, "res")
	if res != nil {
		wallpapers := loadInf(res, "wallpaper")
		if wallpapers != nil {
			papers, _ := wallpapers.([]interface{})
			for i := range papers {
				img := loadInf(papers[i], "img")
				if img == nil {
					continue
				}
				imgURL := mapValue(papers[i], "img").(string)
				fmt.Printf("set paper with %s\n",imgURL)
				wallpaper.SetFromURL(imgURL)
				time.Sleep(15 * time.Minute)
			}
		}
	}
}

func onStart() (page int) {
	s, err := ioutil.ReadFile("log.txt")
	if err != nil {
		ioutil.WriteFile("log.txt", []byte("0"), 0755)
		return 0
	}
	i, _ := strconv.Atoi(string(s))
	return i
}

func main() {
	page := onStart()
	for {
		id, err := reqSite(site)
		if err != nil {
			time.Sleep(5 * time.Minute)
			continue
		}
		for ; page <= count/20; page++ {
			setPaper(id, page)
			ioutil.WriteFile("log.txt", []byte(strconv.Itoa(page)), 0755)
		}
		page = 0
		ioutil.WriteFile("log.txt", []byte("0"), 0755)
	}
}
