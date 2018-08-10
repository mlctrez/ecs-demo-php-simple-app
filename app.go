package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/web"
)

type Context struct{}

func main() {

	fmt.Println("MAIN ENTRY")

	router := web.New(Context{})
	router.Get("/", func(w web.ResponseWriter, req *web.Request) {
		w.Write([]byte("hello world 3"))
	})

	router.Get("/mysql", func(w web.ResponseWriter, req *web.Request) {

		w.Header().Set("Content-Type", "text/plain")

		var con string

		if false {
			rootCAs := x509.NewCertPool()

			pem, err := ioutil.ReadFile("/rds-combined-ca-bundle.pem")
			if err != nil {
				panic(err)
			}
			if ok := rootCAs.AppendCertsFromPEM(pem); !ok {
				panic(err)
			}

			mysqlHost := os.Getenv("MYSQL_HOST")
			config := &tls.Config{
				RootCAs:    rootCAs,
				ServerName: "",
			}
			mysql.RegisterTLSConfig("custom", config)
			con = fmt.Sprintf("%s:%s@tcp(%s:3306)/mysql?tls=custom", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), mysqlHost)
		} else {
			con = fmt.Sprintf("%s:%s@tcp(%s:3306)/mysql", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))
		}

		db, err := sql.Open("mysql", con)
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err)
		}
		log.Println("after ping")

		// http://go-database-sql.org/
		rows, err := db.Query("select User from user")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var userName string
		for rows.Next() {
			err := rows.Scan(&userName)
			if err != nil {
				panic(err)
			}
			w.Write([]byte(userName))
			w.Write([]byte("\n"))
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		defer db.Close()

	})
	http.ListenAndServe(":80", router)
}
