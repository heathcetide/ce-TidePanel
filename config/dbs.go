package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitDB() (db *gorm.DB, err error) {
	InitConfig()
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)
	db, err = gorm.Open(driverName, args)
	if err != nil {
		return nil, err
	}
	log.Println("database connect succcess")
	db.LogMode(true)
	return db, err
}

func GetDB() (db *gorm.DB, err error) {
	return InitDB()
}

func MakeMigrates(db *gorm.DB, insts []any) error {
	var errors []error
	for _, v := range insts {
		ptr, ok := v.(interface{ new() any })
		if !ok {
			errors = append(errors, fmt.Errorf("invalid type %T for migration", v))
			continue
		}
		if err := db.AutoMigrate(ptr); err != nil {
			errors = append(errors)
			log.Printf("Error migrating %T: %v", v, err)
		} else {
			log.Printf("Successfully migrated %T", v)
		}
	}
	if len(errors) > 0 {
		return errors[0] // 或者使用 multierror 包来返回所有错误
	}
	return nil
}

func InitConfig() {
	rootDir, err := os.Getwd()
	viper.AddConfigPath(rootDir) // 设置配置文件路径为当前目录
	viper.SetConfigName(".env")  // 设置配置文件名为 .env
	viper.SetConfigType("properties")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
