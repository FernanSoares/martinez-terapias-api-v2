package ports

import (
	"martinezterapias/api/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

// ClienteRepository define a interface para operações de persistência de clientes
type ClienteRepository interface {
	FindAll() ([]domain.Cliente, error)
	FindByID(id uuid.UUID) (*domain.Cliente, error)
	FindByNome(nome string) ([]domain.Cliente, error)
	FindByEmail(email string) (*domain.Cliente, error)
	FindByAtivo(ativo bool) ([]domain.Cliente, error)
	Create(cliente *domain.Cliente) error
	Update(cliente *domain.Cliente) error
	SoftDelete(id uuid.UUID) error
}

// ServicoRepository define a interface para operações de persistência de serviços
type ServicoRepository interface {
	FindAll() ([]domain.Servico, error)
	FindByID(id uuid.UUID) (*domain.Servico, error)
	Create(servico *domain.Servico) error
	Update(servico *domain.Servico) error
	Delete(id uuid.UUID) error
}

// AgendamentoRepository define a interface para operações de persistência de agendamentos
type AgendamentoRepository interface {
	FindAll() ([]domain.Agendamento, error)
	FindByID(id uuid.UUID) (*domain.Agendamento, error)
	FindByClienteID(clienteID uuid.UUID) ([]domain.Agendamento, error)
	FindByMassoterapeutaID(massoterapeutaID uuid.UUID) ([]domain.Agendamento, error)
	FindByPeriodo(dataInicio, dataFim time.Time) ([]domain.Agendamento, error)
	FindByStatus(status domain.StatusAgendamento) ([]domain.Agendamento, error)
	Create(agendamento *domain.Agendamento) error
	Update(agendamento *domain.Agendamento) error
	UpdateStatus(id uuid.UUID, status domain.StatusAgendamento) error
	Delete(id uuid.UUID) error
}

// AnamneseRepository define a interface para operações de persistência de fichas de anamnese
type AnamneseRepository interface {
	FindAll() ([]domain.Anamnese, error)
	FindByID(id uuid.UUID) (*domain.Anamnese, error)
	FindByClienteID(clienteID uuid.UUID) ([]domain.Anamnese, error)
	Create(anamnese *domain.Anamnese) error
	Update(anamnese *domain.Anamnese) error
}

// UsuarioRepository define a interface para operações de persistência de usuários do sistema
type UsuarioRepository interface {
	FindAll() ([]domain.Usuario, error)
	FindByID(id uuid.UUID) (*domain.Usuario, error)
	FindByEmail(email string) (*domain.Usuario, error)
	Create(usuario *domain.Usuario) error
	Update(usuario *domain.Usuario) error
	SoftDelete(id uuid.UUID) error
}

// UserRepository define a interface para operações de persistência de usuários (compatibilidade)
type UserRepository interface {
	FindAll() ([]domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id uint) error
}
