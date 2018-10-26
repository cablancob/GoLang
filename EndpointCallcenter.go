package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

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

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "taksio"
)

type client_consult struct {
	Date_from string `json:"date_from"`
	Date_to   string `json:"date_to"`
}

func bd_consult(date_from string, date_to string) string {
	result := ""

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM callcenter.consultdatacallcenter('" + date_from + "','" + date_to + "')")
	if err != nil {
		panic(err)
	}

	rows.Next()
	err = rows.Scan(&result)
	if err != nil {
		panic(err)
	}

	return result

}

func request(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/callcenter" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		Block{
			Try: func() {
				body, err := ioutil.ReadAll(r.Body)

				if err != nil {
					panic(err)
					http.Error(w, "404 not found.", http.StatusNotFound)
					return
				}

				var json_client = client_consult{}
				json.Unmarshal(body, &json_client)
				fmt.Fprintf(w, "%s", bd_consult(json_client.Date_from, json_client.Date_to))
			}, Catch: func(e Exception) {
				fmt.Printf("Error %v\n", e)
			},
		}.Do()
	default:
		fmt.Fprintf(w, "Solo metodo POST es soportado")
	}
}

func main() {

	http.HandleFunc("/callcenter", request)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal(err)
	}

}
