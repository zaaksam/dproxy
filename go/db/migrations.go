package db

import (
	"fmt"

	"github.com/zaaksam/dproxy/go/model"
	"xorm.io/xorm"
)

// Migration 升迁接口
type Migration interface {
	Version() string
	Description() string
	Migrate(*xorm.Engine) error
}

type migration struct {
	version     string
	description string
	fn          func(*xorm.Engine) error
}

// Version 版本号
func (m *migration) Version() string {
	return m.version
}

// Description 描述
func (m *migration) Description() string {
	return m.description
}

// Migrate 升迁
func (m *migration) Migrate(x *xorm.Engine) error {
	return m.fn(x)
}

func newMigration(ver, desc string, fn func(*xorm.Engine) error) Migration {
	return &migration{version: ver, description: desc, fn: fn}
}

func checkMigration() (err error) {
	verMD := &model.Version{
		ID: 1,
	}

	err = Engine.Sync2(verMD)
	if err != nil {
		return
	}

	has, err := Engine.Get(verMD)
	if err == nil && !has {
		verMD.ID = 0
		verMD.Version = "0.0.0"

		_, err = Engine.InsertOne(verMD)
	}
	if err != nil {
		return
	}
	currentVersion := verMD.Version

	var pass bool
	for i, l := 0, len(migrations); i < l; i++ {
		pass, err = VersionComparse(migrations[i].Version(), VERSION_COMPARE_COND_GT, currentVersion)
		if err != nil {
			err = fmt.Errorf("%s 版本数据库升迁： %s ，比较版本号时发生错误：%s", migrations[i].Version(), migrations[i].Description(), err.Error())
			return
		}

		if !pass {
			// 要升迁的版本号，小于当前版本号，跳过
			continue
		}

		err = migrations[i].Migrate(Engine)
		if err != nil {
			err = fmt.Errorf("%s 版本数据库升迁： %s ，处理时发生错误：%s", migrations[i].Version(), migrations[i].Description(), err.Error())
			return
		}

		verMD.Version = migrations[i].Version()
		_, err = Engine.ID(1).Update(verMD)
		if err != nil {
			err = fmt.Errorf("%s 版本数据库升迁： %s ，更新版本号时发生错误：%s", migrations[i].Version(), migrations[i].Description(), err.Error())
			return
		}

		fmt.Println(fmt.Sprintf("%s 版本数据库升迁： %s", migrations[i].Version(), migrations[i].Description()))
	}

	return
}

var migrations = []Migration{
	v0_6_0(),
}
