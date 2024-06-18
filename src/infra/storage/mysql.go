package storage

import (
	"context"
	"fmt"
	"github.com/mjedari/health-checker/app/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL(conf config.MySQL) *MySQL {
	// create a connection to mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", conf.User, conf.Pass, conf.Host, conf.Port, conf.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
