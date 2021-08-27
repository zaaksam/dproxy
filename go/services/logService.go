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
func (*logService) Find(pageIndex, pageSize int, created int64, typ, content string) (list *model.ListModel, err error) {
	da := db.NewDA()
	defer da.Close()

	if created > 0 {
		da.And("Created<?", created)
	}

	if typ != "" {
		da.And("Type=?", typ)
	}

	if content != "" {
		da.And("Content like ?", "%"+content+"%")
	}

	da.Desc("Created")

	list, err = da.GetList(&model.LogModel{}, pageIndex, pageSize)
	return
}

// Delete 删除端口映射数据
func (*logService) Delete(created int64, typ, content string) error {
	if created <= 0 {
		return errors.New("created不能为空")
	}

	da := db.NewDA()
	defer da.Close()

	da.Where("Created<?", created)
	if typ != "" {
		da.And("type=?", typ)
	}
	if content != "" {
		da.And("content=?", "%"+content+"%")
	}

	md := &model.LogModel{}
	n, err := da.Delete(md)
	if err != nil {
		return errors.New("日志记录删除失败：" + err.Error())
	} else if n == 0 {
		return errors.New("日志记录删除失败：没有符合条件的记录")
	}

	return nil
}
