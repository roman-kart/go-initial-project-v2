package managers

import (
	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/utils"
	"go.uber.org/zap"
)

type UserAccountManager struct {
	ConConfig  *config.Config
	Logger     *utils.Logger
	logger     *zap.Logger
	Postgresql *utils.Postgresql
}

func NewUserManager(logger *utils.Logger, postgresql *utils.Postgresql, Config *config.Config) *UserAccountManager {
	return &UserAccountManager{
		ConConfig:  Config,
		Logger:     logger,
		logger:     logger.Logger,
		Postgresql: postgresql,
	}
}

func (m *UserAccountManager) Prepare() error {
	return m.migrate()
}

func (m *UserAccountManager) migrate() error {
	return m.Postgresql.Migrate([]interface{}{UserAccount{}})
}

type UserAccount struct {
	utils.BasicPostgresqlModel
	Nickname string
}
