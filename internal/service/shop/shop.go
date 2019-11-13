package shop

//go:generate mockgen -source=shop.go -package=shop -destination=shop_mock_test.go
import (
	"context"
)

// DataLayer holds all method to access data,
// implemented at Data Access layer
type DataLayer interface {
	GetShopUser(ctx context.Context, userID int64) (Shop, error)
}

// Service is an object for interacting with User Service
type Service struct {
	data DataLayer
}

// New will return new User Service object
// accept parameter interface
func New(d DataLayer) Service {
	return Service{
		data: d,
	}
}

// GetShopByUserID will return shop data from API
func (s Service) GetShopByUserID(ctx context.Context, userID int64) (Shop, error) {
	shop, err := s.data.GetShopUser(ctx, userID)
	return shop, err
}
