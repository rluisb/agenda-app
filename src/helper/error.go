package helper

type CustomError struct {
	Error string `json:"error"`
}

func NewCustomError(err error) CustomError {
	return CustomError{
		Error: err.Error(),
	}
}