package product

import (
	"context"
	"fmt"

	"github.com/essajiwa/teratur/pkg/errors"
	"github.com/essajiwa/teratur/pkg/safesql"
)

type (
	// Data object for data layer needs
	Data struct {
		db   safesql.MasterDB
		stmt map[string]safesql.MasterStatement
	}

	statement struct {
		key   string
		query string
	}
)

const (
	qGetNextProductNumber = "SELECT NEXTVAL('product_num_seq') AS product_num"
	qSaveProduct          = `
	INSERT INTO 
		product (
			name,
			product_number,
			create_by,
			create_time,
			update_time
		)
	VALUES ($1, $2, $3, now(), now()) 
	ON conflict DO NOTHING RETURNING product_number, tokopedia_invoice_number;
`
)

// New will return Data object
func New(masterDB safesql.MasterDB) Data {
	d := Data{
		db: masterDB,
	}
	return d
}

// GetUserByID will query User data from DB
func (d Data) generateProductNumber(ctx context.Context) (string, error) {
	var (
		prodNum    int64
		fmtProdNum string
	)

	err := d.db.QueryRowxContext(ctx, qGetNextProductNumber).Scan(&prodNum)
	if err != nil {
		return fmtProdNum, errors.Wrap(err)
	}

	fmtProdNum = fmt.Sprintf("P%0.8d", prodNum)
	return fmtProdNum, nil

}

// SaveProduct will save the product info with auto mumbering for product number
func (d *Data) SaveProduct(ctx context.Context, name string) error {
	return nil
}
