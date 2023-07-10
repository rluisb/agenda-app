package controllers

import (
	"encoding/json"
	"net/http"
)

type ContactController struct {}

func NewContactController() *ContactController {
	return &ContactController{}
}

func (c ContactController) handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"message":"name is required"})
}