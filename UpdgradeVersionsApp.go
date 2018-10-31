package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type app_version struct {
	Version string `json:"version"`
}

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

//variables produccion
var url_rider string = "http://driver.taksio.net/taksio/private/rider/version"
var url_driver string = "http://driver.taksio.net/taksio/private/driver/version"

func main() {

	rider := "https://play.google.com/store/apps/details?id=com.taksio.apps.client"
	driver := "https://play.google.com/store/apps/details?id=com.taksio.apps.partner"
	app_rider := "RIDER"
	app_driver := "DRIVER"

	for {

		Block{
			Try: func() {
				versionPlayStore_rider := ScrapingHtml(rider)
				versionPlayStore_driver := ScrapingHtml(driver)

				versionBd_rider := GetActualVersion(app_rider)
				versionBd_driver := GetActualVersion(app_driver)

				if versionPlayStore_rider != versionBd_rider && strings.TrimSpace(versionPlayStore_rider) != "" {
					UpdateVersion(app_rider, versionPlayStore_rider)
				}

				if versionPlayStore_driver != versionBd_driver && strings.TrimSpace(versionPlayStore_driver) != "" {
					UpdateVersion(app_driver, versionPlayStore_driver)
				}

				t := time.Now()
				fmt.Println(t.Format("2006-01-02 15:04:05") + ", Play Rider: " + versionPlayStore_rider + ", Play Driver: " + versionPlayStore_driver + ", BD Rider: " + versionBd_rider + ", BD Driver: " + versionBd_driver)
			}, Catch: func(e Exception) {
				fmt.Printf("Error %v\n", e)
			},
		}.Do()
		time.Sleep(300 * time.Second)
	}

}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func ScrapingHtml(s string) string {
	// Request the HTML page.
	var result string

	res, err := http.Get(s)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		println("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	// Find the review items
	doc.Find("div span .htlgb").Each(func(i int, s *goquery.Selection) {
		if isNumeric(strings.Replace(s.Text(), ".", "", -1)) {
			result = s.Text()
		}

	})

	return strings.TrimSpace(result)
}

func GetActualVersion(app string) string {
	result := ""

	if app == "RIDER" {
		rs, err := http.Get(url_rider)
		if err != nil {
			panic(err)
		}

		body, err2 := ioutil.ReadAll(rs.Body)
		if err2 != nil {
			panic(err)
		}

		ver := new(app_version)
		err3 := json.Unmarshal([]byte(strings.TrimSpace(string(body))), &ver)
		if err3 != nil {
			panic(err)
		}
		result = ver.Version

	}

	if app == "DRIVER" {
		rs, err := http.Get(url_driver)
		if err != nil {
			panic(err)
		}

		body, err2 := ioutil.ReadAll(rs.Body)
		if err2 != nil {
			panic(err)
		}
		ver := new(app_version)
		err3 := json.Unmarshal([]byte(strings.TrimSpace(string(body))), &ver)
		if err3 != nil {
			panic(err)
		}
		result = ver.Version

	}

	return strings.TrimSpace(result)
}

func UpdateVersion(app string, actualversion string) {
	if app == "RIDER" {
		var jsonStr = []byte(`{"operation": 1, "name_app_p": "RIDER" , "Version_app_p": "` + actualversion + `"}`)
		req, err := http.NewRequest("POST", url_rider+"?operation=1&name_app_p="+app+"&version_app_p="+actualversion+"", bytes.NewBuffer(jsonStr))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		println("SE ACTUALIZO LA VERSION DEL RIDER, HTTP: ", resp.Status)

	}

	if app == "DRIVER" {
		var jsonStr = []byte(`{"operation": 1, "name_app_p": "DRIVER" , "Version_app_p": "` + actualversion + `"}`)
		req, err := http.NewRequest("POST", url_driver+"?operation=1&name_app_p="+app+"&version_app_p="+actualversion+"", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		println("SE ACtUALIZO LA VERSION DEL DRIVER, HTTP: ", resp.Status)
	}
}
