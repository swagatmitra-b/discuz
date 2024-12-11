package utils

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, page string, data interface{}) {
	path := fmt.Sprintf("cmd/static/%s.html", page)
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "Error in parsing html file", 500)
		return
	}

	tmpl.Execute(w, data)
}
