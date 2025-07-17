package services

import (
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
)

// userService implementa a interface UserService
type userService struct {
	repo ports.UserRepository
}

// NewUserService cria uma nova instância de UserService
func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{
		repo: repo,
	}
}

// GetAll retorna todos os usuários
func (s *userService) GetAll() ([]domain.User, error) {
	return s.repo.FindAll()
}

// GetByID retorna um usuário pelo ID
func (s *userService) GetByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

// Create cria um novo usuário
func (s *userService) Create(user *domain.User) error {
	// Aqui poderíamos adicionar validações de negócio
	return s.repo.Create(user)
}

// Update atualiza um usuário existente
func (s *userService) Update(user *domain.User) error {
	// Aqui poderíamos adicionar validações de negócio
	return s.repo.Update(user)
}

// Delete remove um usuário pelo ID
func (s *userService) Delete(id uint) error {
	return s.repo.Delete(id)
}
