package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/zaaksam/dproxy/go/config"
	"github.com/zaaksam/dproxy/go/db"
	"github.com/zaaksam/dproxy/go/model"
)

// Info 普通级别日志
func Info(v ...interface{}) {
	insert("Info", v...)
}

// Debug 调试级别日志
func Debug(v ...interface{}) {
	insert("Debug", v...)
}

// Error 错误级别日志
func Error(v ...interface{}) {
	insert("Error", v...)
}

// Warning 警告级别日志
func Warning(v ...interface{}) {
	insert("Warning", v...)
}

// Critical 严重级别日志
func Critical(v ...interface{}) {
	insert("Critical", v...)
}

func insert(typ string, v ...interface{}) {
	vLen := len(v)
	if vLen == 0 {
		return
	}

	content := fmt.Sprintf(strings.Repeat(" %v", vLen), v...)

	if config.AppConf.Region != "" {
		content = config.AppConf.Region + "：" + content
	}

	md := &model.LogModel{
		Type:    typ,
		Content: content,
		Created: time.Now().Unix(),
	}

	da := db.NewDA()
	defer da.Close()

	_, err := da.Insert(md)
	if err != nil {
		j, _ := json.Marshal(md)
		logs.Critical(fmt.Sprintf("日志 [%s] 写入log日志错误：", string(j)), err)
	}
}
