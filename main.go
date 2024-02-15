package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const webDir = "web"
const pageDir = "web/pages"

func parsePage(path string, data interface{}) (*template.Template, error) {
	layoutTmpl, err := template.ParseFiles(fmt.Sprintf("%s/layout/master.html", webDir))
	if err != nil {
		return nil, err
	}

	pageTmpl, err := template.Must(layoutTmpl.Clone()).ParseFiles(path)
	if err != nil {
		return nil, err
	}

	return pageTmpl, nil
}

func createPageHandler(filePath string, _ fs.FileInfo, err error) error {
	if err != nil {
		return err
	}

	isHtmlPage := strings.HasSuffix(filePath, ".html")
	routePath := strings.TrimSuffix(filePath, ".html")
	routePath = strings.TrimPrefix(routePath, pageDir)
	routePath = strings.TrimSuffix(routePath, "index")

	fmt.Println("routePath", routePath, filePath)

	if isHtmlPage {
		http.HandleFunc(routePath, func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := parsePage(filePath, nil)
			if err != nil {
				log.Fatal(err)
			}

			if err := tmpl.Execute(w, nil); err != nil {
				log.Fatal(err)
			}
		})
	}

	return nil
}

func main() {
	if err := filepath.Walk(pageDir, createPageHandler); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
