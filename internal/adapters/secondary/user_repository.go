package secondary

import (
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"

	"gorm.io/gorm"
)

// UserModel é o modelo de persistência para usuários
type UserModel struct {
	gorm.Model
	Name     string `gorm:"size:100;not null"`
	Email    string `gorm:"size:100;uniqueIndex;not null"`
	Password string `gorm:"size:255;not null"`
}

// toEntity converte um modelo de persistência para uma entidade de domínio
func (m *UserModel) toEntity() domain.User {
	return domain.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// fromEntity converte uma entidade de domínio para um modelo de persistência
func fromEntity(entity *domain.User) *UserModel {
	return &UserModel{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Name:     entity.Name,
		Email:    entity.Email,
		Password: entity.Password,
	}
}

// userRepository implementa a interface UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository cria uma nova instância de UserRepository
func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

// FindAll retorna todos os usuários
func (r *userRepository) FindAll() ([]domain.User, error) {
	var models []UserModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]domain.User, len(models))
	for i, model := range models {
		users[i] = model.toEntity()
	}

	return users, nil
}

// FindByID retorna um usuário pelo ID
func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var model UserModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	user := model.toEntity()
	return &user, nil
}

// Create cria um novo usuário
func (r *userRepository) Create(user *domain.User) error {
	model := fromEntity(user)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	*user = model.toEntity()
	return nil
}

// Update atualiza um usuário existente
func (r *userRepository) Update(user *domain.User) error {
	model := fromEntity(user)
	if err := r.db.Save(model).Error; err != nil {
		return err
	}

	*user = model.toEntity()
	return nil
}

// Delete remove um usuário pelo ID
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&UserModel{}, id).Error
}
