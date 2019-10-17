package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderJSON(t *testing.T) {
	testCases := []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
		status  int
	}{
		{
			name: "OK",
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := Response{}
				resp.RenderJSON(w, r)
			},
			status: http.StatusOK,
		},
		{
			name: "bad request error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := Response{}
				resp.SetError(fmt.Errorf("bad request"), http.StatusBadRequest)
				resp.RenderJSON(w, r)
			},
			status: http.StatusBadRequest,
		},
		{
			name: "default http code on error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := Response{}
				resp.SetError(fmt.Errorf("internal server error"))
				resp.RenderJSON(w, r)
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.handler))
			defer server.Close()

			resp, err := http.Get(server.URL)
			require.NoError(t, err)

			require.Equal(t, tc.status, resp.StatusCode)
			require.Contains(t, resp.Header.Get("Content-Type"), "application/json")
		})
	}

}

func TestErrorMsg(t *testing.T) {
	const (
		someErrorMsg = "some error msg"
	)
	testCases := []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
		status  int
		err     Error
	}{
		{
			name: "with error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := Response{}
				resp.SetError(fmt.Errorf(someErrorMsg), http.StatusBadRequest)
				resp.RenderJSON(w, r)
			},
			status: http.StatusBadRequest,
			err: Error{
				Msg:    someErrorMsg,
				Status: true,
			},
		},
		{
			name: "no error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := Response{}
				resp.RenderJSON(w, r)
			},
			status: http.StatusOK,
			err:    Error{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tc.handler))
			defer server.Close()

			resp, err := http.Get(server.URL)
			require.NoError(t, err)

			require.Equal(t, tc.status, resp.StatusCode)

			var respBody Response
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			require.NoError(t, err)

			require.Equal(t, tc.err, respBody.Error)
		})
	}

}
