package services

import (
	"fmt"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

// agendamentoService implementa a interface AgendamentoService
type agendamentoService struct {
	repo        ports.AgendamentoRepository
	clienteRepo ports.ClienteRepository
	servicoRepo ports.ServicoRepository
	usuarioRepo ports.UsuarioRepository
}

// NewAgendamentoService cria uma nova instância de AgendamentoService
func NewAgendamentoService(
	repo ports.AgendamentoRepository,
	clienteRepo ports.ClienteRepository,
	servicoRepo ports.ServicoRepository,
	usuarioRepo ports.UsuarioRepository,
) ports.AgendamentoService {
	return &agendamentoService{
		repo:        repo,
		clienteRepo: clienteRepo,
		servicoRepo: servicoRepo,
		usuarioRepo: usuarioRepo,
	}
}

// GetAll retorna todos os agendamentos
func (s *agendamentoService) GetAll() ([]domain.Agendamento, error) {
	return s.repo.FindAll()
}

// GetByID retorna um agendamento pelo ID
func (s *agendamentoService) GetByID(id uuid.UUID) (*domain.Agendamento, error) {
	return s.repo.FindByID(id)
}

// GetByClienteID retorna todos os agendamentos de um cliente
func (s *agendamentoService) GetByClienteID(clienteID uuid.UUID) ([]domain.Agendamento, error) {
	return s.repo.FindByClienteID(clienteID)
}

// GetByMassoterapeutaID retorna todos os agendamentos de um massoterapeuta
func (s *agendamentoService) GetByMassoterapeutaID(massoterapeutaID uuid.UUID) ([]domain.Agendamento, error) {
	return s.repo.FindByMassoterapeutaID(massoterapeutaID)
}

// GetByPeriodo retorna agendamentos em um período de tempo
func (s *agendamentoService) GetByPeriodo(dataInicio, dataFim time.Time) ([]domain.Agendamento, error) {
	return s.repo.FindByPeriodo(dataInicio, dataFim)
}

// GetByStatus retorna agendamentos por status
func (s *agendamentoService) GetByStatus(status domain.StatusAgendamento) ([]domain.Agendamento, error) {
	return s.repo.FindByStatus(status)
}

// Create cria um novo agendamento
func (s *agendamentoService) Create(agendamento *domain.Agendamento) error {
	// Verifica se o cliente existe
	cliente, err := s.clienteRepo.FindByID(agendamento.ClienteID)
	if err != nil || cliente == nil {
		return fmt.Errorf("cliente não encontrado: %w", err)
	}

	// Verifica se o serviço existe
	servico, err := s.servicoRepo.FindByID(agendamento.ServicoID)
	if err != nil || servico == nil {
		return fmt.Errorf("serviço não encontrado: %w", err)
	}

	// Verifica se o massoterapeuta existe
	massoterapeuta, err := s.usuarioRepo.FindByID(agendamento.MassoterapeutaID)
	if err != nil || massoterapeuta == nil || massoterapeuta.Perfil != domain.PerfilMassoterapeuta {
		return fmt.Errorf("massoterapeuta não encontrado: %w", err)
	}

	// Verifica disponibilidade de horário
	agendamentosNoPeriodo, err := s.repo.FindByPeriodo(
		agendamento.DataHora,
		agendamento.DataHora.Add(time.Duration(servico.DuracaoMinutos)*time.Minute),
	)
	if err != nil {
		return err
	}

	// Verifica se há conflito de horário para o massoterapeuta
	for _, a := range agendamentosNoPeriodo {
		if a.MassoterapeutaID == agendamento.MassoterapeutaID && a.Status != domain.StatusCancelado {
			return fmt.Errorf("horário já ocupado para este massoterapeuta")
		}
	}

	// Preenche o valor cobrado a partir do serviço se não for especificado
	if agendamento.ValorCobrado == 0 {
		agendamento.ValorCobrado = servico.Valor
	}

	// Gera um UUID se não foi fornecido
	if agendamento.ID == uuid.Nil {
		agendamento.ID = uuid.New()
	}

	// Define o status como "agendado" por padrão
	if agendamento.Status == "" {
		agendamento.Status = domain.StatusAgendado
	}

	return s.repo.Create(agendamento)
}

// Update atualiza um agendamento existente
func (s *agendamentoService) Update(agendamento *domain.Agendamento) error {
	return s.repo.Update(agendamento)
}

// UpdateStatus atualiza apenas o status de um agendamento
func (s *agendamentoService) UpdateStatus(id uuid.UUID, status domain.StatusAgendamento) error {
	return s.repo.UpdateStatus(id, status)
}

// Delete remove um agendamento
func (s *agendamentoService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// SolicitarReagendamento registra uma solicitação de reagendamento
func (s *agendamentoService) SolicitarReagendamento(id uuid.UUID, clienteID uuid.UUID) error {
	// Obtém o agendamento
	agendamento, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("agendamento não encontrado: %w", err)
	}

	// Verifica se o agendamento pertence ao cliente solicitante
	if agendamento.ClienteID != clienteID {
		return fmt.Errorf("não autorizado a reagendar este agendamento")
	}

	// Verifica se o cliente existe
	_, err = s.clienteRepo.FindByID(clienteID)
	if err != nil {
		return fmt.Errorf("cliente não encontrado: %w", err)
	}

	// Verifica se o serviço existe
	_, err = s.servicoRepo.FindByID(agendamento.ServicoID)
	if err != nil {
		return fmt.Errorf("serviço não encontrado: %w", err)
	}

	// Apenas registra a solicitação de reagendamento sem envio de mensagem
	// Essa funcionalidade de envio de mensagem via WhatsApp foi removida
	agendamento.Status = "reagendamento_solicitado"
	err = s.repo.Update(agendamento)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status do agendamento: %w", err)
	}

	return nil
}
