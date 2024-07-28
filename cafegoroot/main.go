package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type IndexPageData struct {
	Username string
	Products []Product
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	sampleUsername := "Matthew"
	sampleProducts := getProducts()
	samplePageData := IndexPageData{Username: sampleUsername, Products: sampleProducts}
	err = tmpl.Execute(w, samplePageData)
	if err != nil {
		log.Fatal(err)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	// Get the product ID
	reqPath := r.URL.Path
	splitPath := strings.Split(reqPath, "/")
	elemCount := len(splitPath)
	// Do note that this will be a string.
	productId := splitPath[elemCount-1]
	// Need to convert from string to int
	intId, err := strconv.Atoi(productId)
	if err != nil {
		log.Fatal(err)
	}
	// Predeclare a product
	var product Product
	// Check each product for whether it matches the given ID
	for _, p := range getProducts() {
		if p.Id == intId {
			product = p
			break
		}
	}
	// If the for loop failed, then product will be the "zero-value" of the Product struct
	if product == (Product{}) {
		log.Fatal("Can't find product with that ID")
	}
	// Template rendering
	tmpl, err := template.ParseFiles("./templates/product.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, product)
	if err != nil {
		log.Fatal(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/login.html")
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == "POST" {
		rUsername := r.FormValue("username")
		cookie := http.Cookie{Name: "cafego_username", Value: rUsername}
		http.SetCookie(w, &cookie)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/product/", productHandler)
	http.HandleFunc("/login/", loginHandler)
	http.ListenAndServe(":5000", nil)
}
