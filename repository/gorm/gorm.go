package gorm

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

//var DB *gorm.DB

type ConnConfig struct {
	Dsn         string
	MaxIdleConn time.Duration
	MaxOpenConn int
	MaxLifeTime time.Duration
	Logger      logger.Interface
}

func New(cfg *ConnConfig) (*gorm.DB, error) {
	DB, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       cfg.Dsn, // DSN data source name
		DefaultStringSize:         256,     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger:         cfg.Logger,
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		cfg.Logger.Error(context.Background(), "connect MySQL failed", err)
		return nil, err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		cfg.Logger.Error(context.Background(), "get sqlDB failed", err)
		return nil, err
	}

	sqlDB.SetConnMaxIdleTime(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime)

	cfg.Logger.Info(context.Background(), "init done")
	return DB, nil
}
