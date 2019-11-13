package user

//go:generate mockgen -source=user.go -package=user -destination=user_mock_test.go
import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	shopSvc "github.com/essajiwa/teratur/internal/service/shop"
	userSvc "github.com/essajiwa/teratur/internal/service/user"
	"github.com/essajiwa/teratur/pkg/response"
)

// IUserSvc should be implemented by service
type IUserSvc interface {
	GetUserByID(ctx context.Context, userID int64) (userSvc.User, error)
}

// IShopSvc interfacing with shop service
type IShopSvc interface {
	GetShopByUserID(ctx context.Context, userID int64) (shopSvc.Shop, error)
}

type (
	// Handler will hold dependency on this User handler
	Handler struct {
		usrSvc  IUserSvc
		shopSvc IShopSvc
	}

	userResp struct {
		Name     string `json:"name"`
		ID       int64  `json:"id"`
		Status   int    `json:"status"`
		ShopName string `json:"shop_name"`
	}
)

// New for user handler initialization
func New(u IUserSvc, s IShopSvc) *Handler {
	return &Handler{
		usrSvc:  u,
		shopSvc: s,
	}
}

// GetUserHandler return user data
func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var (
		resp   *response.Response
		userID int64
		err    error
		errRes response.Error
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

		errRes = response.Error{
			Code:   101,
			Msg:    "Data Not Found",
			Status: true,
		}

		if strings.Contains(err.Error(), "service") {

			errRes = response.Error{
				Code:   201,
				Msg:    "Failed to process request due to server error",
				Status: true,
			}
		}

		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
		resp.Error = errRes
		return
	}

	// call shop service
	shop, err := h.shopSvc.GetShopByUserID(context.Background(), userID)
	if err != nil {
		errRes = response.Error{
			Code:   301,
			Msg:    "Have no Shop",
			Status: true,
		}

		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
		resp.Error = errRes
		return

	}

	resp.Data = userResp{
		Name:     u.Name,
		ID:       u.ID,
		Status:   u.Status,
		ShopName: shop.Name,
	}

	log.Printf("[INFO] %s %s\n", r.Method, r.URL)
}
