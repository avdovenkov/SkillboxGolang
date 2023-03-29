package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"http_service/pkg/service"
	"http_service/pkg/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ApiMock struct {
	mock.Mock
	r *chi.Mux
}

func (m *ApiMock) GetUser(ctx context.Context, id int64) (*user.User, error) {
	m.Called(id)
	user := &user.User{Id: id}
	return user, nil
}
func (m *ApiMock) CreateUser(ctx context.Context, u *user.User) error {
	fmt.Println("Пользователь успешно создан")
	return nil
}
func (m *ApiMock) UpdateUser(ctx context.Context, u *user.User) error {
	return nil
}
func (m *ApiMock) DeleteUser(ctx context.Context, id int64) error {
	return nil
}
func (m *ApiMock) GetUserFriend(ctx context.Context, id int64) ([]*user.Friend, error) {
	return nil, nil
}

func TestGetUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		var testid int64 = 2
		testuser := user.User{Id: testid}
		apimock := new(ApiMock)
		apimock.r = chi.NewRouter()
		service := &service.Service{
			apimock,
		}
		apimock.On("GetUser", testid).Return(testuser)
		apiservice := ApiUser{ctx, service}
		rr := httptest.NewRecorder()
		apimock.r.Get("/{userID}", apiservice.GetOneUser)
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%d", testid), nil)
		assert.NoError(t, err)
		apimock.r.ServeHTTP(rr, request)
		respBody, err := json.Marshal(testuser)
		assert.NoError(t, err)
		assert.Equal(t, 200, rr.Code)
		assert.Equal(t, string(respBody)+"\n", rr.Body.String())
	})
}
func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		var testid int64 = 1
		testuser := user.User{Id: testid, Name: "Володя", Age: 25}
		apimock := new(ApiMock)
		apimock.r = chi.NewRouter()
		service := &service.Service{
			apimock,
		}
		apimock.On("CreateUser", testuser).Return(nil)
		apiservice := ApiUser{ctx, service}
		rr := httptest.NewRecorder()
		apimock.r.Post("/", apiservice.CreateUser)
		b, _ := json.Marshal(testuser)
		request, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
		request.Header.Set("Content-Type", "application/json")
		assert.NoError(t, err)
		apimock.r.ServeHTTP(rr, request)
		respBody, err := json.Marshal(testuser)
		assert.NoError(t, err)
		assert.Equal(t, 201, rr.Code)
		assert.Equal(t, string(respBody)+"\n", rr.Body.String())
	})
}
