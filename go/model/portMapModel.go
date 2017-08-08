package model

// PortMapModel 端口映射模型
type PortMapModel struct {
	ID         int64  `xorm:"pk" json:"id"`
	Title      string `json:"title"`
	TargetIP   string `json:"targetIP"`
	TargetPort int    `json:"targetPort"`
	SourceIP   string `json:"sourceIP"`
	SourcePort int    `json:"sourcePort"`
	UserID     string `json:"userID"`
	UserName   string `json:"userName"`
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
	Deleted    int64  `xorm:"->" json:"-"`
	IsStart    bool   `xorm:"-" json:"isStart"`
}

// TableName 表名
func (*PortMapModel) TableName() string {
	return "portMap"
}

// NewItems 实现 db.IModel 接口
func (*PortMapModel) NewItems() interface{} {
	items := new([]*PortMapModel)
	*items = make([]*PortMapModel, 0)
	return items
}
