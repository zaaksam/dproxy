package services

import (
	"errors"

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

// Delete 删除端口映射数据
func (*logService) Delete(typ string, created int64, content string) error {
	if created <= 0 {
		return errors.New("created不能为空")
	}

	session := db.Engine.Where("Created<?", created)
	if typ != "" {
		session.And("type=?", typ)
	}
	if content != "" {
		session.And("content=?", "%"+content+"%")
	}

	md := &model.LogModel{}
	n, err := session.Delete(md)
	if err != nil {
		return errors.New("日志记录删除失败：" + err.Error())
	} else if n == 0 {
		return errors.New("日志记录删除失败：没有符合条件的记录")
	}

	return nil
}
