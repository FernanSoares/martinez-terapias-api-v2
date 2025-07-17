package services

import (
	"fmt"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockClienteRepo para os testes do ClienteService
type MockClienteRepositoryForService struct {
	mock.Mock
}

func (m *MockClienteRepositoryForService) FindAll() ([]domain.Cliente, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForService) FindByID(id uuid.UUID) (*domain.Cliente, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForService) FindByNome(nome string) ([]domain.Cliente, error) {
	args := m.Called(nome)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForService) FindByEmail(email string) (*domain.Cliente, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForService) FindByAtivo(ativo bool) ([]domain.Cliente, error) {
	args := m.Called(ativo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *MockClienteRepositoryForService) Create(cliente *domain.Cliente) error {
	args := m.Called(cliente)
	return args.Error(0)
}

func (m *MockClienteRepositoryForService) Update(cliente *domain.Cliente) error {
	args := m.Called(cliente)
	return args.Error(0)
}

func (m *MockClienteRepositoryForService) SoftDelete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// ClienteServiceSuite é a suite de testes para ClienteService
type ClienteServiceSuite struct {
	suite.Suite
	mockRepo *MockClienteRepositoryForService
	service  ports.ClienteService
}

// SetupTest é executado antes de cada teste
func (s *ClienteServiceSuite) SetupTest() {
	s.mockRepo = new(MockClienteRepositoryForService)
	s.service = NewClienteService(s.mockRepo)
}

// TestGetAll testa o método GetAll do ClienteService
func (s *ClienteServiceSuite) TestGetAll() {
	now := time.Now()
	clientes := []domain.Cliente{
		{
			ID:             uuid.New(),
			NomeCompleto:   "João Silva",
			CPF:            "123.456.789-00",
			Email:          "joao@exemplo.com",
			Telefone:       "(11) 98765-4321",
			DataNascimento: now.Add(-365 * 24 * time.Hour),
			DataCadastro:   now.Add(-30 * 24 * time.Hour),
			Ativo:          true,
		},
		{
			ID:             uuid.New(),
			NomeCompleto:   "Maria Oliveira",
			CPF:            "987.654.321-00",
			Email:          "maria@exemplo.com",
			Telefone:       "(11) 91234-5678",
			DataNascimento: now.Add(-730 * 24 * time.Hour), // 2 anos atrás
			DataCadastro:   now.Add(-15 * 24 * time.Hour),  // 15 dias atrás
			Ativo:          true,
		},
	}

	// Configurar comportamento esperado do mock
	s.mockRepo.On("FindAll").Return(clientes, nil)

	// Chamar o método a ser testado
	result, err := s.service.GetAll()

	// Verificar resultados
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(result))
	assert.Equal(s.T(), clientes[0].NomeCompleto, result[0].NomeCompleto)
	assert.Equal(s.T(), clientes[0].CPF, result[0].CPF)
	assert.Equal(s.T(), clientes[1].NomeCompleto, result[1].NomeCompleto)
	assert.Equal(s.T(), clientes[1].CPF, result[1].CPF)

	// Verificar que o mock foi chamado conforme esperado
	s.mockRepo.AssertExpectations(s.T())
}

// TestGetByID testa o método GetByID do ClienteService
func (s *ClienteServiceSuite) TestGetByID() {
	// Criar ID e cliente de exemplo
	id := uuid.New()
	now := time.Now()
	cliente := &domain.Cliente{
		ID:             id,
		NomeCompleto:   "José Santos",
		CPF:            "123.456.789-00",
		Email:          "jose@exemplo.com",
		Telefone:       "(11) 2345-6789",
		DataNascimento: now.Add(-40 * 365 * 24 * time.Hour), // 40 anos atrás
		DataCadastro:   now.Add(-60 * 24 * time.Hour),       // 60 dias atrás
		Ativo:          true,
	}

	// Configurar comportamento esperado do mock
	s.mockRepo.On("FindByID", id).Return(cliente, nil)

	// Chamar o método a ser testado
	result, err := s.service.GetByID(id)

	// Verificar resultados
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), id, result.ID)
	assert.Equal(s.T(), cliente.NomeCompleto, result.NomeCompleto)
	assert.Equal(s.T(), cliente.CPF, result.CPF)
	assert.Equal(s.T(), cliente.Email, result.Email)

	// Verificar que o mock foi chamado conforme esperado
	s.mockRepo.AssertExpectations(s.T())
}

// TestCreate testa o método Create do ClienteService
func (s *ClienteServiceSuite) TestCreate() {
	// Criar cliente de exemplo
	now := time.Now()
	cliente := &domain.Cliente{
		NomeCompleto:   "Pedro Almeida",
		CPF:            "123.456.789-00",
		Email:          "pedro@exemplo.com",
		Telefone:       "(11) 98765-4321",
		DataNascimento: now.Add(-25 * 365 * 24 * time.Hour), // 25 anos atrás
		Ativo:          true,
	}

	// Configurar comportamento esperado do mock
	// Primeiro deve verificar se o e-mail já existe
	s.mockRepo.On("FindByEmail", cliente.Email).Return(nil, fmt.Errorf("não encontrado"))
	// Depois criar o cliente
	s.mockRepo.On("Create", mock.AnythingOfType("*domain.Cliente")).Return(nil)

	// Chamar o método a ser testado
	err := s.service.Create(cliente)

	// Verificar resultados
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), uuid.Nil, cliente.ID)       // Verifica que um ID foi gerado
	assert.False(s.T(), cliente.DataCadastro.IsZero()) // Verifica que a data de cadastro foi definida

	// Verificar que o mock foi chamado conforme esperado
	s.mockRepo.AssertExpectations(s.T())
}

// TestUpdate testa o método Update do ClienteService
func (s *ClienteServiceSuite) TestUpdate() {
	// Usar um ID fixo para o teste
	id := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	now := time.Now()

	// Cliente com dados atualizados
	clienteAtualizado := &domain.Cliente{
		ID:             id,
		NomeCompleto:   "Ana Maria Silva", // Nome atualizado
		CPF:            "123.456.789-00",
		Email:          "anamaria@exemplo.com", // Email atualizado
		Telefone:       "(11) 99999-8888",      // Telefone atualizado
		DataNascimento: now.Add(-30 * 365 * 24 * time.Hour),
		DataCadastro:   now.Add(-90 * 24 * time.Hour),
		Ativo:          true,
	}

	// Configurar o mock - apenas Update é chamado, sem FindByID
	s.mockRepo.On("Update", mock.AnythingOfType("*domain.Cliente")).Return(nil)

	// Chamar o método a ser testado
	err := s.service.Update(clienteAtualizado)

	// Verificar resultados
	assert.Nil(s.T(), err)

	// Verificar que o mock foi chamado conforme esperado
	s.mockRepo.AssertExpectations(s.T())
}

// Execute a suite de testes
func TestClienteServiceSuite(t *testing.T) {
	suite.Run(t, new(ClienteServiceSuite))
}
