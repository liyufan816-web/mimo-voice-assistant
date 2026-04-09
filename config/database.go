package config

import (
	"fmt"
	"log"
	"os"

	"github.com/liyufan816-web/mimo-voice-assistant/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
	Loc      string
}

// LoadDatabaseConfig 加载数据库配置
func LoadDatabaseConfig() *DatabaseConfig {
	godotenv.Load()
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "mimo_tts"),
		Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		Loc:      getEnv("DB_LOC", "Local"),
	}
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
		c.Loc,
	)
}

// InitDB 初始化数据库连接
func InitDB(cfg *DatabaseConfig) error {
	var err error
	DB, err = gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}

	// 获取底层 SQL 数据库对象进行配置
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库对象失败：%v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 自动迁移数据表
	err = AutoMigrate()
	if err != nil {
		return fmt.Errorf("自动迁移数据表失败：%v", err)
	}

	log.Println("数据库连接成功")
	return nil
}

// AutoMigrate 自动迁移数据表
func AutoMigrate() error {
	err := DB.AutoMigrate(
		&model.User{},
		&model.HistoryItem{},
		&model.TTSRecord{},
	)
	if err != nil {
		return fmt.Errorf("自动迁移失败：%v", err)
	}
	log.Println("数据表创建/迁移成功")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}


// SetLoc 设置 GORM 的时区
func SetLoc(loc string) {
	os.Setenv("DB_LOC", loc)
}
