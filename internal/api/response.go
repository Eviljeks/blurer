package api

import "net/http"

// swagger:model
type Error struct {
	// HTTP status code
	//
	// required: true
	// example: 403
	Status int `json:"status"`

	// User message
	//
	// required: true
	// example: Forbidden
	Message string `json:"user_message,omitempty"`

	// Internal error code
	//
	// example: 1
	Code int `json:"code"`
}

// swagger:model
type Response struct {
	// Result
	//
	// required: true
	Result interface{} `json:"result"`
	// Errors
	//
	// required: true
	Errors []Error `json:"errors"`
}

func ErrorBadRequest(err error) Response {
	return httpError(err, http.StatusBadRequest)
}

func ErrorNotFound(err error) Response {
	return httpError(err, http.StatusNotFound)
}

func ErrorInternalServer(err error) Response {
	return httpError(err, http.StatusInternalServerError)
}

func httpError(err error, status int) Response {
	return Response{
		Result: nil,
		Errors: []Error{
			{
				Status:  status,
				Code:    0,
				Message: err.Error(),
			},
		},
	}
}

func OK(result interface{}) Response {
	return Response{
		Result: result,
		Errors: nil,
	}
}
