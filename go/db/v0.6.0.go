package db

import (
	"github.com/zaaksam/dproxy/go/model"
	"xorm.io/xorm"
)

func v0_6_0() Migration {
	return newMigration(
		"0.6.0",
		"portmap 表添加 Region、State 字段；",
		func(x *xorm.Engine) (err error) {
			portMapMD := &model.PortMapModel{}

			sqlRaw := "alter table `" + portMapMD.TableName() + "` add column `Region` varchar(50) not null comment '区域' after `ID`;"
			sqlRaw += "alter table `" + portMapMD.TableName() + "` add column `State` varchar(50) not null comment '状态' after `Deleted`;"

			_, err = x.Exec(sqlRaw)
			return
		},
	)
}
