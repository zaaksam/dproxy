package db

import (
	"errors"
	"strconv"
	"strings"
)

type VersionCompareCond string

const (
	VERSION_COMPARE_COND_EQ  VersionCompareCond = "eq"  // 等于
	VERSION_COMPARE_COND_GT  VersionCompareCond = "gt"  // 大于
	VERSION_COMPARE_COND_GTE VersionCompareCond = "gte" // 大于或等于
	VERSION_COMPARE_COND_LT  VersionCompareCond = "lt"  // 小于
	VERSION_COMPARE_COND_LTE VersionCompareCond = "lte" // 小于或等于
)

// Value 原值
func (cond VersionCompareCond) Value() string {
	return string(cond)
}

// IsVaild 是否有效
func (cond VersionCompareCond) IsVaild() bool {
	return cond.ToString() != "未知"
}

// ToString 转字符串
func (cond VersionCompareCond) ToString() (str string) {
	if cond == VERSION_COMPARE_COND_EQ {
		str = "等于"
	} else if cond == VERSION_COMPARE_COND_GT {
		str = "大于"
	} else if cond == VERSION_COMPARE_COND_GTE {
		str = "大于或等于"
	} else if cond == VERSION_COMPARE_COND_LT {
		str = "小于"
	} else if cond == VERSION_COMPARE_COND_LTE {
		str = "小于或等于"
	} else {
		str = "未知"
	}

	return
}

// VersionComparse 版本号比较
func VersionComparse(leftVer string, cond VersionCompareCond, rightVer string) (pass bool, err error) {
	leftVer = strings.ReplaceAll(strings.ToLower(strings.TrimSpace(leftVer)), "v", "")
	leftVers := strings.Split(leftVer, ".")

	rightVer = strings.ReplaceAll(strings.ToLower(strings.TrimSpace(rightVer)), "v", "")
	rightVers := strings.Split(rightVer, ".")

	l := len(leftVers)

	if l != len(rightVers) {
		err = errors.New("版本号格式不一致")
		return
	}

	pad := "00000"
	padLen := len(pad)
	padLenStr := strconv.Itoa(padLen)
	leftVerFormat := ""
	rightVerFormat := ""

	for i := 0; i < l; i++ {
		if len(leftVers[i]) > padLen {
			err = errors.New("leftVer 位数长度超出：" + padLenStr)
			return
		}

		if len(rightVers[i]) > padLen {
			err = errors.New("rightVer 位数长度超出：" + padLenStr)
			return
		}

		leftVerFormat += pad[:padLen-len(leftVers[i])] + leftVers[i]
		rightVerFormat += pad[:padLen-len(rightVers[i])] + rightVers[i]
	}

	if cond == VERSION_COMPARE_COND_EQ && leftVerFormat == rightVerFormat {
		pass = true
	} else if cond == VERSION_COMPARE_COND_GT && leftVerFormat > rightVerFormat {
		pass = true
	} else if cond == VERSION_COMPARE_COND_GTE && leftVerFormat >= rightVerFormat {
		pass = true
	} else if cond == VERSION_COMPARE_COND_LT && leftVerFormat < rightVerFormat {
		pass = true
	} else if cond == VERSION_COMPARE_COND_LTE && leftVerFormat <= rightVerFormat {
		pass = true
	}

	return
}
