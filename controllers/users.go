package controllers

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type Users struct {
	Log       *zerolog.Logger
	Templates struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Exec(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "email: %v \npassword: %v", r.FormValue("email"), r.FormValue("password"))
}
