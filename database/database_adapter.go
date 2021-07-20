package database

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DatabaseAdapter struct {
	*gorm.DB
}

func NewMySQLDatabaseAdapter(server, user, password string, port int, database string, minConn, maxConn int, useLocalTime bool) (*DatabaseAdapter, error) {
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

	return newDatabaseAdapter(mysql.Open(dataSource), minConn, maxConn)
}

func NewPostgresDatabaseAdapter(server, user, password string, port int, database string, minConn, maxConn int, timeZonePtr *string) (*DatabaseAdapter, error) {
	timeZone := "Asia/Seoul"
	if timeZonePtr != nil {
		timeZone = *timeZonePtr
	}

	dataSource := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		server,
		user,
		password,
		database,
		port,
		timeZone,
	)

	return newDatabaseAdapter(postgres.Open(dataSource), minConn, maxConn)
}

func newDatabaseAdapter(dialector gorm.Dialector, minConn, maxConn int) (*DatabaseAdapter, error) {
	conn, err := gorm.Open(dialector, &gorm.Config{
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
