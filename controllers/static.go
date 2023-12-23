package controllers

import (
	"net/http"

	"github.com/joshuabl97/facegramster/ui"
)

func StaticHandler(tpl ui.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Exec(w, nil)
	}
}
