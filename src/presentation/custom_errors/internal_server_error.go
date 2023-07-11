package custom_errors

type InternalServerError struct {
}

func NewInternalServerError() InternalServerError {
	return InternalServerError{}
}

func (e InternalServerError) Build() map[string]string {
	return map[string]string{"message": "internal server error"}
}