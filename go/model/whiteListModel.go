package model

// WhiteListModel 白名单模型
type WhiteListModel struct {
	ID       int64  `xorm:"pk" json:"id"`
	IP       string `json:"ip"`
	UserID   string `json:"userID"`
	UserName string `json:"userName"`
	Expired  int64  `json:"expired"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
	Deleted  int64  `xorm:"->" json:"-"`
}

// TableName 表名
func (*WhiteListModel) TableName() string {
	return "whiteList"
}

// NewItems 实现 db.IModel 接口
func (*WhiteListModel) NewItems() interface{} {
	items := new([]*WhiteListModel)
	*items = make([]*WhiteListModel, 0)
	return items
}
