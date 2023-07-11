package helpers

import (
	"net/http"

	"github.com/rluisb/agenda-app/src/presentation/protocols"
)

func BadRequest(message map[string]string) *protocols.HttpResponse {
	return protocols.NewHttpResponse(http.StatusBadRequest, message)
}

func InternalServerError(message map[string]string) *protocols.HttpResponse {
	return protocols.NewHttpResponse(http.StatusInternalServerError, message)
}