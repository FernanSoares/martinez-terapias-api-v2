package primary

import (
	"net/http"
	"time"

	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ClienteHandler implementa os handlers HTTP para clientes
type ClienteHandler struct {
	clienteService ports.ClienteService
}

// NewClienteHandler cria uma nova instância de ClienteHandler
func NewClienteHandler(clienteService ports.ClienteService) *ClienteHandler {
	return &ClienteHandler{
		clienteService: clienteService,
	}
}

// GetAll retorna todos os clientes
// @Summary Listar clientes
// @Description Retorna a lista de todos os clientes
// @Tags clientes
// @Produce json
// @Param nome query string false "Filtrar por nome"
// @Param ativo query bool false "Filtrar por status (ativo/inativo)"
// @Success 200 {array} domain.Cliente
// @Failure 500 {object} map[string]string
// @Router /clientes [get]
func (h *ClienteHandler) GetAll(c *gin.Context) {
	// Verificar parâmetros de consulta
	nome := c.Query("nome")
	ativoParam := c.Query("ativo")

	var clientes []domain.Cliente
	var err error

	// Filtrar por nome se fornecido
	if nome != "" {
		clientes, err = h.clienteService.GetByNome(nome)
	} else if ativoParam != "" {
		ativo := ativoParam == "true"
		clientes, err = h.clienteService.GetByAtivo(ativo)
	} else {
		clientes, err = h.clienteService.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clientes)
}

// GetByID retorna um cliente pelo ID
// @Summary Obter cliente por ID
// @Description Retorna os detalhes de um cliente específico
// @Tags clientes
// @Produce json
// @Param id path string true "ID do cliente"
// @Success 200 {object} domain.Cliente
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clientes/{id} [get]
func (h *ClienteHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	cliente, err := h.clienteService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}

	c.JSON(http.StatusOK, cliente)
}

// Create cria um novo cliente
// @Summary Cadastrar cliente
// @Description Cria um novo registro de cliente
// @Tags clientes
// @Accept json
// @Produce json
// @Param cliente body domain.Cliente true "Dados do cliente"
// @Success 201 {object} domain.Cliente
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clientes [post]
func (h *ClienteHandler) Create(c *gin.Context) {
	// Usar um mapa para receber os dados JSON brutos
	var rawData map[string]interface{}
	if err := c.ShouldBindJSON(&rawData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Criar o objeto cliente com os dados recebidos
	cliente := domain.Cliente{}

	// Preencher os campos básicos
	if nomeCompleto, ok := rawData["nome_completo"].(string); ok {
		cliente.NomeCompleto = nomeCompleto
	}
	if email, ok := rawData["email"].(string); ok {
		cliente.Email = email
	}
	if cpf, ok := rawData["cpf"].(string); ok {
		cliente.CPF = cpf
	}
	if telefone, ok := rawData["telefone"].(string); ok {
		cliente.Telefone = telefone
	}
	
	// Tratar a data de nascimento especificamente
	if dataNascStr, ok := rawData["data_nascimento"].(string); ok && dataNascStr != "" {
		// Primeiro, tentar fazer o parse no formato ISO YYYY-MM-DD
		dataNasc, err := time.Parse("2006-01-02", dataNascStr)
		if err == nil {
			cliente.DataNascimento = dataNasc
		} else {
			// Se falhar, tentar o formato brasileiro DD/MM/YYYY
			dataNasc, err = time.Parse("02/01/2006", dataNascStr)
			if err == nil {
				cliente.DataNascimento = dataNasc
			} else {
				// Se ainda falhar, usar data zero
				c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use DD/MM/AAAA."})
				return
			}
		}
	} else {
		cliente.DataNascimento = time.Time{}
	}

	// Validações básicas
	if cliente.NomeCompleto == "" || cliente.Email == "" || cliente.Telefone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome, email e telefone são obrigatórios"})
		return
	}

	if err := h.clienteService.Create(&cliente); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cliente)
}

// Update atualiza um cliente existente
// @Summary Atualizar cliente
// @Description Atualiza os dados de um cliente existente
// @Tags clientes
// @Accept json
// @Produce json
// @Param id path string true "ID do cliente"
// @Param cliente body domain.Cliente true "Dados do cliente"
// @Success 200 {object} domain.Cliente
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clientes/{id} [put]
func (h *ClienteHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se o cliente existe
	_, err = h.clienteService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}

	var cliente domain.Cliente
	if err := c.ShouldBindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Garantir que o ID na URL corresponde ao ID no corpo da requisição
	cliente.ID = id

	if err := h.clienteService.Update(&cliente); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cliente)
}

// Delete realiza o soft delete de um cliente
// @Summary Excluir cliente
// @Description Realiza o soft delete de um cliente (marcando como inativo)
// @Tags clientes
// @Produce json
// @Param id path string true "ID do cliente"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clientes/{id} [delete]
func (h *ClienteHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.clienteService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
