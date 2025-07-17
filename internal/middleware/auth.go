package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Constante para armazenar o ID do usuário no contexto do Gin
const UsuarioIDKey = "usuario_id"
const PerfilKey = "perfil"

// Variável para o serviço de usuário, permite substituição para testes
var usuarioService UsuarioService

// UsuarioService define uma interface com o método ValidateToken
type UsuarioService interface {
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

// RequirePerfil é um middleware que verifica se o usuário tem um perfil específico ou um dos perfis permitidos
func RequirePerfil(perfisPermitidos ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se o usuário está autenticado
		_, exists := c.Get(UsuarioIDKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			return
		}

		// Obter o perfil do usuário do contexto
		perfilRaw, exists := c.Get(PerfilKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Perfil não disponível"})
			return
		}

		// Converter para string
		perfil, ok := perfilRaw.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar perfil"})
			return
		}

		// Verificar se o perfil está entre os permitidos
		permitido := false
		for _, p := range perfisPermitidos {
			if perfil == p {
				permitido = true
				break
			}
		}

		if !permitido {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acesso não autorizado para este perfil"})
			return
		}

		c.Next()
	}
}

// JWTAuthMiddleware é um middleware que valida tokens JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter o token do cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cabeçalho de autorização não fornecido"})
			return
		}

		// O token deve estar no formato "Bearer {token}"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			return
		}

		// Usar o serviço para validar o token
		// Se não houver serviço configurado, fallback para validação local
		var claims jwt.MapClaims
		var err error

		if usuarioService != nil {
			claims, err = usuarioService.ValidateToken(tokenString)
		} else {
			// Fallback para validação local quando não há serviço configurado
			const jwtSecret = "seu_segredo_jwt_aqui"
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validar o método de assinatura
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido: " + err.Error()})
				return
			}

			// Extrair claims do token
			if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				claims = tokenClaims
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
				return
			}
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido: " + err.Error()})
			return
		}

		// Processar as claims
		// Extrair o ID e perfil do usuário
		userID, ok := claims["id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token não contém ID de usuário"})
			return
		}

		// Converter string para UUID
		id, err := uuid.Parse(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ID de usuário inválido no token"})
			return
		}

		// Armazenar o ID e perfil do usuário no contexto
		c.Set(UsuarioIDKey, id)
		if perfil, ok := claims["perfil"].(string); ok {
			c.Set(PerfilKey, perfil)
		}

		c.Next()
	}
}
