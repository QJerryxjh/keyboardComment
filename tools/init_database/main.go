package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/Unknwon/goconfig"
)

// 数据库信息模型
type InitDataType struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

// 管理员数据模型
type DbManager struct {
	gorm.Model
	Account  string `gorm:"type:varchar(16);not null;unique"`
	Password string `gorm:"type:varchar(32);not null"`
	Name     string `gorm:"type:varchar(16);not null"`
	Gender   string `gorm:"type:varchar(16)"`
}

// 获取当前运行程序的路径
func GetPath() (path string, err error) {
	return os.Getwd()
}

// 加密信息
func cryptoMd5(str string) string {
	md5Init := md5.New()
	md5Init.Write([]byte(str))
	mdSum := md5Init.Sum(nil)
	// 返回转换成16进制的数据
	return hex.EncodeToString(mdSum)
}

// 配置数据库，有初始化配置文件使用配置文件的配置，否则添加默认配置，并创建配置文件
func Config() (err error) {
	// 获取运行路径
	filePath, err := GetPath()
	if err != nil {
		return err
	}
	initFilePath := filePath + "/conf.ini"

	createdNewFile := false

	// 查看文件信息
	_, err = os.Stat(initFilePath)
	if err != nil {
		// 没有文件就创
		file, err1 := os.Create("conf.ini")
		if err1 != nil {
			return err1
		}
		err2 := file.Close()
		if err2 != nil {
			return err2
		}
		createdNewFile = true
	}
	config, err := goconfig.LoadConfigFile(initFilePath)
	if err != nil {
		return
	}

	rootSection := "root"
	proSection := "pro"

	if createdNewFile {
		// 新创建的配置文件，添加配置值
		// 如果原来就有配置，就会跳过
		config.SetValue(rootSection, "database", "mysql")
		config.SetValue(rootSection, "username", "root")
		config.SetValue(rootSection, "password", "12345678")
		config.SetValue(rootSection, "host", "127.0.0.1")
		config.SetValue(rootSection, "port", "3306")
		config.SetValue(rootSection, "charset", "utf8")
		config.SetValue(rootSection, "parseTime", "True")
		config.SetValue(rootSection, "loc", "Local")

		config.SetValue(proSection, "database", "keyboardComment")
		config.SetValue(proSection, "username", "admin")
		config.SetValue(proSection, "password", "12345678")
		config.SetValue(proSection, "host", "127.0.0.1")
		config.SetValue(proSection, "port", "3306")
		config.SetValue(proSection, "charset", "utf8")
		config.SetValue(proSection, "parseTime", "True")
		config.SetValue(proSection, "loc", "Local")
		err = goconfig.SaveConfigFile(config, "conf.ini")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	rootConf, err := config.GetSection(rootSection)
	if err != nil {
		fmt.Println(err)
		return
	}
	proConf, err := config.GetSection(proSection)
	if err != nil {
		return
	}

	rootPassword := rootConf["password"]

	var initDatabase InitDataType
	initDatabase.Database = proConf["database"]
	initDatabase.Username = proConf["username"]
	initDatabase.Password = proConf["password"]
	initDatabase.Port = proConf["port"]
	initDatabase.Host = proConf["host"]

	err = initDatabase.InitDatabase(rootPassword)
	if err != nil {
		return
	}

	return
}

func (self *InitDataType) InitDatabase(rootPassword string) (err error) {
	db, err := gorm.Open("mysql", "root:"+rootPassword+"@(127.0.0.1)/mysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("开启数据库失败：", err)
		return
	}

	// 如果不存在就创建条目
	err = db.Exec("CREATE USER if not exists '" + self.Username + "'@'%' IDENTIFIED BY'" + self.Password + "';").Error
	defer func() {
		if err != nil {
			// 过程出错。删除条目
			err = db.Exec("drop user '" + self.Username + "'@'%';").Error
			err = db.Exec("drop user '" + self.Username + "'@'localhost';").Error
		}
	}()
	if err != nil {
		return
	}

	// 修改密码
	err = db.Exec("alter user " + self.Username + "@'%' identified with mysql_native_password by '" + self.Password + "';").Error
	if err != nil {
		return
	}

	// 创建数据库
	err = db.Exec("create database if not exists " + self.Database + " charset utf8;").Error
	defer func() {
		if err != nil {
			// 有错误删除数据库
			err = db.Exec("drop database " + self.Database + ";").Error
		}
	}()
	if err != nil {
		return
	}

	// 修改数据库权限，给所有权限
	err = db.Exec("grant all privileges on " + self.Database + ".* to '" + self.Username + "'@'%';").Error
	if err != nil {
		return
	}

	// 刷新
	err = db.Exec("flush privileges;").Error
	if err != nil {
		return
	}

	// 切换数据库
	err = db.Exec("use " + self.Database + ";").Error
	if err != nil {
		return
	}

	// 查看有无管理员表，没有则创建管理员表
	if !db.HasTable(&DbManager{}) {
		err = db.CreateTable(&DbManager{}).Error
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	kcAdmin := DbManager{}
	var count int
	db.Where("account = ?", "kcAdmin").First(&kcAdmin).Count(&count)
	if count > 0 {
		fmt.Println("已经存在管理员")
		// 存在初始管理员
		return
	}

	// 初始化一个管理员
	var newKcAdmin DbManager
	password := "123123"
	securePassword := cryptoMd5(password)
	newKcAdmin.Account = "kcAdmin"
	newKcAdmin.Password = securePassword
	newKcAdmin.Name = "kcAdmin"
	newKcAdmin.Gender = "god is a girl"
	err = db.Create(&newKcAdmin).Error
	if err != nil {
		return
	}
	return
}

func main() {
	err := Config()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("初始化数据库成功")
	}
}
