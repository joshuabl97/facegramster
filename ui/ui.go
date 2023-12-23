package ui

import (
	"net/http"

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
	tpl, err := ui.Parse(filepath)
	if err != nil {
		ui.lg.Error().Err(err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}

	tpl.Exec(w, nil)
}

func (ui *UI) Homepage(w http.ResponseWriter, r *http.Request) {
	ui.execTemplate(w, "../../ui/templates/home.html.tmpl")
}

func (ui *UI) ContactPage(w http.ResponseWriter, r *http.Request) {
	ui.execTemplate(w, "../../ui/templates/contact.html.tmpl")
}
