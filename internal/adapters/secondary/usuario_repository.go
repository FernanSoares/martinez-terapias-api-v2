package secondary

import (
	"fmt"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UsuarioModel é o modelo de persistência para usuários do sistema
type UsuarioModel struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	NomeCompleto string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:100;uniqueIndex;not null"`
	SenhaHash    string    `gorm:"size:255;not null"`
	Telefone     string    `gorm:"size:20"`
	Perfil       string    `gorm:"size:20;not null"`
	Ativo        bool      `gorm:"default:true"`
	DataCadastro time.Time `gorm:"type:timestamp;not null"`
}

// TableName define o nome da tabela no banco de dados
func (UsuarioModel) TableName() string {
	return "usuarios"
}

// toEntity converte um modelo de persistência para uma entidade de domínio
func (m *UsuarioModel) toEntity() domain.Usuario {
	return domain.Usuario{
		ID:           m.ID,
		NomeCompleto: m.NomeCompleto,
		Email:        m.Email,
		SenhaHash:    m.SenhaHash,
		Telefone:     m.Telefone,
		Perfil:       domain.TipoPerfil(m.Perfil),
		Ativo:        m.Ativo,
		DataCadastro: m.DataCadastro,
	}
}

// fromEntity converte uma entidade de domínio para um modelo de persistência
func usuarioFromEntity(entity *domain.Usuario) *UsuarioModel {
	return &UsuarioModel{
		ID:           entity.ID,
		NomeCompleto: entity.NomeCompleto,
		Email:        entity.Email,
		SenhaHash:    entity.SenhaHash,
		Telefone:     entity.Telefone,
		Perfil:       string(entity.Perfil),
		Ativo:        entity.Ativo,
		DataCadastro: entity.DataCadastro,
	}
}

// usuarioRepository implementa a interface UsuarioRepository
type usuarioRepository struct {
	db *gorm.DB
}

// NewUsuarioRepository cria uma nova instância de UsuarioRepository
func NewUsuarioRepository(db *gorm.DB) ports.UsuarioRepository {
	// Garante que a tabela foi criada ou atualizada
	db.AutoMigrate(&UsuarioModel{})
	
	return &usuarioRepository{db: db}
}

// FindAll retorna todos os usuários
func (r *usuarioRepository) FindAll() ([]domain.Usuario, error) {
	var models []UsuarioModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	usuarios := make([]domain.Usuario, len(models))
	for i, model := range models {
		usuarios[i] = model.toEntity()
	}

	return usuarios, nil
}

// FindByID retorna um usuário pelo ID
func (r *usuarioRepository) FindByID(id uuid.UUID) (*domain.Usuario, error) {
	var model UsuarioModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	usuario := model.toEntity()
	return &usuario, nil
}

// FindByEmail retorna um usuário pelo email
func (r *usuarioRepository) FindByEmail(email string) (*domain.Usuario, error) {
	var model UsuarioModel
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}

	usuario := model.toEntity()
	return &usuario, nil
}

// Create cria um novo usuário
func (r *usuarioRepository) Create(usuario *domain.Usuario) error {
	model := usuarioFromEntity(usuario)
	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("erro ao criar usuário: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*usuario = model.toEntity()
	return nil
}

// Update atualiza um usuário existente
func (r *usuarioRepository) Update(usuario *domain.Usuario) error {
	model := usuarioFromEntity(usuario)
	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*usuario = model.toEntity()
	return nil
}

// SoftDelete realiza o soft delete de um usuário
func (r *usuarioRepository) SoftDelete(id uuid.UUID) error {
	// Atualiza apenas o campo ativo para false
	if err := r.db.Model(&UsuarioModel{}).Where("id = ?", id).Update("ativo", false).Error; err != nil {
		return fmt.Errorf("erro ao desativar usuário: %w", err)
	}
	return nil
}
