package shop

import (
	"context"

	"github.com/essajiwa/teratur/internal/service/shop"
	"github.com/essajiwa/teratur/pkg/httpclient"
)

// Data object for data layer needs
type Data struct {
	client  *httpclient.Client
	baseURL string
}

// New return an ads resource
func New(client *httpclient.Client, baseURL string) Data {
	d := Data{
		client:  client,
		baseURL: baseURL,
	}

	return d
}

// GetShopUser will return Shop User information
func (d Data) GetShopUser(ctx context.Context, userID int64) (shop.Shop, error) {
	var shop shop.DataResp
	var endpoint = d.baseURL + "b/5dcbfe27f05d9041253adcb0/2"
	_, err := d.client.GetJSON(ctx, endpoint, nil, &shop)
	if err != nil {
		return shop.Data, err
	}
	return shop.Data, nil
}
