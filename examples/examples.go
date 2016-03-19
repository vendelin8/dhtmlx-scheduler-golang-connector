package main

import (
	"fmt"

	//for MySQL, use the following line:
	_ "github.com/go-sql-driver/mysql"

	//for SQLite, use the following line:
	_ "github.com/mattn/go-sqlite3"

	dsgc "github.com/vendelin8/dhtmlx-scheduler-golang-connector"
	"net/http"
	"os"
	"path"
)

func main() {
	currentDir, err := os.Getwd() //current working directory
	if err != nil {
		panic(err)
	}

	//for MySQL, use the following line, with your credentials:
	//err = dsgc.Open("mysql", "username:password@tcp(host:port)/test?parseTime=true", "")

	//for SQLite, use the following line:
	err = dsgc.Open("sqlite3", path.Join(currentDir, "test"), "")

	if err != nil {
		panic(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path.Join(currentDir, "static")))))

	bind := "0.0.0.0:1212"
	fmt.Printf("listening on %s...\n", bind)
	err = http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

