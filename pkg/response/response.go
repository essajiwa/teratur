package response

import (
	"encoding/json"
	"net/http"
)

// Response defines http response for the client
type Response struct {
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Error      Error       `json:"error"`
	StatusCode int         `json:"-"`
}

// Error defines the error
type Error struct {
	Status bool   `json:"status"` // true if we have error
	Msg    string `json:"msg"`    // error message
	Code   int    `json:"code"`   // error code from application, it is not http status code
}

// SetError set the response to return the given error.
// code is http status code, http.StatusInternalServerError is the default value
func (res *Response) SetError(err error, code ...int) {
	if len(code) > 0 {
		res.StatusCode = code[0]
	} else {
		res.StatusCode = http.StatusInternalServerError
	}

	if err != nil {
		res.Error = Error{
			Msg:    err.Error(),
			Status: true,
		}
	}

}

// RenderJSON writes the http response in JSON format to the client
func (res *Response) RenderJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if res.StatusCode == 0 {
		res.StatusCode = http.StatusOK
	}

	d, err := json.Marshal(res)
	if err != nil {
		res.StatusCode = http.StatusInternalServerError
	}

	w.WriteHeader(res.StatusCode)
	w.Write(d)
}
