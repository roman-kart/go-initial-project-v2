package managers

import (
	"go.uber.org/zap"

	"github.com/roman-kart/go-initial-project/v2/project/config"
	"github.com/roman-kart/go-initial-project/v2/project/tools"
	"github.com/roman-kart/go-initial-project/v2/project/utils"
)

// UserAccountManager do CRUD operations on user accounts.
type UserAccountManager struct {
	Config              *config.Config
	Logger              *utils.Logger
	logger              *zap.Logger
	Postgresql          *utils.Postgresql
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewUserAccountManager creates a new user account manager.
// Using for configuring with wire.
func NewUserAccountManager(
	logger *utils.Logger,
	postgresql *utils.Postgresql,
	config *config.Config,
	errorWrapperCreator tools.ErrorWrapperCreator,
) (*UserAccountManager, error) {
	uam := &UserAccountManager{
		Config:              config,
		Logger:              logger,
		logger:              logger.Logger.Named("UserManager"),
		Postgresql:          postgresql,
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("UserAccountManager"),
	}

	ew := tools.GetErrorWrapper("NewUserAccountManager")

	err := uam.migrate()
	if err != nil {
		return nil, ew(err)
	}

	return uam, nil
}

func (m *UserAccountManager) migrate() error {
	ew := m.ErrorWrapperCreator.GetMethodWrapper("migrate")

	err := m.Postgresql.Migrate([]interface{}{UserAccount{}})
	if err != nil {
		return ew(err)
	}

	return nil
}

// UserAccount contains information of a user.
type UserAccount struct {
	utils.BasicPostgresqlModel
	Nickname string
}
