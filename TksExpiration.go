package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

type tks struct {
	Id      string `json:"id"`
	Expires string `json:"expires"`
}

var url string = "https://promotions.devel.taksio.net/coupon/get"

func ListTks() {

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, _ := sql.Open("postgres", dbinfo)
	defer db.Close()

	rs, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}

	rs.Header.Set("authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpZCI6IjM2N2ZhZTYxLWM0MTItNDU2ZS1iNzM2LWY5YTY5YzNmZWQzMiIsImlhdCI6MTU0MzM0MjM5MX0.rSXpSQ_RczdQtmb96flMG5mhpqe4-yXpIGQGJ7C6j4-AUVQ6M1RfLVi6eLm6ei7DTOYImt1lHbgxUwx7raQTd-AhZpVV7WyHLYJK9rTwH5DsNt0G7SOh0IA9N_27MjaEcvUXY7-72xFObpsDZfoii7G7tsbw8mdIiAB6X0djA2uZW0a_U3fFjCCBkckIBC0WMWTsDdlMGRDgCwkfuY6SqwcWpfS6wKRV3D14TT7uOzzP8NqCOpZGVEVWGkqJSuJMoq6PVNXpDns6j1LxQA3Uc3V_AkM4RcNuifY1No8EO9Ise-rAuqDsMxRr8wHP4VPmiAP5wFuRNBRTk0X-nPVmWw")
	client := &http.Client{}
	resp, err := client.Do(rs)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println("response Body:", string(body))

	//ver := new(tks)
	var tmp []tks
	err3 := json.Unmarshal([]byte(strings.TrimSpace(string(body))), &tmp)
	if err3 != nil {
		panic(err3)
	}

	for _, element := range tmp {
		fmt.Println(element.Id + " " + element.Expires)
		db.Exec("INSERT INTO public.tks_despues(id, expires)VALUES ('" + element.Id + "', '" + element.Expires + "');")
	}

}

func main() {
	ListTks()
}
