package services

import (
	"github.com/zaaksam/dproxy/go/db"
	"github.com/zaaksam/dproxy/go/model"
)

// Log 日志服务对象
var Log logService

type logService struct{}

// Find 查询白名单列表信息
func (*logService) Find(pageIndex, pageSize int, typ, content string) (list *model.ListModel, err error) {
	session := db.NewSession()

	if typ != "" {
		session.And("Type=?", typ)
	}

	if content != "" {
		session.And("Content like ?", "'%"+content+"%'")
	}

	session.Desc("Created")

	list, err = db.GetList(session, &model.LogModel{}, pageIndex, pageSize)
	return
}
