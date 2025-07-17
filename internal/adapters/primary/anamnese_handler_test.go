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

// MockAnamneseService implementa um mock do AnamneseService para testes
type MockAnamneseService struct {
	mock.Mock
}

func (m *MockAnamneseService) GetAll() ([]domain.Anamnese, error) {
	args := m.Called()
	return args.Get(0).([]domain.Anamnese), args.Error(1)
}

func (m *MockAnamneseService) GetByID(id uuid.UUID) (*domain.Anamnese, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Anamnese), args.Error(1)
}

func (m *MockAnamneseService) GetByClienteID(clienteID uuid.UUID) ([]domain.Anamnese, error) {
	args := m.Called(clienteID)
	return args.Get(0).([]domain.Anamnese), args.Error(1)
}

func (m *MockAnamneseService) Create(anamnese *domain.Anamnese) error {
	args := m.Called(anamnese)
	return args.Error(0)
}

func (m *MockAnamneseService) Update(anamnese *domain.Anamnese) error {
	args := m.Called(anamnese)
	return args.Error(0)
}

func (m *MockAnamneseService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// AnamneseHandlerSuite define a suite de testes para AnamneseHandler
type AnamneseHandlerSuite struct {
	suite.Suite
	mockService *MockAnamneseService
	handler     *AnamneseHandler
	router      *gin.Engine
}

// SetupTest prepara o ambiente para cada teste
func (s *AnamneseHandlerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockService = new(MockAnamneseService)
	s.handler = NewAnamneseHandler(s.mockService)
	s.router = gin.New()
}

// TestGetAll testa o endpoint de listar todas as anamneses
func (s *AnamneseHandlerSuite) TestGetAll() {
	anamneses := []domain.Anamnese{
		{
			ID:        uuid.New(),
			ClienteID: uuid.New(),
		},
		{
			ID:        uuid.New(),
			ClienteID: uuid.New(),
		},
	}

	s.mockService.On("GetAll").Return(anamneses, nil)

	s.router.GET("/anamnese", s.handler.GetAll)
	req, _ := http.NewRequest("GET", "/anamnese", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)

	var respAnamneses []domain.Anamnese
	err := json.Unmarshal(resp.Body.Bytes(), &respAnamneses)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), len(anamneses), len(respAnamneses))
}

// TestGetAll_Error testa o erro no endpoint de listar todas as anamneses
func (s *AnamneseHandlerSuite) TestGetAll_Error() {
	s.mockService.On("GetAll").Return([]domain.Anamnese{}, errors.New("erro ao buscar anamneses"))

	s.router.GET("/anamnese", s.handler.GetAll)
	req, _ := http.NewRequest("GET", "/anamnese", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
}

// TestGetByID testa o endpoint de buscar anamnese por ID
func (s *AnamneseHandlerSuite) TestGetByID() {
	id := uuid.New()
	anamnese := &domain.Anamnese{
		ID:        id,
		ClienteID: uuid.New(),
	}

	s.mockService.On("GetByID", id).Return(anamnese, nil)

	s.router.GET("/anamnese/:id", s.handler.GetByID)
	req, _ := http.NewRequest("GET", "/anamnese/"+id.String(), nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)

	var respAnamnese domain.Anamnese
	err := json.Unmarshal(resp.Body.Bytes(), &respAnamnese)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), anamnese.ID, respAnamnese.ID)
}

// TestGetByID_InvalidID testa o erro de ID inválido
func (s *AnamneseHandlerSuite) TestGetByID_InvalidID() {
	s.router.GET("/anamnese/:id", s.handler.GetByID)
	req, _ := http.NewRequest("GET", "/anamnese/not-a-uuid", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

// TestGetByID_NotFound testa o erro de anamnese não encontrada
func (s *AnamneseHandlerSuite) TestGetByID_NotFound() {
	id := uuid.New()
	s.mockService.On("GetByID", id).Return(nil, errors.New("anamnese não encontrada"))

	s.router.GET("/anamnese/:id", s.handler.GetByID)
	req, _ := http.NewRequest("GET", "/anamnese/"+id.String(), nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
}

// TestGetByClienteID testa o endpoint de buscar anamnese por cliente ID
func (s *AnamneseHandlerSuite) TestGetByClienteID() {
	clienteID := uuid.New()
	anamneses := []domain.Anamnese{
		{
			ID:        uuid.New(),
			ClienteID: clienteID,
		},
	}

	s.mockService.On("GetByClienteID", clienteID).Return(anamneses, nil)

	s.router.GET("/clientes-anamnese/:cliente_id", s.handler.GetByClienteID)
	req, _ := http.NewRequest("GET", "/clientes-anamnese/"+clienteID.String(), nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)

	var respAnamneses []domain.Anamnese
	err := json.Unmarshal(resp.Body.Bytes(), &respAnamneses)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), len(anamneses), len(respAnamneses))
	assert.Equal(s.T(), anamneses[0].ID, respAnamneses[0].ID)
	assert.Equal(s.T(), clienteID, respAnamneses[0].ClienteID)
}

// TestGetByClienteID_InvalidID testa o erro de ID de cliente inválido
func (s *AnamneseHandlerSuite) TestGetByClienteID_InvalidID() {
	s.router.GET("/clientes-anamnese/:cliente_id", s.handler.GetByClienteID)
	req, _ := http.NewRequest("GET", "/clientes-anamnese/not-a-uuid", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

// TestGetByClienteID_NotFound testa o erro de anamnese de cliente não encontrada
func (s *AnamneseHandlerSuite) TestGetByClienteID_NotFound() {
	clienteID := uuid.New()
	s.mockService.On("GetByClienteID", clienteID).Return([]domain.Anamnese{}, errors.New("anamnese não encontrada"))

	s.router.GET("/clientes-anamnese/:cliente_id", s.handler.GetByClienteID)
	req, _ := http.NewRequest("GET", "/clientes-anamnese/"+clienteID.String(), nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
}

// TestCreate testa o endpoint de criação de anamnese
func (s *AnamneseHandlerSuite) TestCreate() {
	clienteID := uuid.New()
	anamnese := domain.Anamnese{
		ClienteID:        clienteID,
		QueixaPrincipal:   "Dor nas costas",
		HistoricoDoencas:  "Hipertensão",
		MedicamentosEmUso: "Losartana",
		Alergias:          "Nenhuma",
		CirurgiasPrevias:  "Nenhuma",
	}

	s.mockService.On("Create", mock.AnythingOfType("*domain.Anamnese")).Return(nil)

	s.router.POST("/clientes-anamnese/:cliente_id", s.handler.Create)
	body, _ := json.Marshal(anamnese)
	req, _ := http.NewRequest("POST", "/clientes-anamnese/"+clienteID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusCreated, resp.Code)
}

// TestCreate_InvalidID testa o erro de ID inválido na criação
func (s *AnamneseHandlerSuite) TestCreate_InvalidID() {
	anamnese := domain.Anamnese{
		QueixaPrincipal: "Dor nas costas",
	}

	s.router.POST("/clientes-anamnese/:cliente_id", s.handler.Create)
	body, _ := json.Marshal(anamnese)
	req, _ := http.NewRequest("POST", "/clientes-anamnese/not-a-uuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

// TestCreate_Error testa o erro na criação de anamnese
func (s *AnamneseHandlerSuite) TestCreate_Error() {
	clienteID := uuid.New()
	anamnese := domain.Anamnese{
		ClienteID:       clienteID,
		QueixaPrincipal: "Dor nas costas",
	}

	s.mockService.On("Create", mock.AnythingOfType("*domain.Anamnese")).Return(errors.New("erro ao criar anamnese"))

	s.router.POST("/clientes-anamnese/:cliente_id", s.handler.Create)
	body, _ := json.Marshal(anamnese)
	req, _ := http.NewRequest("POST", "/clientes-anamnese/"+clienteID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

// TestUpdate testa o endpoint de atualização de anamnese
func (s *AnamneseHandlerSuite) TestUpdate() {
	id := uuid.New()
	anamnese := domain.Anamnese{
		ID:              id,
		ClienteID:       uuid.New(),
		QueixaPrincipal: "Dor nas costas atualizada",
	}

	s.mockService.On("GetByID", id).Return(&anamnese, nil)
	s.mockService.On("Update", mock.AnythingOfType("*domain.Anamnese")).Return(nil)

	s.router.PUT("/anamnese/:id", s.handler.Update)
	body, _ := json.Marshal(anamnese)
	req, _ := http.NewRequest("PUT", "/anamnese/"+id.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

// TestAnamneseHandlerSuite executa a suite de testes
func TestAnamneseHandlerSuite(t *testing.T) {
	suite.Run(t, new(AnamneseHandlerSuite))
}
