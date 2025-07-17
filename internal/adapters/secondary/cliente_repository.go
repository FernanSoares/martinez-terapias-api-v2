package secondary

import (
	"fmt"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ClienteModel é o modelo de persistência para clientes
type ClienteModel struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	NomeCompleto  string    `gorm:"size:100;not null"`
	CPF           string    `gorm:"size:14;uniqueIndex"`
	Email         string    `gorm:"size:100;uniqueIndex;not null"`
	Telefone      string    `gorm:"size:20;not null"`
	DataNascimento string    `gorm:"type:date"`
	DataCadastro  string    `gorm:"type:timestamp;not null"`
	Ativo         bool      `gorm:"default:true"`
}

// TableName define o nome da tabela no banco de dados
func (ClienteModel) TableName() string {
	return "clientes"
}

// toEntity converte um modelo de persistência para uma entidade de domínio
func (m *ClienteModel) toEntity() domain.Cliente {
	var dataNascimento, dataCadastro domain.TimeOrZero
	dataNascimento.UnmarshalJSON([]byte(`"` + m.DataNascimento + `"`))
	dataCadastro.UnmarshalJSON([]byte(`"` + m.DataCadastro + `"`))

	return domain.Cliente{
		ID:            m.ID,
		NomeCompleto:  m.NomeCompleto,
		CPF:           m.CPF,
		Email:         m.Email,
		Telefone:      m.Telefone,
		DataNascimento: dataNascimento.Time,
		DataCadastro:  dataCadastro.Time,
		Ativo:         m.Ativo,
	}
}

// fromEntity converte uma entidade de domínio para um modelo de persistência
func clienteFromEntity(entity *domain.Cliente) *ClienteModel {
	return &ClienteModel{
		ID:            entity.ID,
		NomeCompleto:  entity.NomeCompleto,
		CPF:           entity.CPF,
		Email:         entity.Email,
		Telefone:      entity.Telefone,
		DataNascimento: entity.DataNascimento.Format("2006-01-02"),
		DataCadastro:  entity.DataCadastro.Format("2006-01-02 15:04:05"),
		Ativo:         entity.Ativo,
	}
}

// clienteRepository implementa a interface ClienteRepository
type clienteRepository struct {
	db *gorm.DB
}

// NewClienteRepository cria uma nova instância de ClienteRepository
func NewClienteRepository(db *gorm.DB) ports.ClienteRepository {
	// Garante que a tabela foi criada ou atualizada
	db.AutoMigrate(&ClienteModel{})
	
	return &clienteRepository{db: db}
}

// FindAll retorna todos os clientes
func (r *clienteRepository) FindAll() ([]domain.Cliente, error) {
	var models []ClienteModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	clientes := make([]domain.Cliente, len(models))
	for i, model := range models {
		clientes[i] = model.toEntity()
	}

	return clientes, nil
}

// FindByID retorna um cliente pelo ID
func (r *clienteRepository) FindByID(id uuid.UUID) (*domain.Cliente, error) {
	var model ClienteModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	cliente := model.toEntity()
	return &cliente, nil
}

// FindByNome retorna clientes pelo nome (busca parcial)
func (r *clienteRepository) FindByNome(nome string) ([]domain.Cliente, error) {
	var models []ClienteModel
	if err := r.db.Where("nome_completo LIKE ?", "%"+nome+"%").Find(&models).Error; err != nil {
		return nil, err
	}

	clientes := make([]domain.Cliente, len(models))
	for i, model := range models {
		clientes[i] = model.toEntity()
	}

	return clientes, nil
}

// FindByEmail retorna um cliente pelo email
func (r *clienteRepository) FindByEmail(email string) (*domain.Cliente, error) {
	var model ClienteModel
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}

	cliente := model.toEntity()
	return &cliente, nil
}

// FindByAtivo retorna clientes pelo status ativo/inativo
func (r *clienteRepository) FindByAtivo(ativo bool) ([]domain.Cliente, error) {
	var models []ClienteModel
	if err := r.db.Where("ativo = ?", ativo).Find(&models).Error; err != nil {
		return nil, err
	}

	clientes := make([]domain.Cliente, len(models))
	for i, model := range models {
		clientes[i] = model.toEntity()
	}

	return clientes, nil
}

// Create cria um novo cliente
func (r *clienteRepository) Create(cliente *domain.Cliente) error {
	model := clienteFromEntity(cliente)
	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("erro ao criar cliente: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*cliente = model.toEntity()
	return nil
}

// Update atualiza um cliente existente
func (r *clienteRepository) Update(cliente *domain.Cliente) error {
	model := clienteFromEntity(cliente)
	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("erro ao atualizar cliente: %w", err)
	}

	// Atualizar a entidade com os dados do modelo
	*cliente = model.toEntity()
	return nil
}

// SoftDelete realiza o soft delete de um cliente
func (r *clienteRepository) SoftDelete(id uuid.UUID) error {
	// Atualiza apenas o campo ativo para false
	if err := r.db.Model(&ClienteModel{}).Where("id = ?", id).Update("ativo", false).Error; err != nil {
		return fmt.Errorf("erro ao desativar cliente: %w", err)
	}
	return nil
}
