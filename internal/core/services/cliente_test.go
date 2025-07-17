package services

import (
	"errors"
	"martinezterapias/api/internal/core/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// SimpleMockClienteRepo é um mock simples para testes
type SimpleMockClienteRepo struct {
	mock.Mock
}

func (m *SimpleMockClienteRepo) FindAll() ([]domain.Cliente, error) {
	args := m.Called()
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *SimpleMockClienteRepo) FindByID(id uuid.UUID) (*domain.Cliente, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cliente), args.Error(1)
}

func (m *SimpleMockClienteRepo) FindByNome(nome string) ([]domain.Cliente, error) {
	args := m.Called(nome)
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *SimpleMockClienteRepo) FindByEmail(email string) (*domain.Cliente, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cliente), args.Error(1)
}

func (m *SimpleMockClienteRepo) FindByAtivo(ativo bool) ([]domain.Cliente, error) {
	args := m.Called(ativo)
	return args.Get(0).([]domain.Cliente), args.Error(1)
}

func (m *SimpleMockClienteRepo) Create(cliente *domain.Cliente) error {
	args := m.Called(cliente)
	if cliente.ID == uuid.Nil {
		cliente.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *SimpleMockClienteRepo) Update(cliente *domain.Cliente) error {
	args := m.Called(cliente)
	return args.Error(0)
}

func (m *SimpleMockClienteRepo) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *SimpleMockClienteRepo) SoftDelete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestSimpleClienteService é uma função simples para testar o serviço de cliente
func TestSimpleClienteService(t *testing.T) {
	// Criar mock
	mockRepo := new(SimpleMockClienteRepo)
	
	// Criar o serviço com o mock
	service := NewClienteService(mockRepo)
	
	// Teste do GetAll
	t.Run("TestGetAll", func(t *testing.T) {
		// Dados de exemplo
		clientes := []domain.Cliente{
			{
				ID:            uuid.New(),
				NomeCompleto:  "Cliente 1",
				Email:         "cliente1@exemplo.com",
				Telefone:      "1199998888",
				DataNascimento: time.Now().AddDate(-30, 0, 0),
				DataCadastro:  time.Now(),
				Ativo:         true,
			},
			{
				ID:            uuid.New(),
				NomeCompleto:  "Cliente 2",
				Email:         "cliente2@exemplo.com",
				Telefone:      "1199997777",
				DataNascimento: time.Now().AddDate(-25, 0, 0),
				DataCadastro:  time.Now(),
				Ativo:         true,
			},
		}
		
		// Configurar o mock
		mockRepo.On("FindAll").Return(clientes, nil)
		
		// Executar
		result, err := service.GetAll()
		
		// Verificar resultado
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, clientes[0].NomeCompleto, result[0].NomeCompleto)
		assert.Equal(t, clientes[1].NomeCompleto, result[1].NomeCompleto)
		
		// Verificar chamada do mock
		mockRepo.AssertExpectations(t)
	})
	
	// Teste do GetByID
	t.Run("TestGetByID", func(t *testing.T) {
		// Dados de exemplo
		id := uuid.New()
		cliente := &domain.Cliente{
			ID:            id,
			NomeCompleto:  "Cliente Teste",
			Email:         "teste@exemplo.com",
			Telefone:      "1199998888",
			DataNascimento: time.Now().AddDate(-30, 0, 0),
			DataCadastro:  time.Now(),
			Ativo:         true,
		}
		
		// Configurar o mock
		mockRepo.On("FindByID", id).Return(cliente, nil)
		
		// Executar
		result, err := service.GetByID(id)
		
		// Verificar resultado
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, cliente.ID, result.ID)
		assert.Equal(t, cliente.NomeCompleto, result.NomeCompleto)
		
		// Verificar chamada do mock
		mockRepo.AssertExpectations(t)
	})
	
	// Teste do Create
	t.Run("TestCreate", func(t *testing.T) {
		// Dados de exemplo
		cliente := &domain.Cliente{
			NomeCompleto:  "Novo Cliente",
			Email:         "novo@exemplo.com",
			Telefone:      "1199998888",
			DataNascimento: time.Now().AddDate(-28, 0, 0),
		}
		
		// Configurar o mock para verificar se o email já existe
		mockRepo.On("FindByEmail", cliente.Email).Return(nil, errors.New("não encontrado"))
		
		// Configurar o mock para o Create
		mockRepo.On("Create", cliente).Return(nil).Run(func(args mock.Arguments) {
			c := args.Get(0).(*domain.Cliente)
			c.ID = uuid.New()
			c.DataCadastro = time.Now()
		})
		
		// Executar
		err := service.Create(cliente)
		
		// Verificar resultado
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, cliente.ID)
		
		// Verificar chamada do mock
		mockRepo.AssertExpectations(t)
	})
	
	// Teste do Create com email duplicado
	t.Run("TestCreate_DuplicateEmail", func(t *testing.T) {
		// Dados de exemplo
		email := "existente@exemplo.com"
		clienteExistente := &domain.Cliente{
			ID:           uuid.New(),
			NomeCompleto: "Cliente Existente",
			Email:        email,
		}
		
		novoCliente := &domain.Cliente{
			NomeCompleto: "Novo Cliente",
			Email:        email,
			Telefone:     "1199998888",
		}
		
		// Configurar o mock para retornar um cliente existente
		mockRepo.On("FindByEmail", email).Return(clienteExistente, nil)
		
		// Executar
		err := service.Create(novoCliente)
		
		// Verificar resultado (deve ter erro de email duplicado)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "já está em uso")
		
		// Verificar chamada do mock
		mockRepo.AssertExpectations(t)
	})
	
	// Teste do Update
	t.Run("TestUpdate", func(t *testing.T) {
		// Dados de exemplo
		id := uuid.New()
		cliente := &domain.Cliente{
			ID:           id,
			NomeCompleto: "Cliente Atualizado",
			Email:        "atualizado@exemplo.com",
			Telefone:     "1199998888",
		}
		
		// Configurar o mock
		// Note: A implementação atual não verifica a existência do cliente antes de atualizar
		mockRepo.On("Update", cliente).Return(nil)
		
		// Executar
		err := service.Update(cliente)
		
		// Verificar resultado
		assert.NoError(t, err)
		
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("TestDelete", func(t *testing.T) {
		id := uuid.New()
		
		mockRepo.On("SoftDelete", id).Return(nil)
		
		err := service.Delete(id)
		
		assert.NoError(t, err)
		
		mockRepo.AssertExpectations(t)
	})
}
