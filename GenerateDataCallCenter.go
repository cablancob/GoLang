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
	host     = "192.168.1.140"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "taskio"
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
	_, err1 := time.Parse("02/1/2006", s)
	if err1 != nil {
		bolean = false
	} else {
		bolean = true
	}
	return bolean
}

func main() {

	//fecha_actual := time.Now()

	directory := "/home/taksioweb-01/CallCenter/"

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	/*f, err10 := os.Create("/home/taksioweb-01/archivos.txt")
	if err10 != nil {
		log.Fatal(err10)
	}*/

	files, err2 := ioutil.ReadDir(directory)
	if err2 != nil {
		db.Close()
		panic(err2)
	}
	for _, file := range files {
		excelFileName := directory + file.Name()
		//var nombre_archivo = file.Name()

		xlFile, err3 := xlsx.OpenFile(excelFileName)
		if err3 != nil {
			panic(err3)
		}
		for _, sheet := range xlFile.Sheets {
			if isDAte(strings.TrimSpace(sheet.Name)) {
				s := strings.Split(sheet.Name, "-")
				date_ride := strings.TrimSpace(strings.TrimSpace(s[2]) + "-" + strings.TrimSpace(s[1]) + "-" + strings.TrimSpace(s[0]))
				fmt.Println(date_ride)
				fmt.Println("")
				//f.WriteString(date_ride + "\n")
				for _, row := range sheet.Rows {
					//if strings.TrimSpace(row.Cells[1].String()) == "Fecha" {
					if strings.TrimSpace(row.Cells[1].String()) != "" && isDAte2(strings.TrimSpace(row.Cells[1].String())) && strings.TrimSpace(row.Cells[2].String()) != "" {
						driver := strings.TrimSpace(row.Cells[2].String())
						origin := strings.TrimSpace(row.Cells[3].String())
						destination := strings.TrimSpace(row.Cells[4].String())
						rider := strings.TrimSpace(row.Cells[5].String())
						tp_solicitud := strings.TrimSpace(row.Cells[6].String())
						tp_servicio := strings.TrimSpace(row.Cells[7].String())
						nro_transaccion := strings.TrimSpace(row.Cells[8].String())
						mto_tiempo_espera := strings.TrimSpace(strings.Replace(strings.TrimSpace(row.Cells[9].String()), "Bs.", "", -1))
						if mto_tiempo_espera == "" || mto_tiempo_espera == "-" || mto_tiempo_espera == "#REF!" {
							mto_tiempo_espera = "0"
						}
						cant_tks := strings.TrimSpace(row.Cells[10].String())
						if cant_tks == "" {
							cant_tks = "0"
						}
						mto_pago_conductor := strings.TrimSpace(strings.Replace(strings.TrimSpace(row.Cells[13].String()), "Bs.", "", -1))
						if mto_pago_conductor == "" || mto_pago_conductor == "-" || mto_pago_conductor == "#REF!" {
							mto_pago_conductor = "0"
						}
						de_estatus := strings.TrimSpace(row.Cells[14].String())
						mto_pasajero := strings.TrimSpace(strings.Replace(strings.Replace(strings.TrimSpace(row.Cells[16].String()), "Bs.", "", -1), "%", "", -1))
						if mto_pasajero == "" || mto_pasajero == "-" || mto_pasajero == "#REF!" {
							mto_pasajero = "0"
						}

						fmt.Print(date_ride + " | " + driver + " | " + origin + " | " + destination + " | " + rider + " | " + tp_solicitud + " | " + tp_servicio + " | " + nro_transaccion + " | " + mto_tiempo_espera + " | " + cant_tks + " | " + mto_pago_conductor + " | " + de_estatus + " | " + mto_pasajero)
						sql := "SELECT * FROM rider_control.migration('" + date_ride + "','" + driver + "','" + origin + "','" + destination + "','" + rider + "','" + tp_solicitud + "','" + tp_servicio + "','" + nro_transaccion + "','" + mto_tiempo_espera + "','" + cant_tks + "','" + mto_pago_conductor + "','" + de_estatus + "','" + mto_pasajero + "');"
						fmt.Println("")
						fmt.Println(sql)
						_, err4 := db.Exec(sql)
						if err4 != nil {
							db.Close()
							panic(err4)
						}
						/*texto := ""
						for _, cell := range row.Cells {
							text := cell.String()
							fmt.Print(text + "|")
							texto = texto + text + "|"
						}
						texto = texto + "\n"
						f.WriteString(texto)*/
						fmt.Println("")
					}
				}
			}
		}

	}
	/*	_, err5 := db.Exec("UPDATE callcenter.callcenterdata A SET driver = B.uuid FROM callcenter.userdata B WHERE A.driver = B.driver")
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

	rows.Close()*/

	db.Close()
}
