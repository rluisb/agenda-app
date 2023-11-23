package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rluisb/agenda-app/src/db"
	"github.com/rluisb/agenda-app/src/helper"
	"github.com/rluisb/agenda-app/src/types"
)



type ContactHandler struct {
	ContactStore db.ContactStore
}

func NewContactHandler(contactStore db.ContactStore) *ContactHandler {
	return &ContactHandler{
		ContactStore: contactStore,
	}
}

func (handler *ContactHandler) HandlePostContact(w http.ResponseWriter, r *http.Request) {
	var params types.CreateContactParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}
	if errors := params.Validate(); len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errors)
		return
	}

	contact := types.NewContactFromParams(params)
	insertedContact, err := handler.ContactStore.CreateContact(r.Context(), contact)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(helper.NewCustomError(err))
			return	
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insertedContact)
}

func (handler *ContactHandler) HandleListContacts(w http.ResponseWriter, r *http.Request) {
	queryParams := types.NewGetContactsListQueryParams(r.URL.Query())
	users, err := handler.ContactStore.GetContacts(r.Context(), queryParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (handler *ContactHandler) HandleGetContact(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")		

	user, err := handler.ContactStore.GetContactByID(r.Context()	, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(helper.NewCustomError(err))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (handler *ContactHandler) HandleDeleteContact(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")		

	err := handler.ContactStore.DeleteContact(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *ContactHandler) HandleUpdateContact(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")		
	
	var params *types.UpdateContactParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}

	err := handler.ContactStore.UpdateContact(r.Context(), id, params)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(helper.NewCustomError(err))
			return	
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}
	updatedContact, err := handler.ContactStore.GetContactByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedContact)	
}

func (handler *ContactHandler) HandlePatchContact(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")		
	
	var params *types.UpdateContactParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}

	err := handler.ContactStore.PatchContact(r.Context(), id, params)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(helper.NewCustomError(err))
			return	
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}
	patchedContact, err := handler.ContactStore.GetContactByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(helper.NewCustomError(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patchedContact)	
}

