package ui

import (
	"net/http"
	"text/template"

	"github.com/rs/zerolog"
)

type UI struct {
	lg *zerolog.Logger
}

func New(logger *zerolog.Logger) UI {
	return UI{
		lg: logger,
	}
}

func (ui *UI) Homepage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tpl, err := template.ParseFiles("../../ui/templates/home.html.tmpl")
	if err != nil {
		ui.lg.Error().Err(err)
		http.Error(w, "could not parse html template", 500)
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		ui.lg.Error().Err(err)
		http.Error(w, "failed to render template", 500)
	}
}

func (ui *UI) ContactPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tpl, err := template.ParseFiles("../../ui/templates/contact.html.tmpl")
	if err != nil {
		ui.lg.Error().Err(err)
		http.Error(w, "could not parse html template", 500)
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		ui.lg.Error().Err(err)
		http.Error(w, "failed to render template", 500)
	}
}
