package dbs

import (
	"errors"
	"os"

	"github.com/Unknwon/goconfig"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func GetPath() (path string, err error) {
	return os.Getwd()
}

func InitEnvironment(env string) (err error) {
	confFilePath, err := GetPath()
	if err != nil {
		return
	}
	flag := false
	_, err = os.Stat(confFilePath + "conf.ini")
	if err != nil {
		file, err := os.Create("conf.ini")
		if err != nil {
			return err
		}
		defer file.Close()
		flag = true
	}
	if flag {
		config, err1 := goconfig.LoadConfigFile(confFilePath + "/conf.ini")
		if err1 != nil {
			return err1
		}
		config.SetValue("root", "database", "mysql")
		config.SetValue("root", "username", "root")
		config.SetValue("root", "password", "12345678")
		config.SetValue("root", "host", "127.0.0.1")
		config.SetValue("root", "port", "3306")
		config.SetValue("root", "charset", "utf8")
		config.SetValue("root", "parseTime", "True")
		config.SetValue("root", "loc", "Local")
		config.SetValue(env, "database", "keyboardComment")
		config.SetValue(env, "username", "admin")
		config.SetValue(env, "password", "12345678")
		config.SetValue(env, "host", "127.0.0.1")
		config.SetValue(env, "port", "3306")
		config.SetValue(env, "charset", "utf8")
		config.SetValue(env, "parseTime", "True")
		config.SetValue(env, "loc", "Local")
		err = goconfig.SaveConfigFile(config, "conf.ini")
		if err != nil {
			return
		}
	}

	getConfig, err := goconfig.LoadConfigFile(confFilePath + "/conf.ini")
	if err != nil {
		return
	}

	username, err := getConfig.GetValue(env, "username")
	if err != nil {
		return
	}
	database, err := getConfig.GetValue(env, "database")
	if err != nil {
		return
	}
	password, err := getConfig.GetValue(env, "password")
	if err != nil {
		return
	}
	host, err := getConfig.GetValue(env, "host")
	if err != nil {
		return
	}
	charset, err := getConfig.GetValue(env, "charset")
	if err != nil {
		return
	}
	port, err := getConfig.GetValue(env, "port")
	if err != nil {
		return
	}
	parseTime, err := getConfig.GetValue(env, "parseTime")
	if err != nil {
		return
	}
	loc, err := getConfig.GetValue(env, "loc")
	if err != nil {
		return
	}

	if username == "" || password == "" || host == "" || charset == "" || parseTime == "" || port == "" || loc == "" || database == "" {
		return errors.New("初始化数据库缺少参数")
	}
	db, err := gorm.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return
	}
	db.DB().SetMaxIdleConns(1024)
	db.DB().SetMaxOpenConns(256)

	DB = db
	return
}

func Close() (err error) {
	err = DB.Close()
	if err != nil {
		return err
	}
	return
}
