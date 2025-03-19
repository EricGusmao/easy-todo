package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EricGusmao/easy-todo/user"
)

type MockAuthService struct {
	User  *user.User
	Error error
}

func (m *MockAuthService) Login(ctx context.Context, r *LoginUserRequest) (string, error) {
	panic("unimplemented")
}

func (m *MockAuthService) Signup(ctx context.Context, r *CreateUserRequest) (string, error) {
	panic("unimplemented")
}

func (m *MockAuthService) UserFromToken(ctx context.Context, token string) (*user.User, error) {
	if token == "valid-token" {
		return m.User, nil
	}
	return nil, errors.New("Invalid token")
}

func TestAuthMiddleware(t *testing.T) {
	mockUser := &user.User{ID: 1, Email: "test@example.com"}
	mockAuthService := &MockAuthService{User: mockUser, Error: nil}
	authMiddleware := NewMiddleware(mockAuthService)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid token",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.token)
			rr := httptest.NewRecorder()

			handler := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("request returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
