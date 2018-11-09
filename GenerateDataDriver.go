package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "taksio"
)

func main() {
	file, err := os.Open("/home/taksioweb-01/UUID_DRIVER.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err1 := sql.Open("postgres", dbinfo)
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()

	scanner := bufio.NewScanner(file)
	if err2 := scanner.Err(); err2 != nil {
		db.Close()
		log.Fatal(err2)
	}

	_, err3 := db.Exec("TRUNCATE TABLE callcenter.userdata RESTART IDENTITY")
	if err3 != nil {
		db.Close()
		panic(err3)
	}

	for scanner.Scan() {
		datos := strings.Split(scanner.Text(), ",")
		if strings.TrimSpace(datos[0]) != "nombre" {
			_, err4 := db.Exec("INSERT INTO callcenter.userdata(driver, uuid) VALUES ('" + strings.TrimSpace(datos[0]) + "', '" + strings.TrimSpace(datos[1]) + "')")
			if err4 != nil {
				db.Close()
				panic(err4)
			}
		}
	}

	db.Close()

}
