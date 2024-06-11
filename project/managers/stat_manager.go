package managers

import (
	"time"

	"github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/utils"
	"go.uber.org/zap"
)

type StatManager struct {
	Config     *config.Config
	Logger     *utils.Logger
	logger     *zap.Logger
	ClickHouse *utils.ClickHouse
}

func NewStatManager(logger *utils.Logger, clickHouse *utils.ClickHouse, config *config.Config) *StatManager {
	return &StatManager{
		Config:     config,
		Logger:     logger,
		logger:     logger.Logger,
		ClickHouse: clickHouse,
	}
}

func (sm *StatManager) Prepare() error {
	err := sm.migrate()
	return err
}

func (sm *StatManager) migrate() error {
	err := sm.ClickHouse.Migrate([]interface{}{&ApplicationStatsModel{}})
	return err
}

func (sm *StatManager) Add(eventName string, eventDateTime time.Time, eventMessage string) error {
	db, err := sm.ClickHouse.GetConnection()
	if err != nil {
		return err
	}
	db = db.Create(&ApplicationStatsModel{
		EventName:     eventName,
		EventDate:     eventDateTime,
		EventDateTime: eventDateTime,
		EventMessage:  eventMessage,
	})
	err = db.Error
	if err != nil {
		return err
	}
	return err
}

// AddSimple add event with current datetime
func (sm *StatManager) AddSimple(eventName string, eventMessage string) error {
	return sm.Add(eventName, time.Now(), eventMessage)
}

type ApplicationStatsModel struct {
	EventName     string    `gorm:"type:String" my_clickhouse:"order_by=1;primary_key=1"`
	EventDate     time.Time `gorm:"type:date" my_clickhouse:"order_by=2;primary_key=2"`
	EventDateTime time.Time `gorm:"type:datetime" my_clickhouse:"order_by=3"`
	EventMessage  string    `gorm:"type:String" my_clickhouse:"order_by=4"`
}
