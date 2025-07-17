package primary

import (
	"bytes"
	"encoding/json"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/middleware"
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

// Mock para o serviço de cliente
type MockClienteService struct {
	mock.Mock
}

func (m *MockClienteService) GetAll() ([]domain.Cliente, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *MockClienteService) GetByID(id uuid.UUID) (*domain.Cliente, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cliente), args.Error(1)
}

func (m *MockClienteService) GetByNome(nome string) ([]domain.Cliente, error) {
	args := m.Called(nome)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *MockClienteService) GetByAtivo(ativo bool) ([]domain.Cliente, error) {
	args := m.Called(ativo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *MockClienteService) Create(cliente *domain.Cliente) error {
	args := m.Called(cliente)
	return args.Error(0)
}

func (m *MockClienteService) Update(cliente *domain.Cliente) error {
	args := m.Called(cliente)
	return args.Error(0)
}

func (m *MockClienteService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Suite de testes para o ClienteHandler
type ClienteHandlerSuite struct {
	suite.Suite
	mockClienteService *MockClienteService
	handler            *ClienteHandler
	router             *gin.Engine
}

// Configuração da suite antes de cada teste
func (s *ClienteHandlerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockClienteService = new(MockClienteService)
	s.handler = NewClienteHandler(s.mockClienteService)
	s.router = gin.Default()

	// Configurar rotas
	s.router.GET("/clientes", s.handler.GetAll)
	s.router.GET("/clientes/:id", s.handler.GetByID)
	s.router.POST("/clientes", s.handler.Create)
	s.router.PUT("/clientes/:id", s.handler.Update)
	s.router.DELETE("/clientes/:id", s.handler.Delete)
}

// Teste do método GetAll
func (s *ClienteHandlerSuite) TestGetAll() {
	// Criar clientes de exemplo
	now := time.Now()
	clientes := []domain.Cliente{
		{
			ID:            uuid.New(),
			NomeCompleto:  "João da Silva",
			Email:         "joao@exemplo.com",
			CPF:           "123.456.789-00",
			Telefone:      "(11) 98765-4321",
			DataNascimento: now.Add(-365 * 24 * time.Hour), // 1 ano atrás
			Ativo:         true,
			DataCadastro:  now.Add(-30 * 24 * time.Hour), // 30 dias atrás
		},
		{
			ID:            uuid.New(),
			NomeCompleto:  "Maria Oliveira",
			Email:         "maria@exemplo.com",
			CPF:           "987.654.321-00",
			Telefone:      "(11) 91234-5678",
			DataNascimento: now.Add(-730 * 24 * time.Hour), // 2 anos atrás
			Ativo:         true,
			DataCadastro:  now.Add(-15 * 24 * time.Hour), // 15 dias atrás
		},
	}

	// Configurar comportamento esperado do mock
	s.mockClienteService.On("GetAll").Return(clientes, nil)

	// Criar requisição e recorder para a resposta
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/clientes", nil)
	
	// Configurar o middleware de autenticação mock
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UsuarioIDKey, uuid.New())
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.GET("/clientes", s.handler.GetAll)
	
	router.ServeHTTP(w, req)

	// Verificar status code
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// Verificar resposta
	var responseGetAll []map[string]interface{}
	errGetAll := json.Unmarshal(w.Body.Bytes(), &responseGetAll)
	assert.Nil(s.T(), errGetAll)
	assert.Equal(s.T(), 2, len(responseGetAll))
	assert.Equal(s.T(), "João da Silva", responseGetAll[0]["nome_completo"])
	assert.Equal(s.T(), "123.456.789-00", responseGetAll[0]["cpf"]) // Verifica que o CPF foi retornado corretamente
	assert.Equal(s.T(), "Maria Oliveira", responseGetAll[1]["nome_completo"])

	// Verificar que o mock foi chamado conforme esperado
	s.mockClienteService.AssertExpectations(s.T())
}

// Teste do método GetByID
func (s *ClienteHandlerSuite) TestGetByID() {
	// Criar ID e cliente de exemplo
	id := uuid.New()
	now := time.Now()
	cliente := &domain.Cliente{
		ID:            id,
		NomeCompleto:  "José Santos",
		Email:         "jose@exemplo.com",
		CPF:           "123.456.789-00",
		Telefone:      "(11) 2345-6789",
		DataNascimento: now.Add(-40 * 365 * 24 * time.Hour), // 40 anos atrás
		Ativo:         true,
		DataCadastro:  now.Add(-60 * 24 * time.Hour), // 60 dias atrás
	}

	// Configurar comportamento esperado do mock
	s.mockClienteService.On("GetByID", id).Return(cliente, nil)

	// Configurar o router com middleware mock para simular autenticação
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UsuarioIDKey, uuid.New())
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.GET("/clientes/:id", s.handler.GetByID)
	
	// Criar requisição e recorder para a resposta
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/clientes/"+id.String(), nil)
	router.ServeHTTP(w, req)

	// Verificar status code
	assert.Equal(s.T(), http.StatusOK, w.Code)

	// Verificar resposta
	var responseGetByID map[string]interface{}
	errGetByID := json.Unmarshal(w.Body.Bytes(), &responseGetByID)
	assert.Nil(s.T(), errGetByID)
	assert.Equal(s.T(), id.String(), responseGetByID["id"])
	assert.Equal(s.T(), "José Santos", responseGetByID["nome_completo"])
	assert.Equal(s.T(), "123.456.789-00", responseGetByID["cpf"]) // Verifica que o CPF foi retornado corretamente
	assert.Equal(s.T(), "(11) 2345-6789", responseGetByID["telefone"])

	// Verificar que o mock foi chamado conforme esperado
	s.mockClienteService.AssertExpectations(s.T())
}

// Teste do método Create
func (s *ClienteHandlerSuite) TestCreate() {
	// Criar dados do cliente usando mapa para corresponder ao formato esperado pela API
	now := time.Now().Add(-25 * 365 * 24 * time.Hour) // 25 anos atrás
	
	clienteData := map[string]interface{}{
		"nome_completo":    "Pedro Almeida",
		"email":           "pedro@exemplo.com",
		"cpf":             "123.456.789-00", // CPF formatado
		"telefone":        "(11) 98765-4321", // Telefone formatado
		"data_nascimento": now.Format("2006-01-02"), // Formato ISO que o backend espera
		"ativo":           true,
	}
	
	// Converter para JSON
	bodyBytes, _ := json.Marshal(clienteData)

	// Configurar comportamento esperado do mock
	s.mockClienteService.On("Create", mock.AnythingOfType("*domain.Cliente")).Return(nil)

	// Configurar o middleware de autenticação mock
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UsuarioIDKey, uuid.New())
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.POST("/clientes", s.handler.Create)

	// Criar requisição e recorder para a resposta
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Verificar status code
	assert.Equal(s.T(), http.StatusCreated, w.Code)

	// Verificar que o mock foi chamado conforme esperado
	s.mockClienteService.AssertExpectations(s.T())
	
	// Verificar resposta
	var responseCreate map[string]interface{}
	errCreate := json.Unmarshal(w.Body.Bytes(), &responseCreate)
	assert.Nil(s.T(), errCreate)
	assert.NotEmpty(s.T(), responseCreate["id"]) // O ID deve ser gerado pelo servidor
}

// Executa a suite de testes
func TestClienteHandlerSuite(t *testing.T) {
	suite.Run(t, new(ClienteHandlerSuite))
}
