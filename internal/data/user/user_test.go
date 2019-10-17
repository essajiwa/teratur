package user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	svc "github.com/essajiwa/teratur/internal/service/user"
	"github.com/essajiwa/teratur/pkg/safesql"
)

const (
	querySelectMock = "SELECT"
)

func mockData() (Data, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return Data{}, mock, err
	}

	// Mocking prepare
	for range readStmt {
		mock.ExpectPrepare(querySelectMock)
	}

	d := New(safesql.NewSlaveDB(mockDB, "postgres"))

	return d, mock, err
}
func TestData_GetUserByID(t *testing.T) {
	d, mock, err := mockData()
	if err != nil {
		t.Fatal("failed init data mocking")
	}

	type args struct {
		ctx    context.Context
		userID int64
	}

	tests := []struct {
		name    string
		d       Data
		args    args
		mock    func(userID int64)
		want    svc.User
		wantErr bool
	}{
		{
			name: "Test 1 [Success]",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			d: d,
			mock: func(userID int64) {
				rows := sqlmock.NewRows([]string{"id", "name", "status"}).AddRow(1, "1", 1)
				mock.ExpectQuery("SELECT (.+)").WithArgs(userID).WillReturnRows(rows)
			},
			want: svc.User{
				ID:     1,
				Name:   "1",
				Status: 1,
			},
		},
		{
			name: "Test 2 [Query Error]",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			d: d,
			mock: func(userID int64) {
				mock.ExpectQuery("SELECT (.+)").WithArgs(userID).
					WillReturnError(errors.New("Expected Error"))
			},
			want:    svc.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mock(tt.args.userID)
			got, err := tt.d.GetUserByID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Data.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
