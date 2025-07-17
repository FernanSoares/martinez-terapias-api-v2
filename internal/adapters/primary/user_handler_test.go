package primary

import (
	"bytes"
	"encoding/json"
	"martinezterapias/api/internal/core/domain"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUserService implementa um mock do UserService para testes
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAll() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) GetByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// UserHandlerSuite define a suite de testes para UserHandler
type UserHandlerSuite struct {
	suite.Suite
	mockService *MockUserService
	handler     *UserHandler
	router      *gin.Engine
}

// SetupTest prepara o ambiente para cada teste
func (s *UserHandlerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockService = new(MockUserService)
	s.handler = NewUserHandler(s.mockService)
	s.router = gin.New()
}

// TestGetAll testa o endpoint de listar todos os usuários
func (s *UserHandlerSuite) TestGetAll() {
	users := []domain.User{
		{
			ID:    1,
			Name:  "user1",
			Email: "user1@example.com",
		},
		{
			ID:    2,
			Name:  "user2",
			Email: "user2@example.com",
		},
	}

	s.mockService.On("GetAll").Return(users, nil)

	s.router.GET("/users", s.handler.GetAll)
	req, _ := http.NewRequest("GET", "/users", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)

	var respUsers []domain.User
	err := json.Unmarshal(resp.Body.Bytes(), &respUsers)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), len(users), len(respUsers))
}

// TestGetByID testa o endpoint de buscar usuário por ID
func (s *UserHandlerSuite) TestGetByID() {
	id := uint(1)
	user := &domain.User{
		ID:    id,
		Name:  "user_test",
		Email: "test@example.com",
	}

	s.mockService.On("GetByID", id).Return(user, nil)

	s.router.GET("/users/:id", s.handler.GetByID)
	req, _ := http.NewRequest("GET", "/users/1", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)

	var respUser domain.User
	err := json.Unmarshal(resp.Body.Bytes(), &respUser)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), user.ID, respUser.ID)
	assert.Equal(s.T(), user.Name, respUser.Name)
}

// TestCreate testa o endpoint de criação de usuário
func (s *UserHandlerSuite) TestCreate() {
	user := domain.User{
		Name:  "new_user",
		Email: "new@example.com",
		Password: "password123",
	}

	s.mockService.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	s.router.POST("/users", s.handler.Create)
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusCreated, resp.Code)
}

// TestUpdate testa o endpoint de atualização de usuário
func (s *UserHandlerSuite) TestUpdate() {
	id := uint(1)
	user := domain.User{
		ID:    id,
		Name:  "updated_user",
		Email: "updated@example.com",
	}

	s.mockService.On("GetByID", id).Return(&user, nil)
	s.mockService.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	s.router.PUT("/users/:id", s.handler.Update)
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/users/"+strconv.FormatUint(uint64(id), 10), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

// TestDelete testa o endpoint de exclusão de usuário
func (s *UserHandlerSuite) TestDelete() {
	id := uint(123)
	s.mockService.On("Delete", id).Return(nil)

	s.router.DELETE("/users/:id", s.handler.Delete)
	req, _ := http.NewRequest("DELETE", "/users/"+strconv.FormatUint(uint64(id), 10), nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNoContent, resp.Code)
}

// TestUserHandlerSuite executa a suite de testes
func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerSuite))
}
