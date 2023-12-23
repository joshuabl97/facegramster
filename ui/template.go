package ui

import (
	"fmt"
	"html/template"
	"net/http"
)

type Template struct {
	htnlTmpl *template.Template
	ui       *UI
}

func (ui *UI) Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		htnlTmpl: tpl,
		ui:       ui,
	}, nil
}

func (t *Template) Exec(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htnlTmpl.Execute(w, data)
	if err != nil {
		t.ui.lg.Error().Err(err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
