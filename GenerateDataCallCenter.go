package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
)

//DESARROLLO
/*const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "taksio"
)*/

//PRODUCCION
const (
	host     = "experimental.taksio.net"
	port     = 5432
	user     = "postgres"
	password = "3xp3r1m3nt4L-t4ks10-2018"
	dbname   = "taksio"
)

func isDAte(s string) bool {
	bolean := true
	_, err := time.Parse("02-01-2006", s)
	if err != nil {
		bolean = false
	} else {
		bolean = true
	}
	return bolean
}

func isDAte2(s string) bool {
	bolean := true
	_, err := time.Parse("02/01/2006", s)
	if err != nil {
		bolean = false
	} else {
		bolean = true
	}
	return bolean
}

func main() {
	directory := "/home/taksioweb-01/CallCenter/"

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err1 := db.Exec("TRUNCATE TABLE callcenter.callcenterdata RESTART IDENTITY")
	if err1 != nil {
		db.Close()
		panic(err1)
	}

	files, err2 := ioutil.ReadDir(directory)
	if err2 != nil {
		db.Close()
		panic(err2)
	}
	for _, file := range files {
		excelFileName := directory + file.Name()
		xlFile, err3 := xlsx.OpenFile(excelFileName)
		if err3 != nil {
			db.Close()
			panic(err3)
		}
		for _, sheet := range xlFile.Sheets {
			if isDAte(strings.TrimSpace(sheet.Name)) {
				for _, row := range sheet.Rows {
					if strings.TrimSpace(row.Cells[1].String()) != "" && isDAte2(strings.TrimSpace(row.Cells[1].String())) {
						if strings.ToUpper(strings.TrimSpace(row.Cells[6].String())) != "APP" && strings.TrimSpace(row.Cells[2].String()) != "" {
							driver := strings.ToUpper(row.Cells[2].String())
							driver = strings.Replace(driver, "Á", "A", -1)
							driver = strings.Replace(driver, "É", "E", -1)
							driver = strings.Replace(driver, "Í", "I", -1)
							driver = strings.Replace(driver, "Ó", "O", -1)
							driver = strings.Replace(driver, "Ú", "U", -1)
							//println(row.Cells[1].String() + " " + row.Cells[2].String() + " " + row.Cells[3].String() + " " + row.Cells[4].String() + " " + row.Cells[5].String() + " " + row.Cells[5].String() + " " + row.Cells[9].String())
							//println(row.Cells[6].String())
							_, err4 := db.Exec("INSERT INTO callcenter.callcenterdata(date_ride, driver, origin, destination, rider, tks) VALUES (to_date('" + strings.TrimSpace(row.Cells[1].String()) + "','DD/MM/YYYY'), '" + strings.TrimSpace(driver) + "', '" + strings.TrimSpace(row.Cells[3].String()) + "', '" + strings.TrimSpace(row.Cells[4].String()) + "', 'a8a90dd8-02fd-4f86-815b-90cab0435d46', '" + strings.TrimSpace(row.Cells[9].String()) + "');")
							if err4 != nil {
								db.Close()
								panic(err4)
							}
						}
					}

				}
			}
		}
	}
	_, err5 := db.Exec("UPDATE callcenter.callcenterdata A SET driver = B.uuid FROM callcenter.userdata B WHERE A.driver = B.driver")
	if err5 != nil {
		db.Close()
		panic(err5)
	}

	rows, err10 := db.Query("SELECT driver from callcenter.callcenterdata WHERE LENGTH(driver) < 36 GROUP BY driver")
	if err10 != nil {
		panic(err10)
	}

	result := ""

	for rows.Next() {
		err11 := rows.Scan(&result)
		if err11 != nil {
			panic(err11)
		}
		println(result)
	}

	rows.Close()

	db.Close()
}
