package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//DESARROLLO
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "taksio"
)

type rider_data struct {
	Name string `json:"name"`
}

var url string = "http://registry.taksio.net/rider/"

func GetNameRider(uid string) string {
	result := ""

	println(url + uid)
	rs, err := http.Get(url + uid)
	if err != nil {
		panic(err)
	}

	body, err2 := ioutil.ReadAll(rs.Body)
	if err2 != nil {
		panic(err)
	}

	if string(body) != "" {
		ver := new(rider_data)
		err3 := json.Unmarshal([]byte(strings.TrimSpace(string(body))), &ver)
		if err3 != nil {
			panic(err)
		}
		result = ver.Name
	} else {
		result = "SIN NOMBRE"
	}

	return result

}

func main() {

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err1 := sql.Open("postgres", dbinfo)
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()

	/*	file, err := os.Open("/home/taksioweb-01/TKS_ALL - TKS_ALL.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if err2 := scanner.Err(); err2 != nil {
			db.Close()
			log.Fatal(err2)
		}

			_, err3 := db.Exec("TRUNCATE TABLE tks RESTART IDENTITY")
			if err3 != nil {
				db.Close()
				panic(err3)
			}

			for scanner.Scan() {
				datos := strings.Split(scanner.Text(), ",")
				if strings.TrimSpace(datos[0]) != "ID" {
					_, err4 := db.Exec("INSERT INTO tks(id_tks, rider, driver, fecha_creacion, fecha_expiracion) VALUES ('" + strings.TrimSpace(datos[0]) + "','" + strings.TrimSpace(datos[1]) + "','" + strings.TrimSpace(datos[2]) + "','" + strings.TrimSpace(datos[11]) + "','" + strings.TrimSpace(datos[9]) + "')")
					if err4 != nil {
						db.Close()
						panic(err4)
					}
				}
			}*/

	rider := ""
	nombre := ""
	fecha_c := ""
	fecha_v := ""
	referencia := ""
	entidad := ""
	tks := ""
	//rows, err10 := db.Query("SELECT rider,count(*), (select coalesce(to_char(max(date),'YYYY-MM-DD'),'0000') from financial.payment_service_rider where user_uuid = A.rider AND status = 'cleared')  FROM tks_full A WHERE driver = '' AND id_tks not in (SELECT id_tks FROM tks) GROUP BY rider ORDER BY count(*) DESC")
	rows, err10 := db.Query("select A.rider, coalesce(B.reference,''), coalesce(B.entity,''), A.C, A.D, A.F from (SELECT rider, to_date(fecha_creacion,'yyyy-mm-dd') as C, to_date(fecha_expiracion,'yyyy-mm-dd') as D, count(*) as F FROM tks_full WHERE driver = '' and id_tks not in (select id_tks from tks) GROUP BY rider, to_date(fecha_creacion,'yyyy-mm-dd'), to_date(fecha_expiracion,'yyyy-mm-dd') order by to_date(fecha_creacion,'yyyy-mm-dd') desc) as A LEFT JOIN (select user_uuid, date as E, reference, entity  from financial.payment_service_rider WHERE status = 'cleared' group by user_uuid, date, reference, entity) as B on A.rider = B.user_uuid and A.C =B.E order by A.rider, A.C asc")
	if err10 != nil {
		panic(err10)
	}

	f, _ := os.Create("/home/taksioweb-01/TkValidos.csv")
	f.WriteString("RIDER,REFERENCIA,ENTIDAD,FECHA CREACION,FECHA VENCIMOENTO,TKS\n")

	for rows.Next() {
		err := rows.Scan(&rider, &referencia, &entidad, &fecha_c, &fecha_v, &tks)
		if err != nil {
			panic(err)
		}
		nombre = GetNameRider(rider)
		t, err := time.Parse("2006-01-02T15:04:05Z", fecha_c)
		t1, err := time.Parse("2006-01-02T15:04:05Z", fecha_v)
		f.WriteString(nombre + "," + referencia + "," + entidad + "," + t.Format("2006-01-02") + "," + t1.Format("2006-01-02") + "," + tks + "\n")
		println(nombre + "," + referencia + "," + entidad + "," + t.Format("2006-01-02") + "," + t1.Format("2006-01-02") + "," + tks)
	}

	db.Close()

}
