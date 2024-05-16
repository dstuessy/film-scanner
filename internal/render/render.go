package render

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

const layoutDir = "web/layout"
const pageDir = "web/pages"
const componentDir = "web/components"

var funcs = template.FuncMap{
	"sub": func(a, b int) int {
		return a - b
	},
}

func RenderPage(w http.ResponseWriter, p string, data interface{}) error {
	layoutName := fmt.Sprintf("%s/master.html", layoutDir)
	layoutBase := path.Base(layoutName)

	layoutTmpl, err := template.New(layoutBase).Funcs(funcs).ParseFiles(layoutName)
	if err != nil {
		return err
	}

	pageTmpl, err := template.Must(layoutTmpl.Clone()).ParseFiles(fmt.Sprintf("%s%s", pageDir, p))
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

func RenderComponent(w http.ResponseWriter, p string, data interface{}) error {
	name := fmt.Sprintf("%s%s", componentDir, p)
	base := path.Base(name)

	cmpTmpl, err := template.New(base).Funcs(funcs).ParseFiles(name)
	if err != nil {
		return err
	}

	if err := cmpTmpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}
