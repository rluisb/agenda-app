package controllers

import "net/http"

type ContactController struct {}

func NewContactController() *ContactController {
	return &ContactController{}
}

func (c ContactController) handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}