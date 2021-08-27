package db

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zaaksam/dproxy/go/config"
	"github.com/zaaksam/dproxy/go/model"

	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

// Engine 数据库引擎
var Engine *xorm.Engine
var tables []names.TableName

// init 初始化
func init() {
	tables = make([]names.TableName, 20)

	var (
		err            error
		driverName     string
		dataSourceName string
	)
	if config.AppConf.MysqlServer == "" {
		dir := path.Dir(os.Args[0])
		dir = filepath.ToSlash(dir) // .../dproxy/go
		dataSourceName = filepath.Join(dir, "dproxy.db")

		_, err = os.Stat(dataSourceName)
		if err != nil && os.IsNotExist(err) {
			_, err = os.Create(dataSourceName)
			if err != nil {
				panic("创建DB文件错误：" + err.Error())
			}
		}

		driverName = "sqlite3"
	} else {
		driverName = "mysql"
		dataSourceName = config.AppConf.MysqlUsername + ":" + config.AppConf.MysqlPassword
		dataSourceName += "@(" + config.AppConf.MysqlServer + ":" + strconv.Itoa(config.AppConf.MysqlPort) + ")/"
		dataSourceName += config.AppConf.MysqlDatabase + "?charset=utf8mb4&loc=Asia%2FShanghai&multiStatements=true"
	}

	Engine, err = xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		panic("创建数据库引擎错误：" + err.Error())
	}

	err = Engine.Ping()
	if err != nil {
		panic("DB连接不通：" + err.Error())
	}

	//结构体命名与数据库一致
	Engine.SetMapper(core.NewCacheMapper(new(core.SameMapper)))

	err = createTable(
		&model.LogModel{},
		&model.WhiteListModel{},
		&model.PortMapModel{},
		&model.RegionModel{},
	)
	if err != nil {
		if driverName == "sqlite3" {
			os.Remove(dataSourceName)
		}
		panic(err)
	}

	err = checkMigration()
	if err != nil {
		panic(err)
	}

	return
}

func createTable(beans ...names.TableName) (err error) {
	// tables, err := Engine.DBMetas()
	// if err != nil {
	// 	panic(fmt.Sprintf("获取DB信息错误：%s", err))
	// }

	var ok bool
	for i, l := 0, len(beans); i < l; i++ {

		ok, err = Engine.IsTableExist(beans[i])
		if err != nil {
			err = fmt.Errorf("检查数据库表[%s]是否存在失败：%s", beans[i].TableName(), err)
			return
		}

		if ok {
			// 表存在，跳过
			continue
		}

		// err = Engine.CreateTables(beans[i])
		err = Engine.Sync2(beans[i])
		if err != nil {
			err = fmt.Errorf("创建数据库表[%s]失败：%s", beans[i].TableName(), err)
			return
		}
	}
	return
}
