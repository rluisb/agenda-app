package custom_errors

import "strings"

type InvalidParamError struct {
	Param string
}

func NewInvalidParamError(param string) InvalidParamError {
	return InvalidParamError{Param: param}
}

func (e InvalidParamError) Build() map[string]string {
	return map[string]string{"message": "invalid " + strings.ToLower(e.Param)}
}