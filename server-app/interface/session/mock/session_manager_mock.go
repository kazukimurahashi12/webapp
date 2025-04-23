package mock

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/stretchr/testify/mock"
)

type MockSessionManager struct {
	mock.Mock
}

func NewMockSessionManager(t *testing.T) *MockSessionManager {
	return &MockSessionManager{}
}

func (m *MockSessionManager) CreateSession(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockSessionManager) GetSession(c *gin.Context) (string, error) {
	args := m.Called(c)
	return args.String(0), args.Error(1)
}

func (m *MockSessionManager) DeleteSession(c *gin.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockSessionManager) UpdateSession(c *gin.Context, newID string) error {
	args := m.Called(c, newID)
	return args.Error(0)
}

var _ session.SessionManager = (*MockSessionManager)(nil)
