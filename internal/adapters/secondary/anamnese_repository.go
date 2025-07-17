package secondary

import (
	"fmt"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AnamneseModel é o modelo de persistência para fichas de anamnese
type AnamneseModel struct {
	gorm.Model
	ID                      uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClienteID               uuid.UUID `gorm:"type:uuid;not null;index"`
	DataPreenchimento       time.Time `gorm:"type:timestamp;not null"`
	QueixaPrincipal         string    `gorm:"size:500"`
	HistoricoDoencas        string    `gorm:"size:500"`
	CirurgiasPrevias        string    `gorm:"size:500"`
	MedicamentosEmUso       string    `gorm:"size:500"`
	Alergias                string    `gorm:"size:500"`
	HabitosDiarios          string    `gorm:"size:500"`
	ObservacoesMassoterapeuta string    `gorm:"size:1000"`
}

// TableName define o nome da tabela no banco de dados
func (AnamneseModel) TableName() string {
	return "anamneses"
}

// toEntity converte um modelo de persistência para uma entidade de domínio
func (m *AnamneseModel) toEntity() domain.Anamnese {
	return domain.Anamnese{
		ID:                      m.ID,
		ClienteID:               m.ClienteID,
		DataPreenchimento:       m.DataPreenchimento,
		QueixaPrincipal:         m.QueixaPrincipal,
		HistoricoDoencas:        m.HistoricoDoencas,
		CirurgiasPrevias:        m.CirurgiasPrevias,
		MedicamentosEmUso:       m.MedicamentosEmUso,
		Alergias:                m.Alergias,
		HabitosDiarios:          m.HabitosDiarios,
		ObservacoesMassoterapeuta: m.ObservacoesMassoterapeuta,
	}
}

// fromEntity converte uma entidade de domínio para um modelo de persistência
func anamneseFromEntity(entity *domain.Anamnese) *AnamneseModel {
	return &AnamneseModel{
		ID:                      entity.ID,
		ClienteID:               entity.ClienteID,
		DataPreenchimento:       entity.DataPreenchimento,
		QueixaPrincipal:         entity.QueixaPrincipal,
		HistoricoDoencas:        entity.HistoricoDoencas,
		CirurgiasPrevias:        entity.CirurgiasPrevias,
		MedicamentosEmUso:       entity.MedicamentosEmUso,
		Alergias:                entity.Alergias,
		HabitosDiarios:          entity.HabitosDiarios,
		ObservacoesMassoterapeuta: entity.ObservacoesMassoterapeuta,
	}
}

// anamneseRepository implementa a interface AnamneseRepository
type anamneseRepository struct {
	db *gorm.DB
}

// NewAnamneseRepository cria uma nova instância de AnamneseRepository
func NewAnamneseRepository(db *gorm.DB) ports.AnamneseRepository {
	// Garante que a tabela foi criada ou atualizada
	db.AutoMigrate(&AnamneseModel{})
	
	return &anamneseRepository{db: db}
}

// FindAll retorna todas as fichas de anamnese
func (r *anamneseRepository) FindAll() ([]domain.Anamnese, error) {
	var models []AnamneseModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	anamneses := make([]domain.Anamnese, len(models))
	for i, model := range models {
		anamneses[i] = model.toEntity()
	}

	return anamneses, nil
}

// FindByID retorna uma ficha de anamnese pelo ID
func (r *anamneseRepository) FindByID(id uuid.UUID) (*domain.Anamnese, error) {
	var model AnamneseModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	anamnese := model.toEntity()
	return &anamnese, nil
}

// FindByClienteID retorna fichas de anamnese de um cliente
func (r *anamneseRepository) FindByClienteID(clienteID uuid.UUID) ([]domain.Anamnese, error) {
	var models []AnamneseModel
	if err := r.db.Where("cliente_id = ?", clienteID).Order("data_preenchimento desc").Find(&models).Error; err != nil {
		return nil, err
	}

	anamneses := make([]domain.Anamnese, len(models))
	for i, model := range models {
		anamneses[i] = model.toEntity()
	}

	return anamneses, nil
}

// Create cria uma nova ficha de anamnese
func (r *anamneseRepository) Create(anamnese *domain.Anamnese) error {
	model := anamneseFromEntity(anamnese)
	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("erro ao criar ficha de anamnese: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*anamnese = model.toEntity()
	return nil
}

// Update atualiza uma ficha de anamnese existente
func (r *anamneseRepository) Update(anamnese *domain.Anamnese) error {
	model := anamneseFromEntity(anamnese)
	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("erro ao atualizar ficha de anamnese: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*anamnese = model.toEntity()
	return nil
}
