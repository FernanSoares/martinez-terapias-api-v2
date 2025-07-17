package services

import (
	"errors"
	"fmt"
	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Chave de assinatura do JWT - em produção, deve ser uma variável de ambiente
const jwtSecret = "seu_segredo_jwt_aqui"

// usuarioService implementa a interface UsuarioService
type usuarioService struct {
	repo            ports.UsuarioRepository
	agendamentoRepo ports.AgendamentoRepository
}

// NewUsuarioService cria uma nova instância de UsuarioService
func NewUsuarioService(repo ports.UsuarioRepository, agendamentoRepo ports.AgendamentoRepository) ports.UsuarioService {
	return &usuarioService{
		repo:            repo,
		agendamentoRepo: agendamentoRepo,
	}
}

// GetAll retorna todos os usuários
func (s *usuarioService) GetAll() ([]domain.Usuario, error) {
	return s.repo.FindAll()
}

// GetByID retorna um usuário pelo ID
func (s *usuarioService) GetByID(id uuid.UUID) (*domain.Usuario, error) {
	return s.repo.FindByID(id)
}

// Registrar cria um novo usuário no sistema
func (s *usuarioService) Registrar(usuario *domain.Usuario, password string) error {
	// Verifica se o email já está em uso
	existingUser, err := s.repo.FindByEmail(usuario.Email)
	if err == nil && existingUser != nil {
		return errors.New("email já está em uso")
	}

	// Gera um hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	// Configura o usuário
	if usuario.ID == uuid.Nil {
		usuario.ID = uuid.New()
	}

	usuario.SenhaHash = string(hashedPassword)
	usuario.Ativo = true
	usuario.DataCadastro = time.Now()

	// Por padrão, novos registros são de clientes, a menos que especificado
	if usuario.Perfil == "" {
		usuario.Perfil = domain.PerfilCliente
	}

	return s.repo.Create(usuario)
}

// Login autentica um usuário e retorna um token JWT
func (s *usuarioService) Login(email string, password string) (*domain.TokenResponse, error) {
	// Busca o usuário pelo email
	usuario, err := s.repo.FindByEmail(email)
	if err != nil || usuario == nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Verifica se o usuário está ativo
	if !usuario.Ativo {
		return nil, errors.New("usuário desativado")
	}

	// Compara a senha
	err = bcrypt.CompareHashAndPassword([]byte(usuario.SenhaHash), []byte(password))
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Cria o token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     usuario.ID.String(),
		"nome":   usuario.NomeCompleto,
		"email":  usuario.Email,
		"perfil": usuario.Perfil,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token válido por 24 horas
	})

	// Assina o token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	// Cria a resposta
	response := &domain.TokenResponse{}
	response.Token = tokenString
	response.Usuario.ID = usuario.ID
	response.Usuario.Nome = usuario.NomeCompleto
	response.Usuario.Perfil = usuario.Perfil

	return response, nil
}

// Update atualiza um usuário existente
func (s *usuarioService) Update(usuario *domain.Usuario) error {
	// Busca o usuário atual para manter campos sensíveis
	existingUser, err := s.repo.FindByID(usuario.ID)
	if err != nil {
		return err
	}

	// Mantém o hash da senha se não for alterado
	if usuario.SenhaHash == "" {
		usuario.SenhaHash = existingUser.SenhaHash
	}

	return s.repo.Update(usuario)
}

// Delete realiza o soft delete de um usuário
func (s *usuarioService) Delete(id uuid.UUID) error {
	return s.repo.SoftDelete(id)
}

// GetPerfilCliente retorna o perfil de qualquer usuário autenticado
func (s *usuarioService) GetPerfilCliente(id uuid.UUID) (*domain.Usuario, error) {
	usuario, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Remover campos sensíveis
	usuario.SenhaHash = ""

	return usuario, nil
}

// GetByEmail retorna um usuário pelo email
func (s *usuarioService) GetByEmail(email string) (*domain.Usuario, error) {
	return s.repo.FindByEmail(email)
}

// GetAgendamentosCliente retorna os agendamentos de um cliente
func (s *usuarioService) GetAgendamentosCliente(id uuid.UUID, usuarioSolicitanteID uuid.UUID) ([]domain.Agendamento, error) {
	// Verifica se o usuário existe
	usuario, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// Verifica se o usuário solicitante é admin ou massoterapeuta
	usuarioSolicitante, err := s.repo.FindByID(usuarioSolicitanteID)
	if err != nil {
		return nil, err
	}
	
	// Permite acesso se for o próprio cliente ou se for admin/massoterapeuta
	if usuario.Perfil != domain.PerfilCliente && 
	   id != usuarioSolicitanteID && 
	   usuarioSolicitante.Perfil != domain.PerfilAdmin && 
	   usuarioSolicitante.Perfil != domain.PerfilMassoterapeuta {
		return nil, errors.New("acesso negado")
	}

	// Busca os agendamentos
	return s.agendamentoRepo.FindByClienteID(id)
}
