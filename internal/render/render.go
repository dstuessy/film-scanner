package render

import (
	"fmt"
	"html/template"
	"net/http"
)

const layoutDir = "web/layout"
const pageDir = "web/pages"
const componentDir = "web/components"

func RenderPage(w http.ResponseWriter, path string, data interface{}) error {
	layoutTmpl, err := template.ParseFiles(fmt.Sprintf("%s/master.html", layoutDir))
	if err != nil {
		return err
	}

	pageTmpl, err := template.Must(layoutTmpl.Clone()).ParseFiles(fmt.Sprintf("%s%s", pageDir, path))
	if err != nil {
		return err
	}

	cmpGlob := fmt.Sprintf("%s/*.html", componentDir)
	if _, err := pageTmpl.ParseGlob(cmpGlob); err != nil {
		return err
	}

	if err := pageTmpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}

func RenderComponent(w http.ResponseWriter, path string, data interface{}) error {
	cmpTmpl, err := template.ParseFiles(fmt.Sprintf("%s%s", componentDir, path))
	if err != nil {
		return err
	}

	if err := cmpTmpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}
