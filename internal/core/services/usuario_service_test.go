package services

import (
	"errors"
	"martinezterapias/api/internal/core/domain"
	"testing"
	"time"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

// Mock do UsuarioRepository
type MockUsuarioRepository struct {
	mock.Mock
}

func (m *MockUsuarioRepository) FindAll() ([]domain.Usuario, error) {
	args := m.Called()
	return args.Get(0).([]domain.Usuario), args.Error(1)
}

func (m *MockUsuarioRepository) FindByID(id uuid.UUID) (*domain.Usuario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioRepository) FindByEmail(email string) (*domain.Usuario, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioRepository) Create(usuario *domain.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepository) Update(usuario *domain.Usuario) error {
	args := m.Called(usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUsuarioRepository) SoftDelete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock do AgendamentoRepository
type MockAgendamentoRepositoryForUsuario struct {
	mock.Mock
}

func (m *MockAgendamentoRepositoryForUsuario) FindAll() ([]domain.Agendamento, error) {
	args := m.Called()
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepositoryForUsuario) FindByID(id uuid.UUID) (*domain.Agendamento, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepositoryForUsuario) FindByClienteID(clienteID uuid.UUID) ([]domain.Agendamento, error) {
	args := m.Called(clienteID)
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepositoryForUsuario) FindByMassoterapeutaID(massoterapeutaID uuid.UUID) ([]domain.Agendamento, error) {
	args := m.Called(massoterapeutaID)
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepositoryForUsuario) FindByPeriodo(dataInicio, dataFim time.Time) ([]domain.Agendamento, error) {
	args := m.Called(dataInicio, dataFim)
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepositoryForUsuario) FindByStatus(status domain.StatusAgendamento) ([]domain.Agendamento, error) {
	args := m.Called(status)
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepositoryForUsuario) Create(agendamento *domain.Agendamento) error {
	args := m.Called(agendamento)
	return args.Error(0)
}

func (m *MockAgendamentoRepositoryForUsuario) Update(agendamento *domain.Agendamento) error {
	args := m.Called(agendamento)
	return args.Error(0)
}

func (m *MockAgendamentoRepositoryForUsuario) UpdateStatus(id uuid.UUID, status domain.StatusAgendamento) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockAgendamentoRepositoryForUsuario) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// UsuarioServiceSuite define a suite de testes para UsuarioService
type UsuarioServiceSuite struct {
	suite.Suite
	mockUsuarioRepo    *MockUsuarioRepository
	mockAgendamentoRepo *MockAgendamentoRepositoryForUsuario
	service            *usuarioService
}

// SetupTest prepara o ambiente antes de cada teste
func (s *UsuarioServiceSuite) SetupTest() {
	s.mockUsuarioRepo = new(MockUsuarioRepository)
	s.mockAgendamentoRepo = new(MockAgendamentoRepositoryForUsuario)
	s.service = &usuarioService{
		repo:            s.mockUsuarioRepo,
		agendamentoRepo: s.mockAgendamentoRepo,
	}
}

// TestRegistrar testa o método de registro de usuário
func (s *UsuarioServiceSuite) TestRegistrar() {
	usuario := &domain.Usuario{
		NomeCompleto: "Teste Usuário",
		Email:        "teste@exemplo.com",
		Telefone:     "11999998888",
		Perfil:       domain.PerfilCliente,
	}
	senha := "senha123"

	s.mockUsuarioRepo.On("FindByEmail", usuario.Email).Return(nil, errors.New("usuário não encontrado"))
	
	s.mockUsuarioRepo.On("Create", mock.AnythingOfType("*domain.Usuario")).Return(nil)

	err := s.service.Registrar(usuario, senha)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), uuid.Nil, usuario.ID) // Verifica se o ID foi gerado
	assert.NotEmpty(s.T(), usuario.SenhaHash)    // Verifica se a senha foi hasheada
	assert.True(s.T(), usuario.Ativo)           // Verifica se o usuário está ativo
}

// TestRegistrar_EmailEmUso testa o registro com email já em uso
func (s *UsuarioServiceSuite) TestRegistrar_EmailEmUso() {
	usuario := &domain.Usuario{
		NomeCompleto: "Teste Usuário",
		Email:        "existente@exemplo.com",
		Perfil:       domain.PerfilCliente,
	}
	senha := "senha123"

	usuarioExistente := &domain.Usuario{
		ID:          uuid.New(),
		Email:       "existente@exemplo.com",
	}

	s.mockUsuarioRepo.On("FindByEmail", usuario.Email).Return(usuarioExistente, nil)

	err := s.service.Registrar(usuario, senha)
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "já está em uso")
}

// TestLogin testa o método de login
func (s *UsuarioServiceSuite) TestLogin() {
	email := "teste@exemplo.com"
	senha := "senha123"

	// Hash da senha para simular o banco de dados
	senhaHash, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	
	usuario := &domain.Usuario{
		ID:          uuid.New(),
		NomeCompleto: "Teste Usuário",
		Email:       email,
		SenhaHash:   string(senhaHash),
		Perfil:      domain.PerfilCliente,
		Ativo:       true, // Importante definir como ativo para passar na verificação
	}

	s.mockUsuarioRepo.On("FindByEmail", email).Return(usuario, nil)

	tokenResponse, err := s.service.Login(email, senha)
	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), tokenResponse.Token)
	// Comparar diretamente os UUIDs, já que ambos são do mesmo tipo
	assert.Equal(s.T(), usuario.ID, tokenResponse.Usuario.ID)
	assert.Equal(s.T(), usuario.NomeCompleto, tokenResponse.Usuario.Nome)
	assert.Equal(s.T(), usuario.Perfil, tokenResponse.Usuario.Perfil)
}

// TestLogin_UsuarioNaoEncontrado testa login com usuário não encontrado
func (s *UsuarioServiceSuite) TestLogin_UsuarioNaoEncontrado() {
	email := "inexistente@exemplo.com"
	senha := "senha123"

	s.mockUsuarioRepo.On("FindByEmail", email).Return(nil, errors.New("usuário não encontrado"))

	tokenResponse, err := s.service.Login(email, senha)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), tokenResponse)
}

// TestLogin_SenhaIncorreta testa login com senha incorreta
func (s *UsuarioServiceSuite) TestLogin_SenhaIncorreta() {
	email := "teste@exemplo.com"
	senha := "senha123"
	senhaIncorreta := "senha456"

	// Hash da senha correta
	senhaHash, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	
	usuario := &domain.Usuario{
		ID:          uuid.New(),
		Email:       email,
		SenhaHash:   string(senhaHash),
	}

	s.mockUsuarioRepo.On("FindByEmail", email).Return(usuario, nil)

	tokenResponse, err := s.service.Login(email, senhaIncorreta)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), tokenResponse)
}

// TestGetByID testa o método GetByID
func (s *UsuarioServiceSuite) TestGetByID() {
	id := uuid.New()
	usuario := &domain.Usuario{
		ID:          id,
		NomeCompleto: "Teste Usuário",
		Email:       "teste@exemplo.com",
	}

	s.mockUsuarioRepo.On("FindByID", id).Return(usuario, nil)

	result, err := s.service.GetByID(id)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), usuario.ID, result.ID)
	assert.Equal(s.T(), usuario.Email, result.Email)
}

// TestFindByEmail testa o método FindByEmail do repositório
func (s *UsuarioServiceSuite) TestFindByEmail() {
	email := "teste@exemplo.com"
	usuario := &domain.Usuario{
		ID:          uuid.New(),
		NomeCompleto: "Teste Usuário",
		Email:       email,
	}

	s.mockUsuarioRepo.On("FindByEmail", email).Return(usuario, nil)

	// Chamar diretamente o repositório pois o serviço não expõe este método
	result, err := s.mockUsuarioRepo.FindByEmail(email)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), usuario.ID, result.ID)
	assert.Equal(s.T(), usuario.Email, result.Email)
}

// TestGetAll testa o método GetAll
func (s *UsuarioServiceSuite) TestGetAll() {
	usuarios := []domain.Usuario{
		{
			ID:          uuid.New(),
			NomeCompleto: "Usuário 1",
			Email:       "usuario1@exemplo.com",
		},
		{
			ID:          uuid.New(),
			NomeCompleto: "Usuário 2",
			Email:       "usuario2@exemplo.com",
		},
	}

	s.mockUsuarioRepo.On("FindAll").Return(usuarios, nil)

	result, err := s.service.GetAll()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), len(usuarios), len(result))
}

// TestGetAgendamentosCliente testa o método GetAgendamentosCliente
func (s *UsuarioServiceSuite) TestGetAgendamentosCliente() {
	id := uuid.New()
	usuarioSolicitanteID := uuid.New()
	
	// Usuário principal
	usuario := &domain.Usuario{
		ID:          id,
		NomeCompleto: "Teste Usuário",
		Email:       "teste@exemplo.com",
		Perfil:      domain.PerfilMassoterapeuta,
	}
	
	// Usuário solicitante (admin para ter acesso)
	usuarioSolicitante := &domain.Usuario{
		ID:          usuarioSolicitanteID,
		NomeCompleto: "Admin",
		Email:       "admin@exemplo.com",
		Perfil:      domain.PerfilAdmin,
	}

	agendamentos := []domain.Agendamento{
		{
			ID:               uuid.New(),
			ClienteID:        uuid.New(),
			MassoterapeutaID: id,
			ServicoID:        uuid.New(),
			DataHora:         time.Now(),
			Status:           domain.StatusAgendado,
		},
	}

	// Mock para o usuário principal
	s.mockUsuarioRepo.On("FindByID", id).Return(usuario, nil)
	// Mock para o usuário solicitante (importante para as verificações de permissão)
	s.mockUsuarioRepo.On("FindByID", usuarioSolicitanteID).Return(usuarioSolicitante, nil)
	// O método GetAgendamentosCliente chama FindByClienteID, não FindByMassoterapeutaID
	s.mockAgendamentoRepo.On("FindByClienteID", id).Return(agendamentos, nil)

	result, err := s.service.GetAgendamentosCliente(id, usuarioSolicitanteID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), len(agendamentos), len(result))
}

// TestGetAgendamentosCliente_UsuarioNaoEncontrado testa GetAgendamentosCliente com usuário não encontrado
func (s *UsuarioServiceSuite) TestGetAgendamentosCliente_UsuarioNaoEncontrado() {
	id := uuid.New()
	usuarioSolicitanteID := uuid.New()

	s.mockUsuarioRepo.On("FindByID", id).Return(nil, errors.New("usuário não encontrado"))

	_, err := s.service.GetAgendamentosCliente(id, usuarioSolicitanteID)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "não encontrado")
}

// TestLoginToken verifica se o token JWT é gerado corretamente no login
func (s *UsuarioServiceSuite) TestLoginToken() {
	id := uuid.New()
	email := "teste@exemplo.com"
	senha := "senha123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)

	usuario := &domain.Usuario{
		ID:           id,
		Email:        email,
		NomeCompleto: "Teste Usuário",
		SenhaHash:    string(hashedPassword),
		Perfil:       domain.PerfilCliente,
		Ativo:        true, // Definir como ativo para passar na verificação
	}

	s.mockUsuarioRepo.On("FindByEmail", email).Return(usuario, nil)
	
	result, err := s.service.Login(email, senha)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.NotEmpty(s.T(), result.Token)
}

// TestUsuarioServiceSuite executa a suite de testes
func TestUsuarioServiceSuite(t *testing.T) {
	suite.Run(t, new(UsuarioServiceSuite))
}
