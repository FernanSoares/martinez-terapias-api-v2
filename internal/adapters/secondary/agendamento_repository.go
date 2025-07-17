package secondary

import (
	"fmt"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgendamentoModel é o modelo de persistência para agendamentos
type AgendamentoModel struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	DataHora         time.Time `gorm:"type:timestamp;not null;index"`
	ClienteID        uuid.UUID `gorm:"type:uuid;not null;index"`
	ServicoID        uuid.UUID `gorm:"type:uuid;not null"`
	MassoterapeutaID uuid.UUID `gorm:"type:uuid;not null;index"`
	ValorCobrado     float64   `gorm:"type:decimal(10,2);not null"`
	Status           string    `gorm:"size:20;not null;index"`
	Observacoes      string    `gorm:"size:500"`
}

// TableName define o nome da tabela no banco de dados
func (AgendamentoModel) TableName() string {
	return "agendamentos"
}

// toEntity converte um modelo de persistência para uma entidade de domínio
func (m *AgendamentoModel) toEntity() domain.Agendamento {
	return domain.Agendamento{
		ID:               m.ID,
		DataHora:         m.DataHora,
		ClienteID:        m.ClienteID,
		ServicoID:        m.ServicoID,
		MassoterapeutaID: m.MassoterapeutaID,
		ValorCobrado:     m.ValorCobrado,
		Status:           domain.StatusAgendamento(m.Status),
		Observacoes:      m.Observacoes,
	}
}

// fromEntity converte uma entidade de domínio para um modelo de persistência
func agendamentoFromEntity(entity *domain.Agendamento) *AgendamentoModel {
	return &AgendamentoModel{
		ID:               entity.ID,
		DataHora:         entity.DataHora,
		ClienteID:        entity.ClienteID,
		ServicoID:        entity.ServicoID,
		MassoterapeutaID: entity.MassoterapeutaID,
		ValorCobrado:     entity.ValorCobrado,
		Status:           string(entity.Status),
		Observacoes:      entity.Observacoes,
	}
}

// agendamentoRepository implementa a interface AgendamentoRepository
type agendamentoRepository struct {
	db *gorm.DB
}

// NewAgendamentoRepository cria uma nova instância de AgendamentoRepository
func NewAgendamentoRepository(db *gorm.DB) ports.AgendamentoRepository {
	// Garante que a tabela foi criada ou atualizada
	db.AutoMigrate(&AgendamentoModel{})

	return &agendamentoRepository{db: db}
}

// FindAll retorna todos os agendamentos
func (r *agendamentoRepository) FindAll() ([]domain.Agendamento, error) {
	var models []AgendamentoModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	agendamentos := make([]domain.Agendamento, len(models))
	for i, model := range models {
		agendamentos[i] = model.toEntity()
	}

	return agendamentos, nil
}

// FindByID retorna um agendamento pelo ID
func (r *agendamentoRepository) FindByID(id uuid.UUID) (*domain.Agendamento, error) {
	var model AgendamentoModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	agendamento := model.toEntity()
	return &agendamento, nil
}

// FindByClienteID retorna agendamentos de um cliente
func (r *agendamentoRepository) FindByClienteID(clienteID uuid.UUID) ([]domain.Agendamento, error) {
	var models []AgendamentoModel
	if err := r.db.Where("cliente_id = ?", clienteID).Find(&models).Error; err != nil {
		return nil, err
	}

	agendamentos := make([]domain.Agendamento, len(models))
	for i, model := range models {
		agendamentos[i] = model.toEntity()
	}

	return agendamentos, nil
}

// FindByMassoterapeutaID retorna agendamentos de um massoterapeuta
func (r *agendamentoRepository) FindByMassoterapeutaID(massoterapeutaID uuid.UUID) ([]domain.Agendamento, error) {
	var models []AgendamentoModel
	if err := r.db.Where("massoterapeuta_id = ?", massoterapeutaID).Find(&models).Error; err != nil {
		return nil, err
	}

	agendamentos := make([]domain.Agendamento, len(models))
	for i, model := range models {
		agendamentos[i] = model.toEntity()
	}

	return agendamentos, nil
}

// FindByPeriodo retorna agendamentos em um período de tempo
func (r *agendamentoRepository) FindByPeriodo(dataInicio, dataFim time.Time) ([]domain.Agendamento, error) {
	var models []AgendamentoModel
	if err := r.db.Where("data_hora BETWEEN ? AND ?", dataInicio, dataFim).Find(&models).Error; err != nil {
		return nil, err
	}

	agendamentos := make([]domain.Agendamento, len(models))
	for i, model := range models {
		agendamentos[i] = model.toEntity()
	}

	return agendamentos, nil
}

// FindByStatus retorna agendamentos por status
func (r *agendamentoRepository) FindByStatus(status domain.StatusAgendamento) ([]domain.Agendamento, error) {
	var models []AgendamentoModel
	if err := r.db.Where("status = ?", status).Find(&models).Error; err != nil {
		return nil, err
	}

	agendamentos := make([]domain.Agendamento, len(models))
	for i, model := range models {
		agendamentos[i] = model.toEntity()
	}

	return agendamentos, nil
}

// Create cria um novo agendamento
func (r *agendamentoRepository) Create(agendamento *domain.Agendamento) error {
	model := agendamentoFromEntity(agendamento)
	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("erro ao criar agendamento: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*agendamento = model.toEntity()
	return nil
}

// Update atualiza um agendamento existente
func (r *agendamentoRepository) Update(agendamento *domain.Agendamento) error {
	model := agendamentoFromEntity(agendamento)
	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("erro ao atualizar agendamento: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*agendamento = model.toEntity()
	return nil
}

// UpdateStatus atualiza apenas o status de um agendamento
func (r *agendamentoRepository) UpdateStatus(id uuid.UUID, status domain.StatusAgendamento) error {
	if err := r.db.Model(&AgendamentoModel{}).Where("id = ?", id).Update("status", string(status)).Error; err != nil {
		return fmt.Errorf("erro ao atualizar status do agendamento: %w", err)
	}
	return nil
}

// Delete remove um agendamento
func (r *agendamentoRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&AgendamentoModel{}, id).Error; err != nil {
		return fmt.Errorf("erro ao excluir agendamento: %w", err)
	}
	return nil
}
