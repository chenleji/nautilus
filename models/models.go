package models

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/chenleji/nautilus/helper"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DefaultDB  = "default"
	DriverName = "mysql"
)

func init() {
	orm.Debug = false

	conf := helper.GetConfInst().LoadConf()

	dbURL := conf.DBUrl
	dbPort := conf.DBPort
	dbUser := conf.DBUser
	dbPassword := conf.DBPwd
	dbName := helper.Utils{}.GetAppName()

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPassword, dbURL, dbPort, dbName)
	maxIdle := 30
	maxConn := 50
	logs.Debug("dataSource:", dataSource)

	orm.RegisterDriver(DriverName, orm.DRMySQL)
	err := orm.RegisterDataBase(DefaultDB, DriverName, dataSource, maxIdle, maxConn)
	if err != nil {
		logs.Debug(err)
	}

	orm.RegisterModel(new(Model))
	orm.RunSyncdb(DefaultDB, false, true)
}

type Model struct {
	Id int
}

func DatabaseHealth() bool {
	if _, err := orm.NewOrm().Raw("SHOW TABLES", ).Exec(); err != nil {
		logs.Error(err)
		return false
	}

	return true
}
