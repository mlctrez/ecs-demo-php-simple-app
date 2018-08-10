package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"crypto/tls"
	"crypto/x509"
	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/web"
	"io/ioutil"
)

type Context struct{}

func main() {

	router := web.New(Context{})
	router.Get("/", func(w web.ResponseWriter, req *web.Request) {
		w.Write([]byte("hello world 3"))
	})

	router.Get("/mysql", func(w web.ResponseWriter, req *web.Request) {

		rootCAs := x509.NewCertPool()

		pem, err := ioutil.ReadFile("/rds-combined-ca-bundle.pem")
		if err != nil {
			panic(err)
		}
		if ok := rootCAs.AppendCertsFromPEM(pem); !ok {
			panic(err)
		}

		mysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCAs})

		con := fmt.Sprintf("%s:%s@%s:3306/mysql?tls=custom", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))

		db, err := sql.Open("mysql", con)
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err)
		}

		defer db.Close()

	})
	http.ListenAndServe(":80", router)
}
