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

type MockAgendamentoService struct {
	mock.Mock
}

func (m *MockAgendamentoService) GetAll() ([]domain.Agendamento, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoService) GetByID(id uuid.UUID) (*domain.Agendamento, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoService) GetByClienteID(clienteID uuid.UUID) ([]domain.Agendamento, error) {
	args := m.Called(clienteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoService) GetByMassoterapeutaID(massoterapeutaID uuid.UUID) ([]domain.Agendamento, error) {
	args := m.Called(massoterapeutaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoService) GetByPeriodo(dataInicio, dataFim time.Time) ([]domain.Agendamento, error) {
	args := m.Called(dataInicio, dataFim)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoService) GetByStatus(status domain.StatusAgendamento) ([]domain.Agendamento, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoService) Create(agendamento *domain.Agendamento) error {
	args := m.Called(agendamento)
	return args.Error(0)
}

func (m *MockAgendamentoService) Update(agendamento *domain.Agendamento) error {
	args := m.Called(agendamento)
	return args.Error(0)
}

func (m *MockAgendamentoService) UpdateStatus(id uuid.UUID, status domain.StatusAgendamento) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockAgendamentoService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAgendamentoService) SolicitarReagendamento(id uuid.UUID, clienteID uuid.UUID) error {
	args := m.Called(id, clienteID)
	return args.Error(0)
}

type AgendamentoHandlerSuite struct {
	suite.Suite
	mockService *MockAgendamentoService
	handler     *AgendamentoHandler
	router      *gin.Engine
}

func (s *AgendamentoHandlerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockService = new(MockAgendamentoService)
	s.handler = NewAgendamentoHandler(s.mockService)
	s.router = gin.Default()

	s.router.GET("/agendamentos", s.handler.GetAll)
	s.router.GET("/agendamentos/:id", s.handler.GetByID)
	s.router.POST("/agendamentos", s.handler.Create)
	s.router.PUT("/agendamentos/:id", s.handler.Update)
	s.router.DELETE("/agendamentos/:id", s.handler.Delete)
	s.router.PATCH("/agendamentos/:id/status", s.handler.UpdateStatus)
	s.router.POST("/agendamentos/:id/solicitar-reagendamento", s.handler.SolicitarReagendamento)
}

func (s *AgendamentoHandlerSuite) TestGetAll() {
	now := time.Now()
	agendamentos := []domain.Agendamento{
		{
			ID:               uuid.New(),
			ClienteID:        uuid.New(),
			MassoterapeutaID: uuid.New(),
			ServicoID:        uuid.New(),
			DataHora:         now.Add(24 * time.Hour),
			Status:           domain.StatusConfirmado,
		},
		{
			ID:               uuid.New(),
			ClienteID:        uuid.New(),
			MassoterapeutaID: uuid.New(),
			ServicoID:        uuid.New(),
			DataHora:         now.Add(48 * time.Hour),
			Status:           domain.StatusAgendado,
		},
	}

	s.mockService.On("GetAll").Return(agendamentos, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/agendamentos", nil)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UsuarioIDKey, uuid.New())
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.GET("/agendamentos", s.handler.GetAll)

	router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(response))
	assert.Equal(s.T(), string(domain.StatusConfirmado), response[0]["status"])
	assert.Equal(s.T(), string(domain.StatusAgendado), response[1]["status"])
}

func (s *AgendamentoHandlerSuite) TestGetByID() {

	id := uuid.New()
	clienteID := uuid.New()
	massoterapeutaID := uuid.New()
	servicoID := uuid.New()
	now := time.Now()
	agendamento := &domain.Agendamento{
		ID:               id,
		ClienteID:        clienteID,
		MassoterapeutaID: massoterapeutaID,
		ServicoID:        servicoID,
		DataHora:         now.Add(24 * time.Hour),
		Status:           domain.StatusConfirmado,
	}

	s.mockService.On("GetByID", id).Return(agendamento, nil)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UsuarioIDKey, uuid.New())
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.GET("/agendamentos/:id", s.handler.GetByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/agendamentos/"+id.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), id.String(), response["id"])
	assert.Equal(s.T(), string(domain.StatusConfirmado), response["status"])

	s.mockService.AssertExpectations(s.T())
}

func (s *AgendamentoHandlerSuite) TestCreate() {

	clienteID := uuid.New()
	massoterapeutaID := uuid.New()
	servicoID := uuid.New()
	agendamentoID := uuid.New()
	dataHora := time.Now().Add(24 * time.Hour)

	agendamento := &domain.Agendamento{
		ID:               agendamentoID,
		ClienteID:        clienteID,
		MassoterapeutaID: massoterapeutaID,
		ServicoID:        servicoID,
		DataHora:         dataHora,
		Status:           domain.StatusAgendado,
		Observacoes:      "Observação de teste",
	}

	bodyBytes, _ := json.Marshal(agendamento)

	s.mockService.On("Create", mock.AnythingOfType("*domain.Agendamento")).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/agendamentos", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UsuarioIDKey, uuid.New())
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.POST("/agendamentos", s.handler.Create)

	router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	s.mockService.AssertExpectations(s.T())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(s.T(), err)
	assert.Contains(s.T(), response, "id")
}

func (s *AgendamentoHandlerSuite) TestUpdate() {
	agendamentoID := uuid.New()
	clienteID := uuid.New()
	massoterapeutaID := uuid.New()
	servicoID := uuid.New()
	dataHora := time.Now().Add(24 * time.Hour)

	agendamento := &domain.Agendamento{
		ID:               agendamentoID,
		ClienteID:        clienteID,
		MassoterapeutaID: massoterapeutaID,
		ServicoID:        servicoID,
		DataHora:         dataHora,
		Status:           domain.StatusConfirmado,
		Observacoes:      "Observação atualizada",
	}

	s.mockService.On("GetByID", agendamentoID).Return(agendamento, nil)
	s.mockService.On("Update", mock.AnythingOfType("*domain.Agendamento")).Return(nil)

	bodyBytes, _ := json.Marshal(agendamento)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UsuarioIDKey, uuid.New())
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.PUT("/agendamentos/:id", s.handler.Update)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/agendamentos/"+agendamentoID.String(), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	s.mockService.AssertExpectations(s.T())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(s.T(), err)
	assert.Contains(s.T(), response, "id")
}

func (s *AgendamentoHandlerSuite) TestUpdateStatus() {
	agendamentoID := uuid.New()
	novoStatus := domain.StatusRealizado

	requestBody := map[string]string{
		"status": string(novoStatus),
	}
	bodyBytes, _ := json.Marshal(requestBody)

	s.mockService.On("UpdateStatus", agendamentoID, novoStatus).Return(nil)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UsuarioIDKey, uuid.New())
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.PATCH("/agendamentos/:id/status", s.handler.UpdateStatus)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/agendamentos/"+agendamentoID.String()+"/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	s.mockService.AssertExpectations(s.T())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(s.T(), err)
	assert.Contains(s.T(), response, "message")
}

func (s *AgendamentoHandlerSuite) TestDelete() {
	agendamentoID := uuid.New()

	s.mockService.On("Delete", agendamentoID).Return(nil)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UsuarioIDKey, uuid.New())
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.DELETE("/agendamentos/:id", s.handler.Delete)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/agendamentos/"+agendamentoID.String(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusNoContent, w.Code)

	s.mockService.AssertExpectations(s.T())
}

func (s *AgendamentoHandlerSuite) TestSolicitarReagendamento() {

	agendamentoID := uuid.New()
	usuarioID := uuid.New()
	clienteID := uuid.New()

	agendamento := &domain.Agendamento{
		ID:        agendamentoID,
		ClienteID: clienteID,
		DataHora:  time.Now(),
		Status:    domain.StatusConfirmado,
	}

	requestBody := map[string]string{}
	bodyBytes, _ := json.Marshal(requestBody)

	s.mockService.On("GetByID", agendamentoID).Return(agendamento, nil)
	s.mockService.On("SolicitarReagendamento", agendamentoID, clienteID).Return(nil)

	router := gin.Default()
	router.Use(func(c *gin.Context) {

		c.Set(middleware.UsuarioIDKey, usuarioID)
		c.Set(middleware.PerfilKey, "admin")
		c.Next()
	})
	router.POST("/agendamentos/:id/solicitar-reagendamento", s.handler.SolicitarReagendamento)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/agendamentos/"+agendamentoID.String()+"/solicitar-reagendamento", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	s.mockService.AssertExpectations(s.T())
}

func TestAgendamentoHandlerSuite(t *testing.T) {
	suite.Run(t, new(AgendamentoHandlerSuite))
}
