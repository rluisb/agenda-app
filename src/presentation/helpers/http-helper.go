package helpers

import (
	"net/http"
	"strings"

	"github.com/rluisb/agenda-app/src/presentation/protocols"
)

func BadRequest(param string) *protocols.HttpResponse {
	return protocols.NewHttpResponse(http.StatusBadRequest, map[string]string{"message": 	strings.ToLower(param) + " is required"})
}