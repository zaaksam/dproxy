package model

// RegionModel 区域数据模型
type RegionModel struct {
	Name          string `json:"name"`
	LastHeartbeat int64  `json:"lastHeartbeat"` // 时间戮，单位：秒
}

// TableName 表名
func (*RegionModel) TableName() string {
	return "region"
}

// NewItems 实现 db.IModel 接口
func (*RegionModel) NewItems() interface{} {
	items := new([]*RegionModel)
	*items = make([]*RegionModel, 0)
	return items
}
