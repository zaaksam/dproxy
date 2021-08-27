package services

import (
	"errors"
	"sync"
	"time"

	_ "github.com/zaaksam/dproxy/go/config"
	"github.com/zaaksam/dproxy/go/db"
	"github.com/zaaksam/dproxy/go/logger"
	"github.com/zaaksam/dproxy/go/model"
)

// WhiteList 白名单服务对象
var WhiteList whiteListService

func init() {
	WhiteList.cache = make(map[string]struct{})

	//监控白名单
	go WhiteList.Watch()
}

type whiteListService struct {
	mx               sync.RWMutex
	cache            map[string]struct{}
	cacheLastCreated int64
}

// Get 获取指定的白名单信息
func (*whiteListService) Get(id int64) (md *model.WhiteListModel, err error) {
	if id <= 0 {
		return nil, errors.New("ID不能为空")
	}

	da := db.NewDA()
	defer da.Close()

	da.Where("Deleted=0 and ID=?", id)

	// if !isExpired {
	// 	session.And("Expired>?", time.Now().Unix())
	// }

	var has bool
	md = &model.WhiteListModel{}
	has, err = da.Get(md)
	if err != nil {
		return nil, errors.New("查找白名单数据出错：" + err.Error())
	}

	if !has {
		return nil, errors.New("没有找到对应的白名单数据")
	}

	return
}

// Find 查询白名单列表信息
func (*whiteListService) Find(pageIndex, pageSize int, ip, userName, isExpired, sortField, sortDesc string) (list *model.ListModel, err error) {
	da := db.NewDA()
	defer da.Close()

	da.Where("Deleted=0")

	if ip != "" {
		da.And("IP like ?", "%"+ip+"%")
	}

	if userName != "" {
		da.And("UserName like ?", "%"+userName+"%")
	}

	if isExpired == "1" {
		da.And("Expired<?", time.Now().Unix())
	} else if isExpired == "0" {
		da.And("Expired>?", time.Now().Unix())
	}

	order := "Created"
	if sortField == "expired" {
		order = "Expired"
	}

	if sortDesc == "0" {
		order += " asc"
	} else {
		order += " desc"
	}

	da.OrderBy(order)

	list, err = da.GetList(&model.WhiteListModel{}, pageIndex, pageSize)
	return
}

// Add 添加新的白名单
func (*whiteListService) Add(md *model.WhiteListModel) (result *model.WhiteListModel, err error) {
	t := time.Now()
	created := t.Unix()

	if md == nil {
		err = errors.New("model对象不能为空")
	} else if md.IP == "" {
		err = errors.New("IP不能为空")
	} else if md.UserID == "" {
		err = errors.New("UserID不能为空")
	} else if md.UserName == "" {
		err = errors.New("UserName不能为空")
	} else if md.Expired > 0 && md.Expired <= created {
		err = errors.New("Expired不能小于当前时间")
	}
	if err != nil {
		return
	}

	da := db.NewDA()
	defer da.Close()

	var cnt int64
	cnt, err = da.Where("Deleted=0 and IP=? and Expired>?", md.IP, created).Count(&model.WhiteListModel{})
	if err != nil {
		err = errors.New("查询IP错误：" + err.Error())
	} else if cnt > 0 {
		err = errors.New("有效白名单已存在")
	}
	if err != nil {
		return
	}

	md.ID = t.UnixNano() / 1000 //微秒单位
	md.Created = created
	md.Updated = md.Created
	if md.Expired <= 0 {
		md.Expired = t.Add(24 * time.Hour).Unix()
	}
	md.Deleted = 0

	var n int64
	n, err = da.Insert(md)
	if err != nil {
		err = errors.New("添加白名单失败：" + err.Error())
		return
	}

	if n <= 0 {
		md = nil
		err = errors.New("白名单插入数据库失败")
		return
	}

	result = md
	return
}

// Update 更新白名单数据
func (*whiteListService) Update(md *model.WhiteListModel) (err error) {
	updated := time.Now().Unix()

	if md == nil {
		err = errors.New("model对象不能为空")
	} else if md.ID <= 0 {
		err = errors.New("ID不能为空")
	} else if md.IP == "" {
		err = errors.New("IP不能为空")
	} else if md.UserID == "" {
		err = errors.New("UserID不能为空")
	} else if md.UserName == "" {
		err = errors.New("UserName不能为空")
	} else if md.Expired <= updated {
		err = errors.New("Expired不能小于当前时间")
	}
	if err != nil {
		return
	}

	da := db.NewDA()
	defer da.Close()

	var cnt int64
	cnt, err = da.Where("Deleted=0 and IP=? and Expired>? and ID!=?", md.IP, updated, md.ID).Count(&model.WhiteListModel{})
	if err != nil {
		err = errors.New("查询IP错误：" + err.Error())
	} else if cnt > 0 {
		err = errors.New("有效白名单已存在")
	}
	if err != nil {
		return
	}

	md.Updated = updated

	var n int64
	n, err = da.Where("Deleted=0 and ID=?", md.ID).Cols("IP", "UserID", "UserName", "Expired", "Updated").Update(md)
	if err != nil {
		err = errors.New("更新白名单失败：" + err.Error())
	} else if n <= 0 {
		err = errors.New("没有符合条件的白名单记录，更新失败")
	}

	return
}

// Delete 删除白名单数据
func (s *whiteListService) Delete(id int64) error {
	md, err := s.Get(id)
	if err != nil {
		return err
	}
	invalidIP := md.IP

	md = &model.WhiteListModel{
		Deleted: time.Now().Unix(),
	}

	da := db.NewDA()
	defer da.Close()

	n, err := da.Where("Deleted=0 and ID=?", id).Cols("Deleted").Update(md)
	if err != nil {
		return errors.New("白名单记录删除失败：" + err.Error())
	} else if n == 0 {
		return errors.New("白名单记录删除失败：没有符合条件的记录")
	}

	//将代理中客户端设置为无效
	Proxy.setClientInvalid(invalidIP)

	return nil
}

// Clear 清理过期
func (s *whiteListService) Clear() error {
	da := db.NewDA()
	defer da.Close()

	da.Where("Expired<?", time.Now().Unix())

	md := &model.WhiteListModel{}
	n, err := da.Delete(md)
	if err != nil {
		return errors.New("清理过期白名单失败：" + err.Error())
	} else if n == 0 {
		return errors.New("清理过期白名单失败：没有符合条件的记录")
	}

	return nil
}

// Check 从缓存中检查IP是否在有效白名单之中
func (s *whiteListService) Check(ip string) bool {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if _, ok := s.cache[ip]; ok {
		return true
	}

	return false
}

// updateCache 更新有效白名单到缓存中
func (s *whiteListService) updateCache(items []*model.WhiteListModel) {
	itemsLen := len(items)
	var lastCreated int64
	if itemsLen > 0 {
		lastCreated = items[0].Created
	}

	//生成新的缓存内容
	newCache := make(map[string]struct{})
	for i := 0; i < itemsLen; i++ {
		newCache[items[i].IP] = struct{}{}
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	if itemsLen == len(s.cache) && lastCreated == s.cacheLastCreated {
		//内容没有变化
		return
	}

	//对比新老缓存差异，在新缓存名单中不存在的，为过期白名单，设置为无效
	for ip := range s.cache {
		if _, ok := newCache[ip]; !ok {
			Proxy.setClientInvalid(ip)
		}
	}

	s.cache = newCache
	s.cacheLastCreated = lastCreated
}

// Watch 监控有效白名单，并更新到缓存中，2 秒循环
func (s *whiteListService) Watch() {
	for {
		list, err := s.Find(1, 500, "", "", "0", "", "")
		if err != nil {
			logger.Error("白名单监控失败：", err)
		} else {
			if items, ok := list.Items.(*[]*model.WhiteListModel); ok {
				s.updateCache(*items)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
