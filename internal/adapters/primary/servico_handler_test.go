package primary

import (
	"bytes"
	"encoding/json"
	"errors"
	"martinezterapias/api/internal/core/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockServicoService struct {
	mock.Mock
}

func (m *MockServicoService) GetAll() ([]domain.Servico, error) {
	args := m.Called()
	return args.Get(0).([]domain.Servico), args.Error(1)
}

func (m *MockServicoService) GetByID(id uuid.UUID) (*domain.Servico, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Servico), args.Error(1)
}

func (m *MockServicoService) Create(servico *domain.Servico) error {
	args := m.Called(servico)
	return args.Error(0)
}

func (m *MockServicoService) Update(servico *domain.Servico) error {
	args := m.Called(servico)
	return args.Error(0)
}

func (m *MockServicoService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type ServicoHandlerSuite struct {
	suite.Suite
	mockService *MockServicoService
	handler     *ServicoHandler
	router      *gin.Engine
}

func (s *ServicoHandlerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockService = new(MockServicoService)
	s.handler = NewServicoHandler(s.mockService)
	s.router = gin.New()
}

func (s *ServicoHandlerSuite) TestGetAll() {
	servicos := []domain.Servico{
		{
			ID:             uuid.New(),
			NomeServico:    "Serviço 1",
			Descricao:      "Descrição 1",
			DuracaoMinutos: 60,
			Valor:          100.0,
		},
		{
			ID:             uuid.New(),
			NomeServico:    "Serviço 2",
			Descricao:      "Descrição 2",
			DuracaoMinutos: 90,
			Valor:          150.0,
		},
	}

	s.mockService.On("GetAll").Return(servicos, nil)

	s.router.GET("/servicos", s.handler.GetAll)
	req, _ := http.NewRequest("GET", "/servicos", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)

	var respServicos []domain.Servico
	err := json.Unmarshal(resp.Body.Bytes(), &respServicos)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), len(servicos), len(respServicos))
}

func (s *ServicoHandlerSuite) TestGetAll_Error() {
	s.mockService.On("GetAll").Return([]domain.Servico{}, errors.New("erro de serviço"))

	s.router.GET("/servicos", s.handler.GetAll)
	req, _ := http.NewRequest("GET", "/servicos", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
}

func (s *ServicoHandlerSuite) TestGetByID() {
	id := uuid.New()
	servico := &domain.Servico{
		ID:             id,
		NomeServico:    "Serviço 1",
		Descricao:      "Descrição 1",
		DuracaoMinutos: 60,
		Valor:          100.0,
	}

	s.mockService.On("GetByID", id).Return(servico, nil)

	s.router.GET("/servicos/:id", s.handler.GetByID)
	req, _ := http.NewRequest("GET", "/servicos/"+id.String(), nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)

	var respServico domain.Servico
	err := json.Unmarshal(resp.Body.Bytes(), &respServico)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), servico.ID, respServico.ID)
	assert.Equal(s.T(), servico.NomeServico, respServico.NomeServico)
}

func (s *ServicoHandlerSuite) TestGetByID_InvalidID() {
	s.router.GET("/servicos/:id", s.handler.GetByID)
	req, _ := http.NewRequest("GET", "/servicos/not-a-uuid", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

func (s *ServicoHandlerSuite) TestGetByID_NotFound() {
	id := uuid.New()
	s.mockService.On("GetByID", id).Return(nil, errors.New("serviço não encontrado"))

	s.router.GET("/servicos/:id", s.handler.GetByID)
	req, _ := http.NewRequest("GET", "/servicos/"+id.String(), nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
}

func (s *ServicoHandlerSuite) TestCreate() {
	servico := domain.Servico{
		NomeServico:    "Novo Serviço",
		Descricao:      "Descrição do novo serviço",
		DuracaoMinutos: 60,
		Valor:          100.0,
	}

	s.mockService.On("Create", mock.AnythingOfType("*domain.Servico")).Return(nil)

	s.router.POST("/servicos", s.handler.Create)
	body, _ := json.Marshal(servico)
	req, _ := http.NewRequest("POST", "/servicos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusCreated, resp.Code)
}

func (s *ServicoHandlerSuite) TestCreate_ValidationErrors() {
	servico := domain.Servico{
		Descricao:      "Descrição do novo serviço",
		DuracaoMinutos: 60,
		Valor:          100.0,
	}

	s.router.POST("/servicos", s.handler.Create)
	body, _ := json.Marshal(servico)
	req, _ := http.NewRequest("POST", "/servicos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)

	servico = domain.Servico{
		NomeServico: "Nome do serviço",
		Descricao:   "Descrição do novo serviço",
		Valor:       100.0,
	}

	body, _ = json.Marshal(servico)
	req, _ = http.NewRequest("POST", "/servicos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)

	servico = domain.Servico{
		NomeServico:    "Nome do serviço",
		Descricao:      "Descrição do novo serviço",
		DuracaoMinutos: 60,
	}

	body, _ = json.Marshal(servico)
	req, _ = http.NewRequest("POST", "/servicos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

func (s *ServicoHandlerSuite) TestUpdate() {
	id := uuid.New()
	servico := domain.Servico{
		ID:             id,
		NomeServico:    "Serviço Atualizado",
		Descricao:      "Descrição atualizada",
		DuracaoMinutos: 90,
		Valor:          120.0,
	}

	s.mockService.On("GetByID", id).Return(&servico, nil)
	s.mockService.On("Update", mock.AnythingOfType("*domain.Servico")).Return(nil)

	s.router.PUT("/servicos/:id", s.handler.Update)
	body, _ := json.Marshal(servico)
	req, _ := http.NewRequest("PUT", "/servicos/"+id.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

func (s *ServicoHandlerSuite) TestDelete() {
	id := uuid.New()
	s.mockService.On("Delete", id).Return(nil)

	s.router.DELETE("/servicos/:id", s.handler.Delete)
	req, _ := http.NewRequest("DELETE", "/servicos/"+id.String(), nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNoContent, resp.Code)
}

func (s *ServicoHandlerSuite) TestDelete_Error() {
	id := uuid.New()
	s.mockService.On("Delete", id).Return(errors.New("erro ao excluir"))

	s.router.DELETE("/servicos/:id", s.handler.Delete)
	req, _ := http.NewRequest("DELETE", "/servicos/"+id.String(), nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
}

func TestServicoHandlerSuite(t *testing.T) {
	suite.Run(t, new(ServicoHandlerSuite))
}
