package utils

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/roman-kart/go-initial-project/v2/project/config"
	"github.com/roman-kart/go-initial-project/v2/project/tools"
)

// Postgresql manipulates connection to Postgresql database.
type Postgresql struct {
	Config              *config.Config
	logger              *zap.Logger
	Logger              *Logger
	db                  *gorm.DB
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewPostgresql creates new instance of [Postgresql].
// Using for configuring with wire.
func NewPostgresql(
	config *config.Config,
	logger *Logger,
	errorWrapperCreator tools.ErrorWrapperCreator,
) (*Postgresql, func(), error) {
	p := &Postgresql{
		Config:              config,
		logger:              logger.Logger.Named("Postgresql"),
		Logger:              logger,
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("Postgresql"),
	}

	ew := tools.GetErrorWrapper("NewPostgresql")

	_, err := p.GetConnection()
	if err != nil {
		return nil, nil, ew(err)
	}

	return p, func() {
		db, err := p.db.DB()
		if err != nil {
			p.logger.Error("Error while getting db connection", zap.Error(ew(err)))
		}

		err = db.Close()
		if err != nil {
			p.logger.Error("Error while closing db connection", zap.Error(ew(err)))
		}
	}, nil
}

// GetConnectionString returns formated connection string.
func (p *Postgresql) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		p.Config.Postgresql.Host,
		p.Config.Postgresql.User,
		p.Config.Postgresql.Password,
		p.Config.Postgresql.Database,
		p.Config.Postgresql.Port,
	)
}

// GetConnection create new connection with caching.
// If connection is not cached, it will be created.
//
//nolint:dupl
func (p *Postgresql) GetConnection() (*gorm.DB, error) {
	ew := p.ErrorWrapperCreator.GetMethodWrapper("GetConnection")
	logger := p.logger.Named("GetConnection")

	if p.db != nil {
		return p.db, nil
	}

	dsn := p.GetConnectionString()

	logger.Info("dsn", zap.String("dsn", dsn))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, ew(err)
	}

	if p.Config.IsDebug {
		db = db.Debug()
	}

	dbInner, err := db.DB()
	if err != nil {
		return nil, ew(err)
	}

	dbInner.SetConnMaxLifetime(time.Second * time.Duration(p.Config.Postgresql.ConnMaxLifetime))
	dbInner.SetConnMaxIdleTime(time.Second * time.Duration(p.Config.Postgresql.ConnMaxIdleTime))
	dbInner.SetMaxIdleConns(p.Config.Postgresql.MaxIdleConns)
	dbInner.SetMaxOpenConns(p.Config.Postgresql.MaxOpenConns)

	p.db = db

	return db, nil
}

// Migrate models to Postgresql.
// Depends on Clickhouse.AutoMigrate parameter of [config.Config].
func (p *Postgresql) Migrate(models []interface{}) error {
	ew := p.ErrorWrapperCreator.GetMethodWrapper("Migrate")
	logger := p.logger.Named("Migrate")

	if !p.Config.Postgresql.AutoMigrate {
		logger.Info("AutoMigrate is disabled")
		return nil
	}

	db, err := p.GetConnection()
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return ew(err)
	}

	for _, model := range models {
		if p.Config.Postgresql.IsNeedToRecreate {
			err := db.Migrator().DropTable(model)
			if err != nil {
				logger.Error("Failed to drop table", zap.Error(err))
				return ew(err)
			}
		}

		err = db.AutoMigrate(model)
		if err != nil {
			logger.Error("Failed to auto migrate", zap.Error(err), zap.Any("model", model))
			return ew(err)
		}
	}

	return nil
}

// BasicPostgresqlModel is a basic model for PostgresSQL database.
type BasicPostgresqlModel struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
