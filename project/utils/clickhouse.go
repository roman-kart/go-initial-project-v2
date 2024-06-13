package utils

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm/schema"

	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"

	cfg "github.com/roman-kart/go-initial-project/project/config"
)

type ClickHouse struct {
	Config *cfg.Config
	logger *zap.Logger
	Logger *Logger
	db     *gorm.DB
}

func NewClickHouse(config *cfg.Config, logger *Logger) *ClickHouse {
	return &ClickHouse{
		Config: config,
		logger: logger.Logger.Named("ClickHouse"),
		Logger: logger,
	}
}

func (c *ClickHouse) GetConnectionString() string {
	return fmt.Sprintf("tcp://%s:%d/%s?username=%s",
		c.Config.Clickhouse.Host,
		c.Config.Clickhouse.Port,
		c.Config.Clickhouse.Database,
		c.Config.Clickhouse.User,
	)
}

func (c *ClickHouse) GetConnection() (*gorm.DB, error) {
	logger := c.logger.Named("GetConnection")

	if c.db != nil {
		return c.db, nil
	}

	dsn := c.GetConnectionString()
	logger.Info("dsn", zap.String("dsn", dsn))
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if c.Config.IsDebug {
		db = db.Debug()
	}

	dbInner, err := db.DB()
	if err != nil {
		return nil, err
	}
	dbInner.SetConnMaxLifetime(time.Second * time.Duration(c.Config.Clickhouse.ConnMaxLifetime))
	dbInner.SetConnMaxIdleTime(time.Second * time.Duration(c.Config.Clickhouse.ConnMaxIdleTime))
	dbInner.SetMaxIdleConns(c.Config.Clickhouse.MaxIdleConns)
	dbInner.SetMaxOpenConns(c.Config.Clickhouse.MaxOpenConns)

	c.db = db
	return db, err
}

func (c *ClickHouse) Migrate(models []interface{}) error {
	logger := c.logger.Named("Migrate")

	if !c.Config.Clickhouse.AutoMigrate {
		logger.Info("AutoMigrate is disabled")
		return nil
	}

	db, err := c.GetConnection()
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return err
	}

	var tableMigrateEntitiesCh []TableMigrateEntityClickhouse
	for _, model := range models {
		tags, err := RetrieveClickhouseTags(model)
		if err != nil {
			logger.Error("Failed to retrieve tags", zap.Error(err))
			return err
		}

		option, err := BuildOptionsFromClickhouseTags(tags)
		if err != nil {
			logger.Error("Failed to build options from tags", zap.Error(err))
			return err
		}

		tableMigrateEntitiesCh = append(tableMigrateEntitiesCh, TableMigrateEntityClickhouse{
			Model:   &model,
			Options: option,
		})
	}

	for _, entity := range tableMigrateEntitiesCh {
		logger.Info("Migrate", zap.String("model", reflect.TypeOf(entity.Model).String()))
		if c.Config.Clickhouse.IsNeedToRecreate {
			logger.Info("Model is need to recreate", zap.String("model", reflect.TypeOf(entity.Model).String()))
			err := db.Migrator().DropTable(entity.Model)
			if err != nil {
				logger.Error("Failed to drop table", zap.Error(err))
				return err
			}
		}
		if err := db.Set("gorm:table_options", entity.Options).AutoMigrate(entity.Model); err != nil {
			logger.Error("Failed to migrate table", zap.Error(err))
			return err
		}
	}

	return nil
}

type ClickhouseFieldTag struct {
	Name  string
	Value string
}

type ClickhouseField struct {
	Name   string
	DbName string
	Tags   []ClickhouseFieldTag
}

func RetrieveClickhouseTags(model interface{}) ([]ClickhouseField, error) {
	s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, err
	}
	modelAndDbNames := make(map[string]string)
	for _, field := range s.Fields {
		dbName := field.DBName
		modelName := field.Name
		modelAndDbNames[modelName] = dbName
	}

	var resultTags []ClickhouseField
	reflectVal := reflect.ValueOf(model)
	for i := 0; i < reflectVal.NumField(); i++ {
		tag := reflectVal.Type().Field(i).Tag.Get("my_clickhouse")
		fieldModelName := reflectVal.Type().Field(i).Name
		dbName := modelAndDbNames[fieldModelName]
		subtags := strings.Split(tag, ";")

		clickhouseField := ClickhouseField{
			Name:   fieldModelName,
			DbName: dbName,
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

func BuildOptionsFromClickhouseTags(fields []ClickhouseField) (string, error) {
	primaryKeys := make(map[string]string)
	orderBy := make(map[string]string)
	for _, field := range fields {
		tags := field.Tags
		for _, tag := range tags {
			if tag.Name == "primary_key" {
				primaryKeys[tag.Value] = field.DbName
			}
			if tag.Name == "order_by" {
				orderBy[tag.Value] = field.DbName
			}
		}
	}

	primaryKeyIds := SortMapKeys(primaryKeys)
	orderByIds := SortMapKeys(orderBy)

	primaryKeyPartStr := ""
	if len(primaryKeyIds) > 0 {
		primaryKeyPartStr = "PRIMARY KEY ("
		for _, primaryKeyId := range primaryKeyIds {
			primaryKeyPartStr += primaryKeys[primaryKeyId] + ", "
		}
		primaryKeyPartStr = primaryKeyPartStr[:len(primaryKeyPartStr)-2] + ")"
	}
	orderByPartStr := ""
	if len(orderByIds) > 0 {
		orderByPartStr = "ORDER BY ("
		for _, orderById := range orderByIds {
			orderByPartStr += orderBy[orderById] + ", "
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

type TableMigrateEntityClickhouse struct {
	Options string
	Model   interface{}
}
