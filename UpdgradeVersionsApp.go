package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type app_version struct {
	Version string `json:"version"`
}

//variables produccion
//var url_rider string = "http://driver.taksio.net/taksio/private/rider/version/"
//var url_driver string = "http://driver.taksio.net/taksio/private/driver/version/"

//variables desarrollo
var url_rider string = "http://localhost:8008/rider/version/"
var url_driver string = "http://localhost:8009/driver/version/"

func main() {

	rider := "https://play.google.com/store/apps/details?id=com.taksio.apps.client"
	driver := "https://play.google.com/store/apps/details?id=com.taksio.apps.partner"
	app_rider := "RIDER"
	app_driver := "DRIVER"

	for {

		versionPlayStore_rider := ScrapingHtml(rider)
		versionPlayStore_driver := ScrapingHtml(driver)

		versionBd_rider := GetActualVersion(app_rider)
		versionBd_driver := GetActualVersion(app_driver)

		println("Play Rider: ", versionPlayStore_rider)
		println("Play Driver: ", versionPlayStore_driver)

		println("BD Rider: ", versionBd_rider)
		println("BD Driver: ", versionBd_driver)

		if versionPlayStore_rider != versionBd_rider {
			UpdateVersion(app_rider, versionPlayStore_rider)
		}

		if versionPlayStore_driver != versionBd_driver {
			UpdateVersion(app_driver, versionPlayStore_driver)
		}

		fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())
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
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
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
		req, err := http.NewRequest("POST", url_driver, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

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
		req, err := http.NewRequest("POST", url_driver, bytes.NewBuffer(jsonStr))
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
