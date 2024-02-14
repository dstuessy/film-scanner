package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const WebDir = "web"

func ParsePage(name string, data interface{}) (*template.Template, error) {
	layoutTmpl, err := template.ParseFiles(fmt.Sprintf("%s/layout/master.html", WebDir))
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/pages/%s.html", WebDir, name)
	pageTmpl, err := template.Must(layoutTmpl.Clone()).ParseFiles(path)
	if err != nil {
		return nil, err
	}
	return pageTmpl, nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl, err := ParsePage("index", nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := indexTmpl.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", IndexHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
