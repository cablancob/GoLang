package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type rider_data struct {
	Name string `json:"name"`
}

var url string = "http://registry.taksio.net/rider/"
var urlPago = "https://promotions.taksio.net/coupon/create/user/"

func GetNameRider(uid string) string {
	result := ""

	rs, err := http.Get(url + uid)
	if err != nil {
		panic(err)
	}

	body, err2 := ioutil.ReadAll(rs.Body)
	if err2 != nil {
		panic(err)
	}

	ver := new(rider_data)
	err3 := json.Unmarshal([]byte(strings.TrimSpace(string(body))), &ver)
	if err3 != nil {
		panic(err)
	}

	result = ver.Name
	return result

}

func main() {
	file, err := os.Open("/home/taksioweb-01/nuevo 5.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	f, err10 := os.Create("/home/taksioweb-01/RidersTks.csv")
	if err10 != nil {
		log.Fatal(err10)
	}
	f.WriteString("UID,NOMBRE,TKS ACTUALES,CANTIDAD DE TKS A SUMAR\n")
	scanner := bufio.NewScanner(file)
	if err2 := scanner.Err(); err2 != nil {
		log.Fatal(err2)
	}

	for scanner.Scan() {
		datos := strings.Split(scanner.Text(), "|")
		uid := strings.TrimSpace(datos[0])
		tks := strings.TrimSpace(datos[1])
		nuevoTks, err30 := strconv.Atoi(tks)
		if err30 != nil {
			panic(err30)
		}
		nuevoTks = nuevoTks*5 - nuevoTks
		nombre := strings.TrimSpace(GetNameRider(uid))

		jsonStr := []byte(`{"reference" : "123456","text" : "marketing/UX purposes","discount" : 100,"category" : "taxi","expires" : 30,"amount" : ` + strconv.Itoa(nuevoTks) + `}`)
		println(string(jsonStr))
		/*req, err := http.NewRequest("POST", urlPago+uid, bytes.NewBuffer(jsonStr))
		req.Header.Add("authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjU0NTM0MmEwLWExOTktNGJhOC1iZjU5LWY4ZWRiYTE4Mjk3MiIsImlwIjoiMTI3LjAuMC4xIiwiaWF0IjoxNTA5NzU0ODQzfQ.lrJSPmkRs3ka8JckL9WZSPTbvl018ZB8EPiRP9AvPS2w9D0La-NvoL7tTDhswAFY6_Od3Ot-jDlcsN9zVuiSmI02ZZ5paUklLDYV9WjMQO0MLbGU8GiB6IwsXYTZqJ8-dn-B46JVhMrZbFLLTfETss6e4r-hzJsK4hXmQObWAfBBbW3QRQXMlAIrbFURhwZdafyd7o7BUdficlb4Sxtl473IypDu9N5RS0gngvmQFKm5nRDgCIQWDHjy20Mr2U8Ola84x8BzS5GswBU44p2z1vjFCZ_UBFauC2Z_8HA-U9UbGq2Q6EoFX3GV91SyjZwFRw40bLcjf-DkdAdE2g22vg")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()*/
		println(uid + ";" + nombre + ";" + tks + ";" + strconv.Itoa(nuevoTks))
		f.WriteString(uid + "," + nombre + "," + tks + "," + strconv.Itoa(nuevoTks) + "\n")

	}

}
