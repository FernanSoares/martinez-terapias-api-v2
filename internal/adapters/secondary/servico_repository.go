package secondary

import (
	"fmt"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServicoModel é o modelo de persistência para serviços
type ServicoModel struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	NomeServico    string    `gorm:"size:100;not null"`
	Descricao      string    `gorm:"size:500"`
	DuracaoMinutos int       `gorm:"not null"`
	Valor          float64   `gorm:"type:decimal(10,2);not null"`
	ImagemURL      string    `gorm:"size:255"`
}

// TableName define o nome da tabela no banco de dados
func (ServicoModel) TableName() string {
	return "servicos"
}

// toEntity converte um modelo de persistência para uma entidade de domínio
func (m *ServicoModel) toEntity() domain.Servico {
	return domain.Servico{
		ID:             m.ID,
		NomeServico:    m.NomeServico,
		Descricao:      m.Descricao,
		DuracaoMinutos: m.DuracaoMinutos,
		Valor:          m.Valor,
		ImagemURL:      m.ImagemURL,
	}
}

// fromEntity converte uma entidade de domínio para um modelo de persistência
func servicoFromEntity(entity *domain.Servico) *ServicoModel {
	return &ServicoModel{
		ID:             entity.ID,
		NomeServico:    entity.NomeServico,
		Descricao:      entity.Descricao,
		DuracaoMinutos: entity.DuracaoMinutos,
		Valor:          entity.Valor,
		ImagemURL:      entity.ImagemURL,
	}
}

// servicoRepository implementa a interface ServicoRepository
type servicoRepository struct {
	db *gorm.DB
}

// NewServicoRepository cria uma nova instância de ServicoRepository
func NewServicoRepository(db *gorm.DB) ports.ServicoRepository {
	// Garante que a tabela foi criada ou atualizada
	db.AutoMigrate(&ServicoModel{})
	
	return &servicoRepository{db: db}
}

// FindAll retorna todos os serviços
func (r *servicoRepository) FindAll() ([]domain.Servico, error) {
	var models []ServicoModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	servicos := make([]domain.Servico, len(models))
	for i, model := range models {
		servicos[i] = model.toEntity()
	}

	return servicos, nil
}

// FindByID retorna um serviço pelo ID
func (r *servicoRepository) FindByID(id uuid.UUID) (*domain.Servico, error) {
	var model ServicoModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	servico := model.toEntity()
	return &servico, nil
}

// Create cria um novo serviço
func (r *servicoRepository) Create(servico *domain.Servico) error {
	model := servicoFromEntity(servico)
	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("erro ao criar serviço: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*servico = model.toEntity()
	return nil
}

// Update atualiza um serviço existente
func (r *servicoRepository) Update(servico *domain.Servico) error {
	model := servicoFromEntity(servico)
	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("erro ao atualizar serviço: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*servico = model.toEntity()
	return nil
}

// Delete remove um serviço
func (r *servicoRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&ServicoModel{}, id).Error; err != nil {
		return fmt.Errorf("erro ao excluir serviço: %w", err)
	}
	return nil
}
