package constant

// PortmapStateType 端口映射状态类型
type PortmapStateType string

const (
	PORTMAP_STATE_STOP_WAIT  PortmapStateType = "stopwait"
	PORTMAP_STATE_STOP       PortmapStateType = "stop"
	PORTMAP_STATE_START_WAIT PortmapStateType = "startwait"
	PORTMAP_STATE_START      PortmapStateType = "start"
)

// Value 原值
func (state PortmapStateType) Value() string {
	return string(state)
}

// IsVaild 是否有效
func (state PortmapStateType) IsVaild() bool {
	return state.ToString() != "未知"
}

// ToString 转字符串
func (state PortmapStateType) ToString() (str string) {
	if state == PORTMAP_STATE_STOP {
		str = "停止"
	} else if state == PORTMAP_STATE_STOP_WAIT {
		str = "等待停止"
	} else if state == PORTMAP_STATE_START {
		str = "运行中"
	} else if state == PORTMAP_STATE_START_WAIT {
		str = "等待运行"
	} else {
		str = "未知"
	}

	return
}
