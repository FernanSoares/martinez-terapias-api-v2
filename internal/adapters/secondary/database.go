package secondary

import (
	"martinezterapias/api/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDatabaseConnection cria uma nova conexão com o banco de dados
func NewDatabaseConnection(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Executar migrações automaticamente
	err = db.AutoMigrate(
		&UserModel{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
