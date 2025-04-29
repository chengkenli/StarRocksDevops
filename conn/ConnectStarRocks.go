package conn

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"starrocks/util"
	"time"
)

func StarRocks(app string) (*gorm.DB, error) {
	var avg util.ConnectParms
	// connect
	if len(util.MetaLink) == 0 {
		return nil, errors.New("config db is null")
	}
	for _, m := range util.MetaLink {
		if m["app"].(string) == app {
			avg = util.ConnectParms{
				Host: m["feip"].(string),
				Port: int(m["feport"].(int32)),
				User: m["user"].(string),
				Pass: m["password"].(string),
			}
		}
	}

	if avg.Host == "" {
		return nil, errors.New("avg is null")
	}

	newLogger := logger.New(nil,
		logger.Config{
			SlowThreshold: time.Second * 1000, // 控制慢SQL阈值
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?parseTime=true&charset=utf8mb4&loc=Local",
		avg.User,
		avg.Pass,
		avg.Host,
		avg.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("Error >:%s", err.Error()))
		return nil, err
	}
	return db, err
}

func StarRocksSingle(app, host string) (*gorm.DB, error) {
	var avg util.ConnectParms
	if len(host) == 0 {
		return nil, errors.New("登录信息有误，或为空。")
	}
	for _, m := range util.MetaLink {
		if m["app"].(string) == app {
			avg = util.ConnectParms{
				Host: host,
				Port: int(m["feport"].(int32)),
				User: m["user"].(string),
				Pass: m["password"].(string),
			}
		}
	}

	if avg.Host == "" {
		return nil, errors.New("avg is null")
	}

	newLogger := logger.New(nil,
		logger.Config{
			SlowThreshold: time.Second * 1000, // 控制慢SQL阈值
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?charset=utf8mb4&parseTime=True&loc=Local", avg.User, avg.Pass, host, avg.Port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return db, err
}
