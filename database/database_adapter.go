package database

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DatabaseAdapter struct {
	*gorm.DB
}

func NewDatabaseAdapter(server, user, password string, port int, database string, minConn, maxConn int, useLocalTime bool) (*DatabaseAdapter, error) {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		user,
		password,
		server,
		port,
		database,
	)
	if useLocalTime {
		dataSource = dataSource + "&loc=Local"
	}

	conn, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{
		Logger: gormlogger.Discard,
	})
	if err != nil {
		return nil, err
	}

	db, err := conn.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get *sql.DB")
	}

	db.SetMaxIdleConns(minConn)
	db.SetMaxOpenConns(maxConn)

	return &DatabaseAdapter{DB: conn}, nil
}

func (a *DatabaseAdapter) FindDatabases() ([]string, error) {
	type Result struct {
		Database string `gorm:"column:Database"`
	}
	var results []*Result

	if err := a.Raw("show databases").Scan(&results).Error; err != nil {
		return nil, err
	}

	var databases []string
	for _, result := range results {
		databases = append(databases, result.Database)
	}
	return databases, nil
}