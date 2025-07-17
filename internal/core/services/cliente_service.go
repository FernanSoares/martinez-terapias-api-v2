package services

import (
	"fmt"
	"time"

	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"

	"github.com/google/uuid"
)

// clienteService implementa a interface ClienteService
type clienteService struct {
	repo ports.ClienteRepository
}

// NewClienteService cria uma nova instância de ClienteService
func NewClienteService(repo ports.ClienteRepository) ports.ClienteService {
	return &clienteService{
		repo: repo,
	}
}

// GetAll retorna todos os clientes
func (s *clienteService) GetAll() ([]domain.Cliente, error) {
	return s.repo.FindAll()
}

// GetByID retorna um cliente pelo ID
func (s *clienteService) GetByID(id uuid.UUID) (*domain.Cliente, error) {
	return s.repo.FindByID(id)
}

// GetByNome retorna clientes pelo nome (busca parcial)
func (s *clienteService) GetByNome(nome string) ([]domain.Cliente, error) {
	return s.repo.FindByNome(nome)
}

// GetByAtivo retorna clientes pelo status ativo/inativo
func (s *clienteService) GetByAtivo(ativo bool) ([]domain.Cliente, error) {
	return s.repo.FindByAtivo(ativo)
}

// Create cria um novo cliente
func (s *clienteService) Create(cliente *domain.Cliente) error {
	// Verificar se o email já existe
	if cliente.Email != "" {
		clienteExistente, err := s.repo.FindByEmail(cliente.Email)
		if err == nil && clienteExistente != nil {
			return fmt.Errorf("email %s já está em uso", cliente.Email)
		}
	}

	// Gera um UUID se não foi fornecido
	if cliente.ID == uuid.Nil {
		cliente.ID = uuid.New()
	}
	
	// Define a data de cadastro como o momento atual
	cliente.DataCadastro = time.Now()
	
	// Define o cliente como ativo por padrão
	cliente.Ativo = true
	
	return s.repo.Create(cliente)
}

// Update atualiza um cliente existente
func (s *clienteService) Update(cliente *domain.Cliente) error {
	return s.repo.Update(cliente)
}

// Delete realiza o soft delete de um cliente
func (s *clienteService) Delete(id uuid.UUID) error {
	return s.repo.SoftDelete(id)
}
