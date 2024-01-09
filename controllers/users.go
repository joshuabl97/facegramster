package controllers

import (
	"fmt"
	"net/http"

	"github.com/joshuabl97/facegramster/models"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Exec(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email, password := r.FormValue("email"), r.FormValue("password")

	user, err := u.UserService.Create(&models.NewUser{
		Email:    email,
		Password: password,
	})
	if err != nil {
		u.UserService.Lg.Error().Err(err).Msg("Error creating user from UserService")
		http.Error(w, "Error creating user", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "User created &+v", user)
}
