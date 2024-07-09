package utils

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/roman-kart/go-initial-project/v2/project/tools"
)

type ClickHouseConfig struct {
	Host               string
	Port               int
	User               string
	Password           string
	Database           string
	IsNeedToRecreate   bool
	AutoMigrate        bool
	IsNeedToInitialize bool
	ConnMaxLifetime    int64
	ConnMaxIdleTime    int64
	MaxIdleConns       int
	MaxOpenConns       int
	IsDebug            bool
}

// ClickHouse manipulates connection to ClickHouse database.
type ClickHouse struct {
	Config              *ClickHouseConfig
	logger              *zap.Logger
	db                  *gorm.DB
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewClickHouse creates new instance of [ClickHouse].
// Using for configuring with wire.
func NewClickHouse(
	config *ClickHouseConfig,
	logger *zap.Logger,
	errorWrapperCreator tools.ErrorWrapperCreator,
) (*ClickHouse, func(), error) {
	c := &ClickHouse{
		Config:              config,
		logger:              logger.Named("ClickHouse"),
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("ClickHouse"),
	}

	ew := tools.GetErrorWrapper("NewClickHouse")

	_, err := c.GetConnection()
	if err != nil {
		return nil, nil, ew(err)
	}

	return c, func() {
		db, err := c.db.DB()
		if err != nil {
			c.logger.Error("Error while getting db connection", zap.Error(err))
		}

		err = db.Close()
		if err != nil {
			c.logger.Error("Error while closing db connection", zap.Error(err))
		}
	}, nil
}

// GetConnectionString returns formated connection string.
func (c *ClickHouse) GetConnectionString() string {
	hostAndPort := net.JoinHostPort(
		c.Config.Host,
		strconv.Itoa(c.Config.Port),
	)

	return fmt.Sprintf("tcp://%s/%s?username=%s",
		hostAndPort,
		c.Config.Database,
		c.Config.User,
	)
}

// GetConnection create connection to DB with caching.
// If connection is not cached, it will be created.
//
//nolint:dupl
func (c *ClickHouse) GetConnection() (*gorm.DB, error) {
	ew := c.ErrorWrapperCreator.GetMethodWrapper("GetConnection")
	logger := c.logger.Named("GetConnection")

	if c.db != nil {
		return c.db, nil
	}

	dsn := c.GetConnectionString()

	logger.Info("dsn", zap.String("dsn", dsn))

	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, ew(err)
	}

	if c.Config.IsDebug {
		db = db.Debug()
	}

	dbInner, err := db.DB()
	if err != nil {
		return nil, ew(err)
	}

	dbInner.SetConnMaxLifetime(time.Second * time.Duration(c.Config.ConnMaxLifetime))
	dbInner.SetConnMaxIdleTime(time.Second * time.Duration(c.Config.ConnMaxIdleTime))
	dbInner.SetMaxIdleConns(c.Config.MaxIdleConns)
	dbInner.SetMaxOpenConns(c.Config.MaxOpenConns)

	c.db = db

	return db, nil
}

// Migrate models to ClickHouse.
// Depends on Clickhouse.AutoMigrate parameter of [cfg.Config].
func (c *ClickHouse) Migrate(models []interface{}) error {
	ew := c.ErrorWrapperCreator.GetMethodWrapper("Migrate")
	logger := c.logger.Named("Migrate")

	if !c.Config.AutoMigrate {
		logger.Info("AutoMigrate is disabled")
		return nil
	}

	db, err := c.GetConnection()
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return ew(err)
	}

	tableMigrateEntitiesCh := []TableMigrateEntityClickhouse{}

	for _, model := range models {
		tags, err := RetrieveClickhouseTags(model)
		if err != nil {
			logger.Error("Failed to retrieve tags", zap.Error(err))
			return ew(err)
		}

		option, err := BuildOptionsFromClickhouseTags(tags)
		if err != nil {
			logger.Error("Failed to build options from tags", zap.Error(err))
			return ew(err)
		}

		tableMigrateEntitiesCh = append(tableMigrateEntitiesCh, TableMigrateEntityClickhouse{
			Model:   &model,
			Options: option,
		})
	}

	for _, entity := range tableMigrateEntitiesCh {
		logger.Info("Migrate", zap.String("model", reflect.TypeOf(entity.Model).String()))

		if c.Config.IsNeedToRecreate {
			logger.Info("Model is need to recreate", zap.String("model", reflect.TypeOf(entity.Model).String()))

			err := db.Migrator().DropTable(entity.Model)
			if err != nil {
				logger.Error("Failed to drop table", zap.Error(err))
				return ew(err)
			}
		}

		if err := db.Set("gorm:table_options", entity.Options).AutoMigrate(entity.Model); err != nil {
			logger.Error("Failed to migrate table", zap.Error(err))
			return ew(err)
		}
	}

	return nil
}

// ClickhouseFieldTag contains name and value of tag.
type ClickhouseFieldTag struct {
	Name  string
	Value string
}

// ClickhouseField contains property name original and in database and set of property's tags.
type ClickhouseField struct {
	Name   string
	DBName string
	Tags   []ClickhouseFieldTag
}

// RetrieveClickhouseTags retrieves tags of model.
func RetrieveClickhouseTags(model interface{}) ([]ClickhouseField, error) {
	ew := tools.GetErrorWrapper("RetrieveClickhouseTags")

	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, ew(err)
	}

	modelAndDBNames := make(map[string]string)

	for _, field := range s.Fields {
		dbName := field.DBName
		modelName := field.Name
		modelAndDBNames[modelName] = dbName
	}

	resultTags := []ClickhouseField{}

	reflectVal := reflect.ValueOf(model)
	if reflectVal.Kind() == reflect.Pointer {
		reflectVal = reflectVal.Elem()
	}

	for i := range reflectVal.NumField() {
		tag := reflectVal.Type().Field(i).Tag.Get("my_clickhouse")
		if tag == "" {
			continue
		}

		fieldModelName := reflectVal.Type().Field(i).Name
		dbName := modelAndDBNames[fieldModelName]
		subtags := strings.Split(tag, ";")

		clickhouseField := ClickhouseField{
			Name:   fieldModelName,
			DBName: dbName,
			Tags:   []ClickhouseFieldTag{},
		}

		for _, subtag := range subtags {
			keyAndVal := strings.Split(subtag, "=")
			key := keyAndVal[0]
			val := keyAndVal[1]
			clickhouseField.Tags = append(clickhouseField.Tags, ClickhouseFieldTag{
				Name:  key,
				Value: val,
			})
		}

		resultTags = append(resultTags, clickhouseField)
	}

	return resultTags, nil
}

// BuildOptionsFromClickhouseTags builds options part CREATE TABLE query tags.
func BuildOptionsFromClickhouseTags(fields []ClickhouseField) (string, error) {
	primaryKeys := make(map[string]string)
	orderBy := make(map[string]string)

	for _, field := range fields {
		tags := field.Tags
		for _, tag := range tags {
			if tag.Name == "primary_key" {
				primaryKeys[tag.Value] = field.DBName
			}

			if tag.Name == "order_by" {
				orderBy[tag.Value] = field.DBName
			}
		}
	}

	primaryKeyIDs := tools.SortMapKeys(primaryKeys)
	orderByIDs := tools.SortMapKeys(orderBy)

	primaryKeyPartStr := ""
	if len(primaryKeyIDs) > 0 {
		primaryKeyPartStr = "PRIMARY KEY ("
		for _, primaryKeyID := range primaryKeyIDs {
			primaryKeyPartStr += primaryKeys[primaryKeyID] + ", "
		}

		primaryKeyPartStr = primaryKeyPartStr[:len(primaryKeyPartStr)-2] + ")"
	}

	orderByPartStr := ""

	if len(orderByIDs) > 0 {
		orderByPartStr = "ORDER BY ("
		for _, orderByID := range orderByIDs {
			orderByPartStr += orderBy[orderByID] + ", "
		}

		orderByPartStr = orderByPartStr[:len(orderByPartStr)-2] + ")"
	}

	options := fmt.Sprintf(`
ENGINE MergeTree 
%s
%s
`, primaryKeyPartStr, orderByPartStr)

	return options, nil
}

// TableMigrateEntityClickhouse contains model for migrate and options part of CREATE TABLE query.
type TableMigrateEntityClickhouse struct {
	Options string
	Model   interface{}
}
