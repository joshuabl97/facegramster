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

func (ui *UI) execTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tpl, err := template.ParseFiles("../../ui/templates/home.html.tmpl")
	if err != nil {
		ui.lg.Error().Err(err)
		http.Error(w, "could not parse html template", http.StatusInternalServerError)
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		ui.lg.Error().Err(err)
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}

func (ui *UI) Homepage(w http.ResponseWriter, r *http.Request) {
	ui.execTemplate(w, "../../ui/templates/home.html.tmpl")
}

func (ui *UI) ContactPage(w http.ResponseWriter, r *http.Request) {
	ui.execTemplate(w, "../../ui/templates/contact.html.tmpl")
}
