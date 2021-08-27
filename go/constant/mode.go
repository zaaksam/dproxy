package constant

// ModeType 模式类型
type ModeType string

const (
	MODE_APIUI ModeType = "apiui"
	MODE_API   ModeType = "api"
	MODE_WEB   ModeType = "web"
	MODE_APP   ModeType = "app"
)

// Value 原值
func (mode ModeType) Value() string {
	return string(mode)
}

// IsVaild 是否有效
func (mode ModeType) IsVaild() bool {
	return mode.ToString() != "未知"
}

// ToString 转字符串
func (mode ModeType) ToString() (str string) {
	if mode == MODE_APIUI {
		str = "API带UI"
	} else if mode == MODE_API {
		str = "API"
	} else if mode == MODE_WEB {
		str = "Web"
	} else if mode == MODE_APP {
		str = "App"
	} else {
		str = "未知"
	}

	return
}
