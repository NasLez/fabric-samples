package dao

import (
	"context"
	"fabric_ebl/sal/config"
	"fabric_ebl/sal/dao/generator/query"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(ctx context.Context, config *config.AppConfig) {
	db, err := gorm.Open(mysql.Open(config.Mysql.Dsn), &gorm.Config{})
	if err != nil {
		logger.CtxErrorf(ctx, "[Init] init db failed, err = %v", err)
		panic(err)
	}
	query.SetDefault(db)
	logger.CtxInfof(ctx, "[Init] init db success")
}
