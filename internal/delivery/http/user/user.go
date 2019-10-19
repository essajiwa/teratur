package user

//go:generate mockgen -source=user.go -package=user -destination=user_mock_test.go
import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	userSvc "github.com/essajiwa/teratur/internal/service/user"
	"github.com/essajiwa/teratur/pkg/response"
)

// IUserSvc should be implemented by service
type IUserSvc interface {
	GetUserByID(ctx context.Context, userID int64) (userSvc.User, error)
}

type (
	// Handler will hold dependency on this User handler
	Handler struct {
		usrSvc IUserSvc
	}

	userResp struct {
		Name   string `json:"name"`
		ID     int64  `json:"id"`
		Status int    `json:"status"`
	}
)

// New for user handler initialization
func New(u IUserSvc) *Handler {
	return &Handler{
		usrSvc: u,
	}
}

// GetUserHandler return user data
func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var (
		resp   *response.Response
		userID int64
		err    error
	)
	resp = &response.Response{}
	defer resp.RenderJSON(w, r)

	vars := mux.Vars(r)
	userID, err = strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return
	}

	u, err := h.usrSvc.GetUserByID(context.Background(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "service") {

		}
		resp.Error = response.Error{
			Code:   101,
			Msg:    "Data Not Found",
			Status: true,
		}

		log.Printf("[ERROR] %s %s\n", r.Method, r.URL)
		log.Println(err)
		return
	}

	resp.Data = userResp{
		Name:   u.Name,
		ID:     u.ID,
		Status: u.Status,
	}
	log.Printf("[INFO] %s %s\n", r.Method, r.URL)
}
