package main

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/ecoball/go-ecoball/explorer/data"
	"github.com/ecoball/go-ecoball/explorer/onlooker"
)

type WebHandle func(w http.ResponseWriter, r *http.Request)

type webserver struct {
	url2handle map[string]WebHandle
}

func (this *webserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	path := r.URL.String()
	switch path {
	case "/":
		t := template.Must(template.ParseFiles("./root.html"))
		t.Execute(w, data.PrintBlock())

	default:
		err = errors.New("unrecognized transaction type")
	}

	if err != nil {
		http.Error(w, "error 500: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	go onlooker.Bystander()
	http.ListenAndServe(":8080", &webserver{})
}
