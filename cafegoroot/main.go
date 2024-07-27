package main

import (
	"html/template"
	"log"
	"net/http"
)

type IndexPageData struct {
	Username string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	sampleUsername := "Matthew"
	samplePageData := IndexPageData{Username: sampleUsername}
	err = tmpl.Execute(w, samplePageData)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":5000", nil)
}
