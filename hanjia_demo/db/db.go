package db

import (
	"hanjia_demo/model"
	"hanjia_demo/utils_log"   // 导入Zap日志包
	"hanjia_demo/utils_viper" // 导入Viper配置包
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化MySQL数据库
func Init() {
	// 1. 从Viper配置中读取MySQL参数
	cfg := utils_viper.GetConfig()
	if cfg == nil {
		utils_log.Fatal("初始化MySQL失败：配置未加载")
	}
	mysqlCfg := cfg.Mysql

	// 2. 配置GORM日志
	var gormLogger logger.Interface
	if utils_viper.IsDev() {
		// 开发环境：显示所有SQL日志
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		// 生产环境：关闭SQL日志，减少开销
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// 3. 连接MySQL（使用配置文件的DSN）
	db, err := gorm.Open(mysql.Open(mysqlCfg.Dsn), &gorm.Config{
		Logger: gormLogger, // 绑定自定义日志
	})
	if err != nil {
		utils_log.Fatal("MySQL连接失败", zap.Error(err)) // 替换log.Fatal为Zap日志
	}

	// 4. 配置连接池（从配置文件读取参数，提升性能）
	sqlDB, err := db.DB()
	if err != nil {
		utils_log.Fatal("获取SQL DB实例失败", zap.Error(err))
	}
	// 设置连接池参数
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)                                    // 最大打开连接数
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)                                    // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlCfg.ConnMaxLifetime) * time.Second) // 连接最大存活时间

	// 5. 表迁移（保留原有逻辑）
	err = db.AutoMigrate(&model.User{}, &model.Article{}, &model.Comment{}, &model.Follow{})
	if err != nil {
		utils_log.Fatal("数据表迁移失败", zap.Error(err))
	}

	// 6. 赋值全局DB并打印成功日志
	DB = db
	utils_log.Info("MySQL数据库初始化完成 ✅",
		zap.String("dsn", mysqlCfg.Dsn),
		zap.Int("max_open_conns", mysqlCfg.MaxOpenConns),
		zap.Int("max_idle_conns", mysqlCfg.MaxIdleConns),
	)
}
