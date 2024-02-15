package database

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"message/config"
	"message/logs"
	"strconv"
	"strings"
	"time"
)

// DB 全局变量 DB 用于存储数据库连接实例
var DB *gorm.DB

// 进行连接重试
var maxRetries = 5
var curRetries = 1

// InitMySQL 用于初始化 MySQL 数据库连接
func InitMySQL() {
	// 从配置中获取数据库连接所需的信息
	username := config.AppConfig.Database.User
	password := config.AppConfig.Database.Pwd
	host := config.AppConfig.Database.Host
	port := config.AppConfig.Database.Port
	database := config.AppConfig.Database.Name
	charset := config.AppConfig.Database.Params.Character

	// 构建 DSN 字符串
	dsn := strings.Join([]string{
		username,
		":",
		password,
		"@tcp(",
		host,
		":",
		strconv.Itoa(port),
		")/",
		database,
		"?charset" + charset + "&parseTime=True",
	}, "")

	// 进行数据库连接
	err := databaseConnect(dsn)
	if err != nil {
		logs.LogError.Errorf(
			"InitMySQL-数据库连接失败！ %s 3秒后尝试连接。正在尝试%d次 剩余%d次。",
			err,
			curRetries,
			maxRetries-curRetries,
		)
		curRetries = maxRetries - curRetries
		if maxRetries-curRetries > 0 {
			time.Sleep(time.Second * 3)
			InitMySQL()
		}
	}
}

// databaseConnect 用于实际连接数据库
func databaseConnect(dsn string) error {
	// 根据 Gin 的模式设置 ORM 日志级别
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	// 使用 GORM 进行数据库连接
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	// 设置数据库连接池参数
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 10)
	sqlDB.SetConnMaxIdleTime(time.Second * 15)

	// 将数据库连接赋值给全局变量 DB
	DB = db

	// 初始化数据库迁移
	InitMigration()

	return err
}
