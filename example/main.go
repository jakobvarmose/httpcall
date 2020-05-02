package main

import (
	"log"
	"net/http"
	"github.com/jakobvarmose/httpcall"
)

func main() {
	fileServer := http.FileServer(http.Dir("."))

	app := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Println(req.Method, req.URL.String())
			res := httpcall.Call(fileServer, req)
			if res.StatusCode == http.StatusNotFound {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Your file was not found\n"))
				return
			}
			httpcall.Write(w, res)
		}),
	}

	log.Println("listening")
	err := app.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
