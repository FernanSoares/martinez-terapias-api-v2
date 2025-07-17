package config

import "os"

// Config contém as configurações da aplicação
type Config struct {
	ServerPort  string
	DatabaseURL string
	JWTSecret   string
}

// NewConfig cria uma nova instância de Config
func NewConfig() *Config {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "martinezterapias.db"
	}

	// WhatsApp configuration removed

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "seu_segredo_jwt_aqui" // Em produção, deve ser uma chave segura
	}

	return &Config{
		ServerPort:  serverPort,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
	}
}
