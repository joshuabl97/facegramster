package controllers

import "net/http"

type Template interface {
	Exec(w http.ResponseWriter, data interface{})
}
