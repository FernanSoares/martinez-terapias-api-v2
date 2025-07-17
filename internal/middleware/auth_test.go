package middleware

import (
	"martinezterapias/api/internal/core/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Variável para os testes
// já declarada em auth.go, então não redeclaramos aqui
// var usuarioService *MockUsuarioService

// MockUsuarioService para simular o serviço de usuário nos testes
type MockUsuarioService struct {
	mock.Mock
}

func (m *MockUsuarioService) Registrar(usuario *domain.Usuario, senha string) error {
	args := m.Called(usuario, senha)
	return args.Error(0)
}

func (m *MockUsuarioService) Login(email, senha string) (*domain.TokenResponse, error) {
	args := m.Called(email, senha)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenResponse), args.Error(1)
}

func (m *MockUsuarioService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}

func (m *MockUsuarioService) GetByID(id uuid.UUID) (*domain.Usuario, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioService) GetByEmail(email string) (*domain.Usuario, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Usuario), args.Error(1)
}

func (m *MockUsuarioService) GetAll() ([]domain.Usuario, error) {
	args := m.Called()
	return args.Get(0).([]domain.Usuario), args.Error(1)
}

func (m *MockUsuarioService) GetAgendamentos(usuarioID uuid.UUID) ([]domain.Agendamento, error) {
	args := m.Called(usuarioID)
	return args.Get(0).([]domain.Agendamento), args.Error(1)
}

// Configurar o ambiente de teste
func setupTestRouter() (*gin.Engine, *MockUsuarioService) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUsuarioService)
	
	// Substituir o serviço original pelo mock
	usuarioService = mockService
	
	router := gin.New()
	return router, mockService
}

// Teste para o middleware JWTAuthMiddleware
func TestJWTAuthMiddleware(t *testing.T) {
	// Configurar o router e o mock
	router, mockService := setupTestRouter()
	
	// Adicionar o middleware e uma rota protegida
	protected := router.Group("/")
	protected.Use(JWTAuthMiddleware())
	protected.GET("/protegido", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	
	// Teste com token válido
	t.Run("Token Válido", func(t *testing.T) {
		// Criar claims válidas para o mock retornar
		usuarioID := uuid.New()
		claims := jwt.MapClaims{
			"id":     usuarioID.String(),
			"email":  "teste@exemplo.com",
			"nome":   "Teste Usuário",
			"perfil": "cliente",
			"exp":    time.Now().Add(time.Hour).Unix(),
		}
		
		// Configurar o mock para retornar as claims quando ValidateToken for chamado
		// O middleware retira o prefixo "Bearer " antes de chamar ValidateToken
		mockService.On("ValidateToken", "valid-token").Return(claims, nil).Once()
		
		// Criar a requisição com o token no cabeçalho
		req, _ := http.NewRequest("GET", "/protegido", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		resp := httptest.NewRecorder()
		
		// Executar a requisição
		router.ServeHTTP(resp, req)
		
		// Verificar o resultado
		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
	
	// Teste sem token
	t.Run("Sem Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protegido", nil)
		resp := httptest.NewRecorder()
		
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
	
	// Teste com formato de token inválido
	t.Run("Formato de Token Inválido", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protegido", nil)
		req.Header.Set("Authorization", "invalid-format")
		resp := httptest.NewRecorder()
		
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
	
	// Teste com token inválido
	t.Run("Token Inválido", func(t *testing.T) {
		// O middleware retira o prefixo "Bearer " antes de chamar ValidateToken
		mockService.On("ValidateToken", "invalid-token").Return(nil, jwt.ErrSignatureInvalid).Once()
		
		req, _ := http.NewRequest("GET", "/protegido", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp := httptest.NewRecorder()
		
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		mockService.AssertExpectations(t)
	})
}

// Teste para o middleware RequirePerfil
func TestRequirePerfil(t *testing.T) {
	// Configurar o router
	router, _ := setupTestRouter()
	
	// Adicionar uma rota protegida com requisito de perfil
	protected := router.Group("/")
	protected.Use(func(c *gin.Context) {
		// Simular o middleware de autenticação que adiciona usuario_id e perfil no contexto
		c.Set(UsuarioIDKey, uuid.New())
		c.Set(PerfilKey, "admin")
		c.Next()
	})
	protected.Use(RequirePerfil("admin"))
	protected.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	
	// Rota que requer múltiplos perfis
	multiPerfil := router.Group("/")
	multiPerfil.Use(func(c *gin.Context) {
		c.Set(UsuarioIDKey, uuid.New())
		c.Set(PerfilKey, "massoterapeuta")
		c.Next()
	})
	multiPerfil.Use(RequirePerfil("admin", "massoterapeuta"))
	multiPerfil.GET("/admin-ou-masso", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	
	// Rota com perfil não autorizado
	naoAutorizado := router.Group("/")
	naoAutorizado.Use(func(c *gin.Context) {
		c.Set(UsuarioIDKey, uuid.New())
		c.Set(PerfilKey, "cliente")
		c.Next()
	})
	naoAutorizado.Use(RequirePerfil("admin"))
	naoAutorizado.GET("/apenas-admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
	
	// Teste com perfil correto
	t.Run("Perfil Correto", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin", nil)
		resp := httptest.NewRecorder()
		
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusOK, resp.Code)
	})
	
	// Teste com um dos perfis permitidos
	t.Run("Um dos Perfis Permitidos", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/admin-ou-masso", nil)
		resp := httptest.NewRecorder()
		
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusOK, resp.Code)
	})
	
	// Teste com perfil não autorizado
	t.Run("Perfil Não Autorizado", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/apenas-admin", nil)
		resp := httptest.NewRecorder()
		
		router.ServeHTTP(resp, req)
		
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})
}

// Teste para resetar a variável global após os testes
func TestMain(m *testing.M) {
	// Salvar a referência original
	original := usuarioService
	
	// Executar os testes
	m.Run()
	
	// Restaurar a referência original
	usuarioService = original
}
