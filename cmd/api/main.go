package main

import (
	"log"
	"martinezterapias/api/internal/adapters/primary"
	"martinezterapias/api/internal/adapters/secondary"
	"martinezterapias/api/internal/core/services"
	"martinezterapias/api/internal/config"
	"martinezterapias/api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Carregar configurações
	cfg := config.NewConfig()

	// Inicializar repositório (adaptador secundário)
	db, err := secondary.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	// Inicializar repositórios para compatibilidade com código existente
	userRepository := secondary.NewUserRepository(db)

	// Inicializar repositórios da clínica
	clienteRepository := secondary.NewClienteRepository(db)
	servicoRepository := secondary.NewServicoRepository(db)
	agendamentoRepository := secondary.NewAgendamentoRepository(db)
	anamneseRepository := secondary.NewAnamneseRepository(db)
	usuarioRepository := secondary.NewUsuarioRepository(db)


	// Inicializar serviços
	userService := services.NewUserService(userRepository) // Manter compatibilidade

	clienteService := services.NewClienteService(clienteRepository)
	servicoService := services.NewServicoService(servicoRepository)
	agendamentoService := services.NewAgendamentoService(
		agendamentoRepository,
		clienteRepository,
		servicoRepository,
		usuarioRepository,
	)
	anamneseService := services.NewAnamneseService(anamneseRepository, clienteRepository)
	usuarioService := services.NewUsuarioService(usuarioRepository, agendamentoRepository)

	// Configurar o router Gin
	router := gin.Default()

	// Configurar CORS para frontend
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Inicializar handlers (adaptadores primários)
	userHandler := primary.NewUserHandler(userService) // Manter compatibilidade
	authHandler := primary.NewAuthHandler(usuarioService)
	clienteHandler := primary.NewClienteHandler(clienteService)
	servicoHandler := primary.NewServicoHandler(servicoService)
	agendamentoHandler := primary.NewAgendamentoHandler(agendamentoService)
	anamneseHandler := primary.NewAnamneseHandler(anamneseService)
	perfilHandler := primary.NewPerfilHandler(usuarioService)

	// Configurar rotas
	api := router.Group("/api")
	{
		// Rotas públicas (sem autenticação)
		api.POST("/registrar", authHandler.Registrar)
		api.POST("/login", authHandler.Login)
		api.POST("/registrar-admin", authHandler.RegistrarAdmin)

		// Rotas protegidas por autenticação
		auth := api.Group("/")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			// Rotas para acessar o próprio perfil e agendamentos (para clientes)
			me := auth.Group("/me")
			{
				me.GET("/perfil", perfilHandler.GetPerfil)
				me.GET("/agendamentos", perfilHandler.GetAgendamentos)
			}

			// Rota para solicitar reagendamento
			auth.POST("/agendamentos/:id/solicitar-reagendamento", agendamentoHandler.SolicitarReagendamento)

			// Rotas que requerem perfil de administrador ou massoterapeuta
			adminOrMasso := auth.Group("/")
			adminOrMasso.Use(middleware.RequirePerfil("admin", "massoterapeuta"))
			{
				// Rotas para clientes
				clientes := adminOrMasso.Group("/clientes")
				{
					clientes.GET("/", clienteHandler.GetAll)
					clientes.GET("/:id", clienteHandler.GetByID)
					clientes.POST("/", clienteHandler.Create)
					clientes.PUT("/:id", clienteHandler.Update)
					clientes.DELETE("/:id", clienteHandler.Delete)

				}

				// Rotas para anamnese associadas a clientes (em grupo separado para evitar conflito com :id)
				clienteAnamnese := adminOrMasso.Group("/clientes-anamnese")
				{
					clienteAnamnese.GET("/:cliente_id", anamneseHandler.GetByClienteID)
					clienteAnamnese.POST("/:cliente_id", anamneseHandler.Create)
				}

				// Rotas para serviços
				servicos := adminOrMasso.Group("/servicos")
				{
					servicos.GET("/", servicoHandler.GetAll)
					servicos.GET("/:id", servicoHandler.GetByID)
					servicos.POST("/", servicoHandler.Create)
					servicos.PUT("/:id", servicoHandler.Update)
					servicos.DELETE("/:id", servicoHandler.Delete)
				}

				// Rotas para agendamentos
				agendamentos := adminOrMasso.Group("/agendamentos")
				{
					agendamentos.GET("/", agendamentoHandler.GetAll)
					agendamentos.GET("/:id", agendamentoHandler.GetByID)
					agendamentos.POST("/", agendamentoHandler.Create)
					agendamentos.PUT("/:id", agendamentoHandler.Update)
					agendamentos.PATCH("/:id/status", agendamentoHandler.UpdateStatus)
					agendamentos.DELETE("/:id", agendamentoHandler.Delete)
				}

				// Rotas para fichas de anamnese
				anamnese := adminOrMasso.Group("/anamnese")
				{
					anamnese.GET("/", anamneseHandler.GetAll)
					anamnese.GET("/:id", anamneseHandler.GetByID)
					anamnese.PUT("/:id", anamneseHandler.Update)
				}
			}

			// Rotas para usuários (compatibilidade com frontend)
			usuarios := auth.Group("/usuarios")
			{
				usuarios.GET("/", userHandler.GetAll) // Reutilizar handler existente
				usuarios.GET("/:id", userHandler.GetByID)
			}

			// Rotas que requerem apenas perfil de administrador
			admin := auth.Group("/admin")
			admin.Use(middleware.RequirePerfil("admin"))
			{
				// Aqui poderíamos adicionar rotas específicas para administradores
				// como gerenciamento de usuários massoterapeutas
			}
		}
	}

	// Mantendo as rotas legadas para compatibilidade
	users := api.Group("/users")
	{
		users.GET("/", userHandler.GetAll)
		users.GET("/:id", userHandler.GetByID)
		users.POST("/", userHandler.Create)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}

	// Iniciar o servidor
	log.Printf("Servidor iniciado na porta %s", cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}
