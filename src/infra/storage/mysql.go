package storage

import (
	"context"
	"fmt"
	"github.com/mjedari/health-checker/app/config"
	"github.com/mjedari/health-checker/infra/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL(conf config.MySQL) (*MySQL, error) {
	ctx := context.TODO()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Pass, conf.Host, conf.Port, conf.Database)
	fmt.Println("THIISSS ISS: ", dsn)
	mysqlRetry, err := utils.Retry(func(ctx context.Context) (any, error) {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			return nil, err
		}

		return db, nil
	}, utils.RetryTimes, utils.RetryDelay)(ctx)

	if err != nil {
		return nil, err
	}
	// here we convert interface datatype to gorm.DB
	db := mysqlRetry.(*gorm.DB)

	return &MySQL{db: db}, nil
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
