package user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestService_GetUserByID(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockData := NewMockDataLayer(controller)
	svc := New(mockData)

	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		mock    func()
		want    User
		wantErr bool
	}{
		{
			name: "Test Success",
			args: args{
				userID: 35007,
				ctx:    context.Background(),
			},
			mock: func() {
				mockData.EXPECT().
					GetUserByID(context.Background(), int64(35007)).
					Return(User{
						ID:     35007,
						Name:   "Dadang",
						Status: 1,
					}, nil)
			},
			want: User{
				ID:     35007,
				Name:   "Dadang",
				Status: 1,
			},
			wantErr: false,
		},
		{
			name: "Test Failed",
			args: args{
				userID: 35007,
				ctx:    context.Background(),
			},
			mock: func() {
				mockData.EXPECT().
					GetUserByID(context.Background(), int64(35007)).
					Return(User{}, errors.New("nil returned"))
			},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.GetUserByID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
