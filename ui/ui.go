package ui

import (
	"net/http"
	"text/template"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles("../../ui/templates/home.html.tmpl")
	if err != nil {
		http.Error(w, "could not parse html template", 500)
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "failed to render template", 500)
	}
}
