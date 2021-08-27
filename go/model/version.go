package model

// Version 版本
type Version struct {
	ID      int64  `json:"id" xorm:"pk autoincr bigint(20) not null 'id'"`
	Version string `json:"version" xorm:"varchar(20) not null 'version' comment('当前版本号')"`
}

// TableName 表别名
func (*Version) TableName() string {
	return "version"
}

// NewItems 实现 db.IModel 接口
func (*Version) NewItems() interface{} {
	items := new([]*Version)
	*items = make([]*Version, 0)
	return items
}
