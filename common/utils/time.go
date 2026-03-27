package utils

import (
	"encoding/json"
	"time"
)

// CustomTime 自定义时间类型，支持多种日期格式
type CustomTime struct {
	time.Time
}

// UnmarshalJSON 实现自定义的 JSON 反序列化，支持多种日期格式
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		ct.Time = time.Time{}
		return nil
	}

	// 去除引号
	str := string(data)
	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	if str == "" {
		ct.Time = time.Time{}
		return nil
	}

	// 支持的日期格式列表
	formats := []string{
		time.RFC3339,                // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05Z07:00", // 带时区
		"2006-01-02T15:04:05",       // 不带时区
		"2006-01-02T15:04:05+08:00", // 固定时区
		"2006-01-02T15:04:05-08:00", // 固定时区
		"2006-01-02 15:04:05",       // 空格分隔
		"2006-01-02",                // 只有日期
		"2006/01/02",                // 斜杠分隔
		"2006/01/02 15:04:05",       // 斜杠分隔带时间
	}

	var err error
	for _, format := range formats {
		ct.Time, err = time.Parse(format, str)
		if err == nil {
			return nil
		}
	}

	return err
}

// MarshalJSON 实现 JSON 序列化
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ct.Time.Format(time.RFC3339))
}

// ToTime 转换为标准 time.Time
func (ct CustomTime) ToTime() time.Time {
	return ct.Time
}

// NewCustomTime 从 time.Time 创建 CustomTime
func NewCustomTime(t time.Time) CustomTime {
	return CustomTime{Time: t}
}
