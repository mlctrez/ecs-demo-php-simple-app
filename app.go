package main

import (
	"github.com/gocraft/web"
	"net/http"
)

type Context struct{}

func main() {

	router := web.New(Context{})
	router.Get("/", func(w web.ResponseWriter, req *web.Request) {
		w.Write([]byte("hello world 2"))
	})
	http.ListenAndServe(":80", router)
}
