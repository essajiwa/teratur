package user

import (
	"context"
	"log"

	svc "github.com/essajiwa/teratur/internal/service/user"
	"github.com/essajiwa/teratur/pkg/safesql"
	"github.com/pkg/errors"
)

type (
	// Data object for data layer needs
	Data struct {
		db   safesql.SlaveDB
		stmt map[string]safesql.SlaveStatement
	}

	statement struct {
		key   string
		query string
	}
)

const (
	getUserByID  = "GetByUserID"
	qGetUserByID = `
	SELECT
		"id",
		"name",
		"status"
	FROM 
		"public"."user" 
	WHERE "id" = $1;`

	getAllUser  = "GetAllUser"
	qGetAllUser = `
	SELECT
		"id",
		"name",
		"status"
	FROM 
		"public"."user";`
)

var (
	readStmt = []statement{
		{getUserByID, qGetUserByID},
		{getAllUser, qGetAllUser},
	}
)

// New will return Data object
func New(slaveDB safesql.SlaveDB) Data {
	d := Data{
		db: slaveDB,
	}

	d.initStmt()
	return d
}

func (d *Data) initStmt() {
	var (
		err   error
		stmts = make(map[string]safesql.SlaveStatement)
	)

	for _, v := range readStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("failed init statement key %v, err : %v", v.key, err)
		}
	}

	d.stmt = stmts
}

// GetUserByID will query User data from DB
func (d Data) GetUserByID(ctx context.Context, userID int64) (svc.User, error) {
	var (
		user svc.User
		err  error
	)

	row := d.stmt[getUserByID].QueryRowxContext(ctx, userID)
	err = row.StructScan(&user)
	if err != nil {
		return user, errors.Wrap(err, "[data][GetByUserID]")
	}

	return user, err
}
