package ui

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

type Template struct {
	htmlTmpl *template.Template
	ui       *UI
}

func (ui *UI) Must(tpl Template, err error) Template {
	if err != nil {
		ui.lg.Fatal().Err(err).Msg("Failed to parse template")
	}
	return tpl
}

func (ui *UI) Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		htmlTmpl: tpl,
		ui:       ui,
	}, nil
}

func (ui *UI) ParseFS(fs fs.FS, pattern string) (Template, error) {
	tpl, err := template.ParseFS(fs, pattern)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTmpl: tpl,
		ui:       ui,
	}, nil
}

func (t *Template) Exec(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTmpl.Execute(w, data)
	if err != nil {
		t.ui.lg.Error().Err(err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
