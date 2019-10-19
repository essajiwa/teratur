package user

//go:generate mockgen -source=user.go -package=user -destination=user_mock_test.go
import (
	"context"

	"github.com/essajiwa/teratur/pkg/errors"
)

// DataLayer holds all method to access data,
// implemented at Data Access layer
type DataLayer interface {
	GetUserByID(ctx context.Context, userID int64) (User, error)
}

// ShopService will do wiring to shop related data, put only method used in this service
type ShopService interface {
	GetShopName(ctx context.Context, shopID int64) (string, error)
}

// Service is an object for interacting with User Service
type Service struct {
	data DataLayer
	shop ShopService
}

// New will return new User Service object
// accept parameter interface
func New(userData DataLayer) Service {
	return Service{
		data: userData,
	}
}

// GetUserByID will return User data filtered by it's ID
func (s Service) GetUserByID(ctx context.Context, userID int64) (User, error) {
	if userID == 3 {
		return User{}, errors.New("pengen aja ngasih error dari service")
	}
	return s.data.GetUserByID(ctx, userID)
}
