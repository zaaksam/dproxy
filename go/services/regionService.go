package services

import (
	"errors"
	"time"

	"github.com/zaaksam/dproxy/go/config"
	"github.com/zaaksam/dproxy/go/constant"
	"github.com/zaaksam/dproxy/go/db"
	"github.com/zaaksam/dproxy/go/logger"
	"github.com/zaaksam/dproxy/go/model"
)

// Region 服务
var Region regionService

type regionService struct{}

func init() {
	// 注册区域信息到数据库，并更新心跳时间，重置心跳超时的区域端口映射数据状态为 stop
	go Region.Watch()
}

// Add 添加新的端口映射
func (s *regionService) Add(md *model.RegionModel) (result *model.RegionModel, err error) {
	if md == nil {
		err = errors.New("model对象不能为空")
		// } else if md.Name == "" {
		// 	err = errors.New("Name不能为空")
	} else if md.LastHeartbeat <= 0 {
		err = errors.New("LastHearbeat不能为空")
	}
	if err != nil {
		return
	}

	da := db.NewDA()
	defer da.Close()

	var n int64
	n, err = da.Insert(md)
	if err != nil {
		err = errors.New("添加区域失败：" + err.Error())
		return
	}

	if n <= 0 {
		err = errors.New("区域插入数据库失败")
		return
	}

	result = md
	return
}

// Update 更新区域数据
func (s *regionService) Update(md *model.RegionModel) (err error) {
	if md == nil {
		err = errors.New("model对象不能为空")
		// } else if md.Name == "" {
		// 	err = errors.New("Name不能为空")
	} else if md.LastHeartbeat <= 0 {
		err = errors.New("LastHearbeat不能为空")
	}
	if err != nil {
		return
	}

	da := db.NewDA()
	defer da.Close()

	var n int64
	n, err = da.Where("Name=?", md.Name).Cols("LastHeartbeat").Update(md)
	if err != nil {
		err = errors.New("更新区域失败：" + err.Error())
	} else if n <= 0 {
		err = errors.New("没有符合条件的区域记录，更新失败")
	}

	return
}

// CheckTimeout 检查区域是否超时
func (s *regionService) CheckTimeout(name string) (isTimeout bool, err error) {
	md, err := s.Get(name)
	if err != nil {
		return
	}

	// 当前减去 30 秒，作为健康在线的范围
	ux := time.Now().Unix() - 30
	if md.LastHeartbeat < ux {
		isTimeout = true
	}
	return
}

// Get 获取
func (s *regionService) Get(name string) (md *model.RegionModel, err error) {
	var has bool
	md, has, err = s.get(name)
	if err != nil {
		return
	}

	if !has {
		err = errors.New("没有找到对应的区域数据")
	}
	return
}

func (s *regionService) get(name string) (md *model.RegionModel, has bool, err error) {
	// if name == "" {
	// 	err = errors.New("Name不能为空")
	// 	return
	// }

	da := db.NewDA()
	defer da.Close()

	da.Where("Name=?", name)

	md = &model.RegionModel{}
	has, err = da.Get(md)
	if err != nil {
		err = errors.New("查找区域数据出错：" + err.Error())
	}

	return
}

// Find 查询区域列表信息
func (s *regionService) Find(pageIndex, pageSize int) (list *model.ListModel, err error) {
	da := db.NewDA()
	defer da.Close()

	da.Asc("Name")
	list, err = da.GetList(&model.RegionModel{}, pageIndex, pageSize)
	return
}

// 查询心跳超时的列表数据
func (s *regionService) findTimeout(pageIndex, pageSize int) (list *model.ListModel, err error) {
	da := db.NewDA()
	defer da.Close()

	// 当前减去 30 秒，作为健康在线的范围
	ux := time.Now().Unix() - 30

	// 100  109-10=99
	// 110  115-10=105
	// 120  120-10=110
	// 120  121-10=111
	// 120  130-10=120
	// 120  131-10=121
	da.Where("LastHeartbeat<?", ux)
	da.Asc("LastHeartbeat")

	list, err = da.GetList(&model.RegionModel{}, pageIndex, pageSize)
	return
}

// Watch 监控心跳过期的区域数据，重置该区域的 portmap state 为 stop
func (s *regionService) Watch() {
	for {
		// 更新心跳时间
		md, has, err := s.get(config.AppConf.Region)
		if err != nil {
			logger.Error("更新心跳时间出错：" + err.Error())
		} else {
			if has {
				md.LastHeartbeat = time.Now().Unix()
				err = s.Update(md)
			} else {
				md = &model.RegionModel{
					Name: config.AppConf.Region,
				}
				md.LastHeartbeat = time.Now().Unix()
				_, err = s.Add(md)
			}
			if err != nil {
				logger.Error("更新心跳时间出错：" + err.Error())
			}
		}

		// 检查过期
		list, err := Region.findTimeout(1, 500)
		if err != nil {
			logger.Error("获取心跳超时的区域列表失败：", err)
		} else {
			if items, ok := list.Items.(*[]*model.RegionModel); ok {
				items2 := *items

				for i, l := 0, len(items2); i < l; i++ {
					err = PortMap.resetState(constant.PORTMAP_STATE_STOP, items2[i].Name)
					if err != nil {
						logger.Error("重置心跳超时的 "+items2[i].Name+" 区域端口映射列表状态失败：", err)
					}
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}
