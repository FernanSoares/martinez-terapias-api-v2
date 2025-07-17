package services

import (
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"

	"github.com/google/uuid"
)

// servicoService implementa a interface ServicoService
type servicoService struct {
	repo ports.ServicoRepository
}

// NewServicoService cria uma nova instância de ServicoService
func NewServicoService(repo ports.ServicoRepository) ports.ServicoService {
	return &servicoService{
		repo: repo,
	}
}

// GetAll retorna todos os serviços
func (s *servicoService) GetAll() ([]domain.Servico, error) {
	return s.repo.FindAll()
}

// GetByID retorna um serviço pelo ID
func (s *servicoService) GetByID(id uuid.UUID) (*domain.Servico, error) {
	return s.repo.FindByID(id)
}

// Create cria um novo serviço
func (s *servicoService) Create(servico *domain.Servico) error {
	// Gera um UUID se não foi fornecido
	if servico.ID == uuid.Nil {
		servico.ID = uuid.New()
	}
	
	return s.repo.Create(servico)
}

// Update atualiza um serviço existente
func (s *servicoService) Update(servico *domain.Servico) error {
	return s.repo.Update(servico)
}

// Delete remove um serviço
func (s *servicoService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
