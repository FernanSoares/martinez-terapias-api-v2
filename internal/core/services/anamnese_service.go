package services

import (
	"fmt"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

// anamneseService implementa a interface AnamneseService
type anamneseService struct {
	repo        ports.AnamneseRepository
	clienteRepo ports.ClienteRepository
}

// NewAnamneseService cria uma nova instância de AnamneseService
func NewAnamneseService(repo ports.AnamneseRepository, clienteRepo ports.ClienteRepository) ports.AnamneseService {
	return &anamneseService{
		repo:        repo,
		clienteRepo: clienteRepo,
	}
}

// GetAll retorna todas as fichas de anamnese
func (s *anamneseService) GetAll() ([]domain.Anamnese, error) {
	return s.repo.FindAll()
}

// GetByID retorna uma ficha de anamnese pelo ID
func (s *anamneseService) GetByID(id uuid.UUID) (*domain.Anamnese, error) {
	return s.repo.FindByID(id)
}

// GetByClienteID retorna todas as fichas de anamnese de um cliente
func (s *anamneseService) GetByClienteID(clienteID uuid.UUID) ([]domain.Anamnese, error) {
	return s.repo.FindByClienteID(clienteID)
}

// Create cria uma nova ficha de anamnese
func (s *anamneseService) Create(anamnese *domain.Anamnese) error {
	// Verifica se o cliente existe
	cliente, err := s.clienteRepo.FindByID(anamnese.ClienteID)
	if err != nil || cliente == nil {
		return fmt.Errorf("cliente não encontrado: %w", err)
	}

	// Gera um UUID se não foi fornecido
	if anamnese.ID == uuid.Nil {
		anamnese.ID = uuid.New()
	}

	// Define a data de preenchimento como o momento atual
	if anamnese.DataPreenchimento.IsZero() {
		anamnese.DataPreenchimento = time.Now()
	}

	return s.repo.Create(anamnese)
}

// Update atualiza uma ficha de anamnese existente
func (s *anamneseService) Update(anamnese *domain.Anamnese) error {
	return s.repo.Update(anamnese)
}
