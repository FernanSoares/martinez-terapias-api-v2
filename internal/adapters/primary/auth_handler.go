package primary

import (
	"net/http"

	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

// AuthHandler implementa os handlers HTTP para autenticação
type AuthHandler struct {
	usuarioService ports.UsuarioService
}

// NewAuthHandler cria uma nova instância de AuthHandler
func NewAuthHandler(usuarioService ports.UsuarioService) *AuthHandler {
	return &AuthHandler{
		usuarioService: usuarioService,
	}
}

// RegistroRequest é a estrutura para requisições de registro de novo usuário
type RegistroRequest struct {
	NomeCompleto string `json:"nome_completo" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Senha        string `json:"senha" binding:"required,min=6"`
	Telefone     string `json:"telefone" binding:"required"`
}

// LoginRequest é a estrutura para requisições de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Senha    string `json:"senha" binding:"required"`
}

// Registrar cadastra um novo usuário (cliente)
func (h *AuthHandler) Registrar(c *gin.Context) {
	var request RegistroRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Criar um novo usuário com perfil de cliente
	usuario := &domain.Usuario{
		NomeCompleto: request.NomeCompleto,
		Email:        request.Email,
		Telefone:     request.Telefone,
		Perfil:       domain.PerfilCliente,
	}

	if err := h.usuarioService.Registrar(usuario, request.Senha); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remover dados sensíveis antes de retornar
	usuario.SenhaHash = ""

	c.JSON(http.StatusCreated, usuario)
}

// RegistrarAdmin cadastra um novo usuário administrador
func (h *AuthHandler) RegistrarAdmin(c *gin.Context) {
	var request RegistroRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar se já existe algum administrador
	usuarios, err := h.usuarioService.GetAll()
	if err == nil && len(usuarios) > 0 {
		// Verificar se já existe um admin
		adminExists := false
		for _, u := range usuarios {
			if u.Perfil == domain.PerfilAdmin {
				adminExists = true
				break
			}
		}
		
		// Já existe um admin, bloquear criação para segurança
		if adminExists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Já existe um administrador no sistema"})
			return
		}
	}

	// Criar um novo usuário com perfil de administrador
	usuario := &domain.Usuario{
		NomeCompleto: request.NomeCompleto,
		Email:        request.Email,
		Telefone:     request.Telefone,
		Perfil:       domain.PerfilAdmin,
	}

	if err := h.usuarioService.Registrar(usuario, request.Senha); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remover dados sensíveis antes de retornar
	usuario.SenhaHash = ""

	c.JSON(http.StatusCreated, usuario)
}

// Login autentica um usuário e retorna um token JWT
func (h *AuthHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenResponse, err := h.usuarioService.Login(request.Email, request.Senha)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}
