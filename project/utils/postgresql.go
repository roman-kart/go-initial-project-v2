package utils

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/roman-kart/go-initial-project/project/config"
)

type Postgresql struct {
	Config *config.Config
	logger *zap.Logger
	Logger *Logger
}

func NewPostgresql(config *config.Config, logger *Logger) *Postgresql {
	return &Postgresql{
		Config: config,
		logger: logger.Logger.Named("Postgresql"),
		Logger: logger,
	}
}

func (p *Postgresql) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		p.Config.Postgresql.Host,
		p.Config.Postgresql.User,
		p.Config.Postgresql.Password,
		p.Config.Postgresql.Database,
		p.Config.Postgresql.Port,
	)
}

func (p *Postgresql) GetConnection() (*gorm.DB, error) {
	logger := p.logger.Named("GetConnection")

	dsn := p.GetConnectionString()
	logger.Info("dsn", zap.String("dsn", dsn))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}
	if p.Config.IsDebug {
		db = db.Debug()
	}
	return db, err
}

func (p *Postgresql) Migrate(models []interface{}) error {
	logger := p.logger.Named("Migrate")

	if !p.Config.Postgresql.AutoMigrate {
		logger.Info("AutoMigrate is disabled")
		return nil
	}

	db, err := p.GetConnection()
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return err
	}

	for _, model := range models {
		if p.Config.Postgresql.IsNeedToRecreate {
			err := db.Migrator().DropTable(model)
			if err != nil {
				logger.Error("Failed to drop table", zap.Error(err))
				return err
			}
		}
		err = db.AutoMigrate(model)
		if err != nil {
			logger.Error("Failed to auto migrate", zap.Error(err), zap.Any("model", model))
			return err
		}
	}

	return nil
}

type BasicPostgresqlModel struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
