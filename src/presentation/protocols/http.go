package protocols

type HttpResponse struct {
	StatusCode int
	Body       interface{}
}

func NewHttpResponse(statusCode int, body interface{}) *HttpResponse {
	return &HttpResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}