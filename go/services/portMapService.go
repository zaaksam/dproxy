package services

import (
	"errors"
	"time"

	"github.com/zaaksam/dproxy/go/db"
	"github.com/zaaksam/dproxy/go/model"
)

// PortMap 端口映射服务对象
var PortMap portMapService

type portMapService struct{}

// Get 获取指定的端口映射信息
func (s *portMapService) Get(id int64) (md *model.PortMapModel, err error) {
	if id <= 0 {
		return nil, errors.New("ID不能为空")
	}

	session := db.NewSession()
	session.Where("Deleted=0 and ID=?", id)

	var has bool
	md = &model.PortMapModel{}
	has, err = session.Get(md)
	if err != nil {
		return nil, errors.New("查找端口映射数据出错：" + err.Error())
	}

	if !has {
		return nil, errors.New("没有找到对应的端口映射数据")
	}

	s.checkIsStart(md)

	return
}

// Find 查询端口映射列表信息
func (s *portMapService) Find(pageIndex, pageSize int, targetIP, sourceIP, userName string) (list *model.ListModel, err error) {
	session := db.NewSession()
	session.Where("Deleted=0")

	if targetIP != "" {
		session.And("TargetIP=?", targetIP)
	}

	if sourceIP != "" {
		session.And("SourceIP=?", sourceIP)
	}

	if userName != "" {
		session.And("UserName=?", userName)
	}

	session.Desc("Created")

	list, err = db.GetList(session, &model.PortMapModel{}, pageIndex, pageSize)

	if mds, ok := list.Items.(*[]*model.PortMapModel); ok {
		s.checkIsStart(*mds...)
	}

	return
}

// Add 添加新的端口映射
func (*portMapService) Add(md *model.PortMapModel) (result *model.PortMapModel, err error) {
	if md == nil {
		err = errors.New("model对象不能为空")
	} else if md.Title == "" {
		err = errors.New("Title不能为空")
	} else if md.TargetIP == "" {
		err = errors.New("TargetIP不能为空")
	} else if md.TargetPort <= 0 {
		err = errors.New("TargetPort不能为空")
	} else if md.SourceIP == "" {
		err = errors.New("SourceIP不能为空")
	} else if md.SourcePort <= 0 {
		err = errors.New("SourcePort不能为空")
	} else if md.UserID == "" {
		err = errors.New("UserID不能为空")
	} else if md.UserName == "" {
		err = errors.New("UserName不能为空")
	}
	if err != nil {
		return
	}

	session := db.NewSession()

	var cnt int64
	cnt, err = session.Where("Deleted=0 and TargetIP=? and TargetPort=? and SourceIP=? and SourcePort=?", md.TargetIP, md.TargetPort, md.SourceIP, md.SourcePort).Count(&model.PortMapModel{})
	if err != nil {
		err = errors.New("查询有效端口映射数据错误：" + err.Error())
	} else if cnt > 0 {
		err = errors.New("有效端口映射已存在")
	}
	if err != nil {
		return
	}

	t := time.Now()
	md.ID = t.UnixNano() / 1000 //微秒单位
	md.Created = t.Unix()
	md.Updated = md.Created
	md.Deleted = 0

	var n int64
	n, err = session.Insert(md)
	if err != nil {
		err = errors.New("添加端口映射失败：" + err.Error())
		return
	}

	if n <= 0 {
		err = errors.New("端口映射插入数据库失败")
		return
	}

	result = md
	return
}

// Update 更新端口映射数据
func (*portMapService) Update(md *model.PortMapModel) (err error) {
	if md == nil {
		err = errors.New("model对象不能为空")
	} else if md.ID <= 0 {
		err = errors.New("ID不能为空")
	} else if md.Title == "" {
		err = errors.New("Title不能为空")
	} else if md.TargetIP == "" {
		err = errors.New("TargetIP不能为空")
	} else if md.TargetPort <= 0 {
		err = errors.New("TargetPort不能为空")
	} else if md.SourceIP == "" {
		err = errors.New("SourceIP不能为空")
	} else if md.SourcePort <= 0 {
		err = errors.New("SourcePort不能为空")
	} else if md.UserID == "" {
		err = errors.New("UserID不能为空")
	} else if md.UserName == "" {
		err = errors.New("UserName不能为空")
	}
	if err != nil {
		return
	}

	session := db.NewSession()

	var cnt int64
	cnt, err = session.Where("Deleted=0 and TargetIP=? and TargetPort=? and SourceIP=? and SourcePort=? and ID != ?", md.TargetIP, md.TargetPort, md.SourceIP, md.SourcePort, md.ID).Count(&model.PortMapModel{})
	if err != nil {
		err = errors.New("查询有效端口映射数据错误：" + err.Error())
	} else if cnt > 0 {
		err = errors.New("有效端口映射已存在")
	}
	if err != nil {
		return
	}

	md.Updated = time.Now().Unix()

	var n int64
	n, err = db.Engine.Where("Deleted=0 and ID=?", md.ID).Cols("Title", "TargetIP", "TargetPort", "SourceIP", "SourcePort", "UserID", "UserName", "Updated").Update(md)
	if err != nil {
		err = errors.New("更新端口映射失败：" + err.Error())
	} else if n <= 0 {
		err = errors.New("没有符合条件的端口映射记录，更新失败")
	}

	return
}

// Delete 删除端口映射数据
func (*portMapService) Delete(id int64) error {
	if id <= 0 {
		return errors.New("ID不能为空")
	}

	md := &model.PortMapModel{
		Deleted: time.Now().Unix(),
	}

	n, err := db.Engine.Where("Deleted=0 and ID=?", id).Cols("Deleted").Update(md)
	if err != nil {
		return errors.New("端口映射记录删除失败：" + err.Error())
	} else if n == 0 {
		return errors.New("端口映射记录删除失败：没有符合条件的记录")
	}

	return nil
}

func (s *portMapService) checkIsStart(mds ...*model.PortMapModel) {
	for i, l := 0, len(mds); i < l; i++ {
		if _, ok := Proxy.proxys[mds[i].ID]; ok {
			mds[i].IsStart = true
		}
	}
}
