package custom_errors

import "strings"

type MissingParamError struct {
	Param string
}

func NewMissingParamError(param string) MissingParamError {
	return MissingParamError{Param: param}
}

func (e MissingParamError) Build() map[string]string {
	return map[string]string{"message": strings.ToLower(e.Param) + " is required"}
}