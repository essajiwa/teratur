package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	userSvc "github.com/essajiwa/teratur/internal/service/user"

	"github.com/gorilla/mux"

	gomock "github.com/golang/mock/gomock"
)

func TestHandler_GetUserHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockData := NewMockIUserSvc(controller)
	svc := New(mockData)
	r := mux.NewRouter()
	r.HandleFunc("/user/{id:[0-9]+}", svc.GetUserHandler).Methods("GET")

	rr, err := http.NewRequest("GET", "/user/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	ww := httptest.NewRecorder()

	mockData.EXPECT().GetUserByID(context.Background(), int64(1)).
		Return(userSvc.User{
			ID:     1,
			Name:   "Dadang",
			Status: 1,
		}, nil)

	r.ServeHTTP(ww, rr)

	resp := ww.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code %d", resp.StatusCode)
	}

}
