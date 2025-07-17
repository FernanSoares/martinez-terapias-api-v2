package primary

import (
	"bytes"
	"encoding/json"
	"martinezterapias/api/internal/core/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock para o serviço de usuário
type MockUsuarioService struct {
	mock.Mock
}

func (m *MockUsuarioService) GetAll() ([]domain.Usuario, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Usuario), args.Error(1)
}

func (m *MockUsuarioService) GetByID(id uuid.UUID) (*domain.Usuario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioService) GetPerfilCliente(id uuid.UUID) (*domain.Usuario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioService) GetAgendamentosCliente(id uuid.UUID, usuarioSolicitanteID uuid.UUID) ([]domain.Agendamento, error) {
	args := m.Called(id, usuarioSolicitanteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockUsuarioService) Registrar(usuario *domain.Usuario, password string) error {
	args := m.Called(usuario, password)
	return args.Error(0)
}

func (m *MockUsuarioService) Login(email string, password string) (*domain.TokenResponse, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenResponse), args.Error(1)
}

func (m *MockUsuarioService) Update(usuario *domain.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Suite de testes para o AuthHandler
type AuthHandlerSuite struct {
	suite.Suite
	mockService *MockUsuarioService
	handler     *AuthHandler
	router      *gin.Engine
}

// SetupTest configura o ambiente de teste antes de cada caso de teste
func (s *AuthHandlerSuite) SetupTest() {
	// Configurar o Gin para modo de teste
	gin.SetMode(gin.TestMode)

	// Inicializar o mock
	s.mockService = new(MockUsuarioService)
	
	// Inicializar o handler com o mock
	s.handler = NewAuthHandler(s.mockService)
	
	// Configurar o router
	s.router = gin.New()
	s.router.POST("/api/registrar", s.handler.Registrar)
	s.router.POST("/api/login", s.handler.Login)
	s.router.POST("/api/registrar-admin", s.handler.RegistrarAdmin)
}

// TestLogin testa o endpoint de login
func (s *AuthHandlerSuite) TestLogin() {
	// Preparar dados de teste
	email := "admin@example.com"
	password := "admin123"
	userID := uuid.New()
	
	// Dados para o request
	loginRequest := LoginRequest{
		Email: email,
		Senha: password,
	}
	
	// Converter para JSON
	jsonValue, _ := json.Marshal(loginRequest)
	
	// Caso 1: Login bem-sucedido
	// Preparar resposta do serviço
	tokenResponse := &domain.TokenResponse{
		Token: "jwt-token-exemplo",
		Usuario: struct {
			ID     uuid.UUID        `json:"id"`
			Nome   string          `json:"nome"`
			Perfil domain.TipoPerfil `json:"perfil"`
		}{
			ID:     userID,
			Nome:   "Administrador",
			Perfil: domain.PerfilAdmin,
		},
	}
	
	// Configurar mock para o caso de sucesso
	s.mockService.On("Login", email, password).Return(tokenResponse, nil).Once()
	
	// Criar request
	req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Criar response recorder
	w := httptest.NewRecorder()
	
	// Executar request
	s.router.ServeHTTP(w, req)
	
	// Verificar resposta
	assert.Equal(s.T(), http.StatusOK, w.Code)
	
	// Verificar corpo da resposta
	var responseBody domain.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), tokenResponse.Token, responseBody.Token)
	assert.Equal(s.T(), tokenResponse.Usuario.ID, responseBody.Usuario.ID)
	assert.Equal(s.T(), tokenResponse.Usuario.Nome, responseBody.Usuario.Nome)
	assert.Equal(s.T(), tokenResponse.Usuario.Perfil, responseBody.Usuario.Perfil)
	
	// Verificar chamadas mock
	s.mockService.AssertExpectations(s.T())
	
	// Resetar para próximo teste
	s.SetupTest()
	
	// Caso 2: Credenciais inválidas
	// Configurar mock para retornar erro
	s.mockService.On("Login", email, password).Return(nil, 
		assert.AnError).Once()
	
	// Criar request
	req, _ = http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Criar response recorder
	w = httptest.NewRecorder()
	
	// Executar request
	s.router.ServeHTTP(w, req)
	
	// Verificar resposta
	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	
	// Verificar erro na resposta
	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.Contains(s.T(), errorResponse["error"], "Credenciais inválidas")
	
	// Verificar chamadas mock
	s.mockService.AssertExpectations(s.T())
	
	// Resetar para próximo teste
	s.SetupTest()
	
	// Caso 3: Erro no formato do JSON
	// Criar request com JSON inválido
	invalidJSON := []byte(`{"email": "admin@example.com"`) // JSON mal formado
	req, _ = http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Criar response recorder
	w = httptest.NewRecorder()
	
	// Executar request
	s.router.ServeHTTP(w, req)
	
	// Verificar resposta
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// TestRegistrar testa o endpoint de registro
func (s *AuthHandlerSuite) TestRegistrar() {
	// Preparar dados de teste
	registroRequest := RegistroRequest{
		NomeCompleto: "Novo Cliente",
		Email:        "cliente@exemplo.com",
		Senha:        "senha123",
		Telefone:     "(11) 99999-9999",
	}
	
	// Converter para JSON
	jsonValue, _ := json.Marshal(registroRequest)
	
	// Caso 1: Registro bem-sucedido
	// Configurar mock para o caso de sucesso
	s.mockService.On("Registrar", mock.AnythingOfType("*domain.Usuario"), registroRequest.Senha).Return(nil).Once()
	
	// Criar request
	req, _ := http.NewRequest(http.MethodPost, "/api/registrar", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Criar response recorder
	w := httptest.NewRecorder()
	
	// Executar request
	s.router.ServeHTTP(w, req)
	
	// Verificar resposta
	assert.Equal(s.T(), http.StatusCreated, w.Code)
	
	// Verificar chamadas mock
	s.mockService.AssertExpectations(s.T())
	
	// Resetar para próximo teste
	s.SetupTest()
	
	// Caso 2: Erro no registro (email já existe)
	// Configurar mock para retornar erro
	s.mockService.On("Registrar", mock.AnythingOfType("*domain.Usuario"), registroRequest.Senha).Return(
		assert.AnError).Once()
	
	// Criar request
	req, _ = http.NewRequest(http.MethodPost, "/api/registrar", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Criar response recorder
	w = httptest.NewRecorder()
	
	// Executar request
	s.router.ServeHTTP(w, req)
	
	// Verificar resposta
	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	
	// Verificar chamadas mock
	s.mockService.AssertExpectations(s.T())
}

// TestRegistrarAdmin testa o endpoint de registro de administrador
func (s *AuthHandlerSuite) TestRegistrarAdmin() {
	// Preparar dados de teste
	registroRequest := RegistroRequest{
		NomeCompleto: "Novo Admin",
		Email:        "admin@exemplo.com",
		Senha:        "senha123",
		Telefone:     "(11) 88888-8888",
	}
	
	// Converter para JSON
	jsonValue, _ := json.Marshal(registroRequest)
	
	// Caso 1: Registro de admin bem-sucedido quando não há admin
	// Configurar mock para retornar lista vazia de usuários (nenhum admin)
	s.mockService.On("GetAll").Return([]domain.Usuario{}, nil).Once()
	
	// Configurar mock para o caso de sucesso no registro
	s.mockService.On("Registrar", mock.AnythingOfType("*domain.Usuario"), registroRequest.Senha).Return(nil).Once()
	
	// Criar request
	req, _ := http.NewRequest(http.MethodPost, "/api/registrar-admin", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Criar response recorder
	w := httptest.NewRecorder()
	
	// Executar request
	s.router.ServeHTTP(w, req)
	
	// Verificar resposta
	assert.Equal(s.T(), http.StatusCreated, w.Code)
	
	// Verificar chamadas mock
	s.mockService.AssertExpectations(s.T())
	
	// Resetar para próximo teste
	s.SetupTest()
	
	// Caso 2: Tentativa de registro quando já existe admin
	// Lista de usuários com um admin
	usuarios := []domain.Usuario{
		{
			ID:           uuid.New(),
			NomeCompleto: "Admin Existente",
			Email:        "admin.existente@exemplo.com",
			Perfil:       domain.PerfilAdmin,
			Ativo:        true,
			DataCadastro: time.Now().Add(-24 * time.Hour),
		},
	}
	
	// Configurar mock para retornar lista com admin
	s.mockService.On("GetAll").Return(usuarios, nil).Once()
	
	// Criar request
	req, _ = http.NewRequest(http.MethodPost, "/api/registrar-admin", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	
	// Criar response recorder
	w = httptest.NewRecorder()
	
	// Executar request
	s.router.ServeHTTP(w, req)
	
	// Verificar resposta
	assert.Equal(s.T(), http.StatusForbidden, w.Code)
	
	// Verificar que o método Registrar não foi chamado
	s.mockService.AssertNotCalled(s.T(), "Registrar")
}

// TestAuthHandlerSuite executa a suite de testes
func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerSuite))
}
