package primary

import (
	"errors"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUsuarioServiceForPerfil implementa um mock do UsuarioService para testes
type MockUsuarioServiceForPerfil struct {
	mock.Mock
}

func (m *MockUsuarioServiceForPerfil) Registrar(usuario *domain.Usuario, senha string) error {
	args := m.Called(usuario, senha)
	return args.Error(0)
}

func (m *MockUsuarioServiceForPerfil) Login(email, senha string) (*domain.TokenResponse, error) {
	args := m.Called(email, senha)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenResponse), args.Error(1)
}

func (m *MockUsuarioServiceForPerfil) ValidateToken(tokenString string) (map[string]interface{}, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockUsuarioServiceForPerfil) GetByID(id uuid.UUID) (*domain.Usuario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioServiceForPerfil) GetByEmail(email string) (*domain.Usuario, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioServiceForPerfil) GetAll() ([]domain.Usuario, error) {
	args := m.Called()
	return args.Get(0).([]domain.Usuario), args.Error(1)
}

func (m *MockUsuarioServiceForPerfil) GetAgendamentos(usuarioID uuid.UUID) ([]domain.Agendamento, error) {
	args := m.Called(usuarioID)
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockUsuarioServiceForPerfil) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUsuarioServiceForPerfil) GetAgendamentosCliente(usuarioID uuid.UUID, clienteID uuid.UUID) ([]domain.Agendamento, error) {
	args := m.Called(usuarioID, clienteID)
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockUsuarioServiceForPerfil) GetPerfilCliente(id uuid.UUID) (*domain.Usuario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioServiceForPerfil) SoftDelete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUsuarioServiceForPerfil) Update(usuario *domain.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

// Estas estruturas e funções já estão definidas no arquivo principal

// PerfilHandlerSuite define a suite de testes para PerfilHandler
type PerfilHandlerSuite struct {
	suite.Suite
	mockService *MockUsuarioServiceForPerfil
	handler     *PerfilHandler
	router      *gin.Engine
}

// SetupTest prepara o ambiente para cada teste
func (s *PerfilHandlerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockService = new(MockUsuarioServiceForPerfil)
	// Usando ports.UsuarioService para o tipo do serviço
	var usuarioService ports.UsuarioService = s.mockService
	s.handler = NewPerfilHandler(usuarioService)
	s.router = gin.New()

	// Configurar middleware para simular autenticação
	s.router.Use(func(c *gin.Context) {
		c.Set("userID", "35289a20-c7c3-4b98-8661-c43f48c45f8e") // ID de usuário fixo para testes
		c.Next()
	})
}

// TestGetPerfil testa o endpoint de obter o próprio perfil
func (s *PerfilHandlerSuite) TestGetPerfil() {
	userID, _ := uuid.Parse("35289a20-c7c3-4b98-8661-c43f48c45f8e")
	usuario := &domain.Usuario{
		ID:          userID,
		NomeCompleto: "Usuário Teste",
		Email:       "teste@exemplo.com",
		Telefone:    "11999998888",
		Perfil:      domain.PerfilCliente,
		Ativo:       true,
	}

	s.mockService.On("GetByID", userID).Return(usuario, nil)

	s.router.GET("/me/perfil", s.handler.GetPerfil)
	req, _ := http.NewRequest("GET", "/me/perfil", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

// TestGetPerfil_Error testa erro ao obter o próprio perfil
func (s *PerfilHandlerSuite) TestGetPerfil_Error() {
	userID, _ := uuid.Parse("35289a20-c7c3-4b98-8661-c43f48c45f8e")

	s.mockService.On("GetByID", userID).Return(nil, errors.New("usuário não encontrado"))

	s.router.GET("/me/perfil", s.handler.GetPerfil)
	req, _ := http.NewRequest("GET", "/me/perfil", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
}

// TestGetAgendamentos testa o endpoint de obter agendamentos do próprio perfil
func (s *PerfilHandlerSuite) TestGetAgendamentos() {
	userID, _ := uuid.Parse("35289a20-c7c3-4b98-8661-c43f48c45f8e")
	agendamentos := []domain.Agendamento{
		{
			ID:        uuid.New(),
			ClienteID: userID,
			Status:    domain.StatusAgendado,
		},
		{
			ID:        uuid.New(),
			ClienteID: userID,
			Status:    domain.StatusConfirmado,
		},
	}

	s.mockService.On("GetAgendamentos", userID).Return(agendamentos, nil)

	s.router.GET("/me/agendamentos", s.handler.GetAgendamentos)
	req, _ := http.NewRequest("GET", "/me/agendamentos", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

// TestGetAgendamentos_Error testa erro ao obter agendamentos do próprio perfil
func (s *PerfilHandlerSuite) TestGetAgendamentos_Error() {
	userID, _ := uuid.Parse("35289a20-c7c3-4b98-8661-c43f48c45f8e")

	s.mockService.On("GetAgendamentos", userID).Return([]domain.Agendamento{}, errors.New("agendamentos não encontrados"))

	s.router.GET("/me/agendamentos", s.handler.GetAgendamentos)
	req, _ := http.NewRequest("GET", "/me/agendamentos", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
}

// TestPerfilHandlerSuite executa a suite de testes
func TestPerfilHandlerSuite(t *testing.T) {
	suite.Run(t, new(PerfilHandlerSuite))
}
