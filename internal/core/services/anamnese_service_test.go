package services

import (
	"errors"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockAnamneseRepository implementa um mock do repositório de anamneses
type MockAnamneseRepository struct {
	mock.Mock
}

func (m *MockAnamneseRepository) FindAll() ([]domain.Anamnese, error) {
	args := m.Called()
	return args.Get(0).([]domain.Anamnese), args.Error(1)
}

func (m *MockAnamneseRepository) FindByID(id uuid.UUID) (*domain.Anamnese, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Anamnese), args.Error(1)
}

func (m *MockAnamneseRepository) FindByClienteID(clienteID uuid.UUID) ([]domain.Anamnese, error) {
	args := m.Called(clienteID)
	return args.Get(0).([]domain.Anamnese), args.Error(1)
}

func (m *MockAnamneseRepository) Create(anamnese *domain.Anamnese) error {
	args := m.Called(anamnese)
	return args.Error(0)
}

func (m *MockAnamneseRepository) Update(anamnese *domain.Anamnese) error {
	args := m.Called(anamnese)
	return args.Error(0)
}

// MockClienteRepositoryForAnamnese implementa um mock do repositório de clientes para o serviço de anamnese
type MockClienteRepositoryForAnamnese struct {
	mock.Mock
}

func (m *MockClienteRepositoryForAnamnese) FindAll() ([]domain.Cliente, error) {
	args := m.Called()
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForAnamnese) FindByID(id uuid.UUID) (*domain.Cliente, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForAnamnese) FindByEmail(email string) (*domain.Cliente, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForAnamnese) FindByNome(nome string) ([]domain.Cliente, error) {
	args := m.Called(nome)
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForAnamnese) FindByAtivo(ativo bool) ([]domain.Cliente, error) {
	args := m.Called(ativo)
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForAnamnese) Create(cliente *domain.Cliente) error {
	args := m.Called(cliente)
	return args.Error(0)
}

func (m *MockClienteRepositoryForAnamnese) Update(cliente *domain.Cliente) error {
	args := m.Called(cliente)
	return args.Error(0)
}

func (m *MockClienteRepositoryForAnamnese) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockClienteRepositoryForAnamnese) SoftDelete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// AnamneseServiceSuite define a suite de testes para AnamneseService
type AnamneseServiceSuite struct {
	suite.Suite
	mockAnamneseRepo *MockAnamneseRepository
	mockClienteRepo  *MockClienteRepositoryForAnamnese
	service          ports.AnamneseService
}

// SetupTest prepara o ambiente antes de cada teste
func (s *AnamneseServiceSuite) SetupTest() {
	s.mockAnamneseRepo = new(MockAnamneseRepository)
	s.mockClienteRepo = new(MockClienteRepositoryForAnamnese)
	s.service = NewAnamneseService(s.mockAnamneseRepo, s.mockClienteRepo)
}

// TestGetAll testa o método GetAll
func (s *AnamneseServiceSuite) TestGetAll() {
	anamneses := []domain.Anamnese{
		{
			ID:                uuid.New(),
			ClienteID:         uuid.New(),
			DataPreenchimento: time.Now(),
		},
		{
			ID:                uuid.New(),
			ClienteID:         uuid.New(),
			DataPreenchimento: time.Now(),
		},
	}

	s.mockAnamneseRepo.On("FindAll").Return(anamneses, nil)

	result, err := s.service.GetAll()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(result))
}

// TestGetByID testa o método GetByID
func (s *AnamneseServiceSuite) TestGetByID() {
	id := uuid.New()
	clienteID := uuid.New()
	anamnese := &domain.Anamnese{
		ID:                id,
		ClienteID:         clienteID,
		DataPreenchimento: time.Now(),
	}

	s.mockAnamneseRepo.On("FindByID", id).Return(anamnese, nil)

	result, err := s.service.GetByID(id)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), id, result.ID)
	assert.Equal(s.T(), clienteID, result.ClienteID)
}

// TestGetByClienteID testa o método GetByClienteID
func (s *AnamneseServiceSuite) TestGetByClienteID() {
	clienteID := uuid.New()
	anamneses := []domain.Anamnese{
		{
			ID:                uuid.New(),
			ClienteID:         clienteID,
			DataPreenchimento: time.Now(),
		},
	}

	s.mockAnamneseRepo.On("FindByClienteID", clienteID).Return(anamneses, nil)

	result, err := s.service.GetByClienteID(clienteID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), clienteID, result[0].ClienteID)
}

// TestCreate testa o método de criação de anamnese
func (s *AnamneseServiceSuite) TestCreate() {
	clienteID := uuid.New()
	cliente := &domain.Cliente{
		ID:           clienteID,
		NomeCompleto: "Teste Cliente",
		Email:        "cliente@teste.com",
	}
	
	anamnese := &domain.Anamnese{
		ClienteID:         clienteID,
		DataPreenchimento: time.Now(),
		QueixaPrincipal:   "Dores nas costas",
		HistoricoDoencas:  "Nenhum",
		CirurgiasPrevias:  "Nenhuma",
		MedicamentosEmUso: "Nenhum",
	}

	s.mockClienteRepo.On("FindByID", clienteID).Return(cliente, nil)
	s.mockAnamneseRepo.On("Create", mock.AnythingOfType("*domain.Anamnese")).Return(nil)

	err := s.service.Create(anamnese)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), uuid.Nil, anamnese.ID) // Verifica se um ID foi gerado
}

// TestCreate_ClienteNaoEncontrado testa a criação quando o cliente não existe
func (s *AnamneseServiceSuite) TestCreate_ClienteNaoEncontrado() {
	clienteID := uuid.New()
	anamnese := &domain.Anamnese{
		ClienteID:         clienteID,
		DataPreenchimento: time.Now(),
	}

	s.mockClienteRepo.On("FindByID", clienteID).Return(nil, errors.New("cliente não encontrado"))

	err := s.service.Create(anamnese)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "cliente não encontrado")
}

// TestUpdate testa o método de atualização de anamnese
func (s *AnamneseServiceSuite) TestUpdate() {
	id := uuid.New()
	clienteID := uuid.New()
	anamnese := &domain.Anamnese{
		ID:                id,
		ClienteID:         clienteID,
		DataPreenchimento: time.Now(),
		QueixaPrincipal:   "Dores atualizadas",
		HistoricoDoencas:  "Atualizado",
		CirurgiasPrevias:  "Nenhuma",
		MedicamentosEmUso: "Atualizado",
	}

	s.mockAnamneseRepo.On("Update", anamnese).Return(nil)

	err := s.service.Update(anamnese)
	assert.Nil(s.T(), err)
}

// TestAnamneseServiceSuite executa a suite de testes
func TestAnamneseServiceSuite(t *testing.T) {
	suite.Run(t, new(AnamneseServiceSuite))
}
