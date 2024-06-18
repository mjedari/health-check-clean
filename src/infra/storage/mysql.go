package storage

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL() *MySQL {
	// create a connection to mysql

	dbName := "health_db"
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("root:2231218/m@tcp(127.0.0.1:3306)/%s?parseTime=true", dbName)), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	return &MySQL{db: db}
}

func (s *MySQL) Create(ctx context.Context, value interface{}) error {
	return s.db.WithContext(ctx).Create(value).Error
}

func (s *MySQL) Read(ctx context.Context, id uint, out interface{}) error {
	return s.db.WithContext(ctx).First(out, id).Error
}

func (s *MySQL) ReadAll(ctx context.Context, out interface{}) error {
	return s.db.WithContext(ctx).Find(out).Error
}

func (s *MySQL) Update(ctx context.Context, value interface{}) error {
	return s.db.WithContext(ctx).Save(value).Error
}

func (s *MySQL) Delete(ctx context.Context, value interface{}) error {
	return s.db.WithContext(ctx).Delete(value).Error
}
