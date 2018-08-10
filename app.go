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

var db *sql.DB

func init() {
	var err error
	db, err = connect()
	if err != nil {
		db = nil
		log.Println("unable to connect to db")
	}
	err = db.Ping()
	if err != nil {
		db = nil
		log.Println("db ping error", err)
	}
}

func connect() (*sql.DB, error) {
	con := fmt.Sprintf("%s:%s@tcp(%s:3306)/mysql", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))
	return sql.Open("mysql", con)
}

func connectSsl() (*sql.DB, error) {
	rootCAs := x509.NewCertPool()

	pem, err := ioutil.ReadFile("rds-combined-ca-bundle.pem")
	if err != nil {
		panic(err)
	}
	if ok := rootCAs.AppendCertsFromPEM(pem); !ok {
		panic("unable to append certs")
	}

	mysqlHost := os.Getenv("MYSQL_HOST")
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:    rootCAs,
		ServerName: mysqlHost,
	})
	con := fmt.Sprintf("%s:%s@tcp(%s:3306)/mysql?tls=custom", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), mysqlHost)
	return sql.Open("mysql", con)
}

func main() {

	fmt.Println("MAIN ENTRY")

	router := web.New(Context{})
	router.Get("/", func(w web.ResponseWriter, req *web.Request) {
		w.Write([]byte("hello world 3"))
	})

	router.Get("/mysql", func(w web.ResponseWriter, req *web.Request) {

		w.Header().Set("Content-Type", "text/plain")

		if db==nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("no db"))
			return
		}

		// http://go-database-sql.org/
		rows, err := db.Query("select User from user")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer rows.Close()

		var userName string
		for rows.Next() {
			err := rows.Scan(&userName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(userName))
			w.Write([]byte("\n"))
		}
		err = rows.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

	})
	http.ListenAndServe(":80", router)
}
