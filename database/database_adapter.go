package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

	conn, err := gorm.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	conn.LogMode(false)
	conn.DB().SetMaxIdleConns(minConn)
	conn.DB().SetMaxOpenConns(maxConn)

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