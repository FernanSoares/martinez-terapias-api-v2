package ports

import (
	"martinezterapias/api/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

// ClienteService define a interface para a lógica de negócio relacionada a clientes
type ClienteService interface {
	GetAll() ([]domain.Cliente, error)
	GetByID(id uuid.UUID) (*domain.Cliente, error)
	GetByNome(nome string) ([]domain.Cliente, error)
	GetByAtivo(ativo bool) ([]domain.Cliente, error)
	Create(cliente *domain.Cliente) error
	Update(cliente *domain.Cliente) error
	Delete(id uuid.UUID) error
}

// ServicoService define a interface para a lógica de negócio relacionada a serviços
type ServicoService interface {
	GetAll() ([]domain.Servico, error)
	GetByID(id uuid.UUID) (*domain.Servico, error)
	Create(servico *domain.Servico) error
	Update(servico *domain.Servico) error
	Delete(id uuid.UUID) error
}

// AgendamentoService define a interface para a lógica de negócio relacionada a agendamentos
type AgendamentoService interface {
	GetAll() ([]domain.Agendamento, error)
	GetByID(id uuid.UUID) (*domain.Agendamento, error)
	GetByClienteID(clienteID uuid.UUID) ([]domain.Agendamento, error)
	GetByMassoterapeutaID(massoterapeutaID uuid.UUID) ([]domain.Agendamento, error)
	GetByPeriodo(dataInicio, dataFim time.Time) ([]domain.Agendamento, error)
	GetByStatus(status domain.StatusAgendamento) ([]domain.Agendamento, error)
	Create(agendamento *domain.Agendamento) error
	Update(agendamento *domain.Agendamento) error
	UpdateStatus(id uuid.UUID, status domain.StatusAgendamento) error
	Delete(id uuid.UUID) error
	SolicitarReagendamento(id uuid.UUID, clienteID uuid.UUID) error
}

// AnamneseService define a interface para a lógica de negócio relacionada a fichas de anamnese
type AnamneseService interface {
	GetAll() ([]domain.Anamnese, error)
	GetByID(id uuid.UUID) (*domain.Anamnese, error)
	GetByClienteID(clienteID uuid.UUID) ([]domain.Anamnese, error)
	Create(anamnese *domain.Anamnese) error
	Update(anamnese *domain.Anamnese) error
}

// UsuarioService define a interface para a lógica de negócio relacionada a usuários do sistema
type UsuarioService interface {
	GetAll() ([]domain.Usuario, error)
	GetByID(id uuid.UUID) (*domain.Usuario, error)
	Registrar(usuario *domain.Usuario, password string) error
	Login(email string, password string) (*domain.TokenResponse, error)
	Update(usuario *domain.Usuario) error
	Delete(id uuid.UUID) error
	GetPerfilCliente(id uuid.UUID) (*domain.Usuario, error)
	GetAgendamentosCliente(id uuid.UUID, usuarioSolicitanteID uuid.UUID) ([]domain.Agendamento, error)
}

// WhatsAppService foi removido

// UserService define a interface para a lógica de negócio relacionada a usuários (compatibilidade)
type UserService interface {
	GetAll() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id uint) error
}
