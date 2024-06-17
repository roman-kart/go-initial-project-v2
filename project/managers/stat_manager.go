package managers

import (
	"time"

	"go.uber.org/zap"

	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/tools"
	"github.com/roman-kart/go-initial-project/project/utils"
)

// StatManager do CRUD operations with statistics.
type StatManager struct {
	Config              *config.Config
	Logger              *utils.Logger
	logger              *zap.Logger
	ClickHouse          *utils.ClickHouse
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewStatManager create new StatManager instance.
// Using for configuring with wire.
func NewStatManager(
	logger *utils.Logger,
	clickHouse *utils.ClickHouse,
	config *config.Config,
	errorWrapperCreator tools.ErrorWrapperCreator,
) (*StatManager, error) {
	sm := &StatManager{
		Config:              config,
		Logger:              logger,
		logger:              logger.Logger,
		ClickHouse:          clickHouse,
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("StatManager"),
	}

	ew := tools.GetErrorWrapper("NewStatManager")

	err := sm.migrate()
	if err != nil {
		return nil, ew(err)
	}

	return sm, nil
}

func (sm *StatManager) migrate() error {
	ew := sm.ErrorWrapperCreator.GetMethodWrapper("migrate")
	err := sm.ClickHouse.Migrate([]interface{}{&ApplicationStatsModel{}})

	return ew(err)
}

// Add new event.
func (sm *StatManager) Add(eventName string, eventDateTime time.Time, eventMessage string) error {
	ew := sm.ErrorWrapperCreator.GetMethodWrapper("Add")

	db, err := sm.ClickHouse.GetConnection()
	if err != nil {
		return ew(err)
	}

	db = db.Create(&ApplicationStatsModel{
		EventName:     eventName,
		EventDate:     eventDateTime,
		EventDateTime: eventDateTime,
		EventMessage:  eventMessage,
	})

	err = db.Error
	if err != nil {
		return ew(err)
	}

	return nil
}

// AddSimple add event with current datetime.
func (sm *StatManager) AddSimple(eventName string, eventMessage string) error {
	return sm.Add(eventName, time.Now(), eventMessage)
}

// ApplicationStatsModel contains statistics data.
type ApplicationStatsModel struct {
	EventName     string    `gorm:"type:String"   my_clickhouse:"order_by=1;primary_key=1"`
	EventDate     time.Time `gorm:"type:date"     my_clickhouse:"order_by=2;primary_key=2"`
	EventDateTime time.Time `gorm:"type:datetime" my_clickhouse:"order_by=3"`
	EventMessage  string    `gorm:"type:String"   my_clickhouse:"order_by=4"` // Can be of any size
}
