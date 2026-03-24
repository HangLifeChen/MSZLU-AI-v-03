package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// BasicSettings 基本设置
type BasicSettings struct {
	SystemName          string `json:"systemName"`
	SystemDescription   string `json:"systemDescription"`
	Language            string `json:"language"`
	Theme               string `json:"theme"`
	EnableNotifications bool   `json:"enableNotifications"`
}

// ModelSettings 模型设置
type ModelSettings struct {
	DefaultProvider string  `json:"defaultProvider"`
	DefaultModel    string  `json:"defaultModel"`
	Temperature     float64 `json:"temperature"`
	MaxTokens       int     `json:"maxTokens"`
	TopP            float64 `json:"topP"`
}

// SecuritySettings 安全设置
type SecuritySettings struct {
	PasswordPolicy   string `json:"passwordPolicy"`
	SessionTimeout   int    `json:"sessionTimeout"`
	Enable2FA        bool   `json:"enable2FA"`
	MaxLoginAttempts int    `json:"maxLoginAttempts"`
}

// NotificationSettings 通知设置
type NotificationSettings struct {
	EmailEnabled     bool   `json:"emailEnabled"`
	SMSEnabled       bool   `json:"smsEnabled"`
	InAppEnabled     bool   `json:"inAppEnabled"`
	SystemTemplate   string `json:"systemTemplate"`
	SecurityTemplate string `json:"securityTemplate"`
}

// AliyunOSSSettings 阿里云 OSS 设置
type AliyunOSSSettings struct {
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Bucket          string `json:"bucket"`
	Endpoint        string `json:"endpoint"`
	PathPrefix      string `json:"pathPrefix"`
}

// QiniuSettings 七牛云设置
type QiniuSettings struct {
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey"`
	Bucket     string `json:"bucket"`
	Zone       string `json:"zone"`
	PathPrefix string `json:"pathPrefix"`
	Domain     string `json:"domain"`
}

// StorageSettings 云存储设置
type StorageSettings struct {
	DefaultProvider string            `json:"defaultProvider"`
	Aliyun          AliyunOSSSettings `json:"aliyun"`
	Qiniu           QiniuSettings     `json:"qiniu"`
}

// SystemSettings 系统设置
type SystemSettings struct {
	BaseModel
	Basic        BasicSettings        `json:"basic" gorm:"column:basic;type:jsonb"`
	Model        ModelSettings        `json:"model" gorm:"column:model;type:jsonb"`
	Security     SecuritySettings     `json:"security" gorm:"column:security;type:jsonb"`
	Notification NotificationSettings `json:"notification" gorm:"column:notification;type:jsonb"`
	Storage      StorageSettings      `json:"storage" gorm:"column:storage;type:jsonb"`
}

// TableName 返回表名
func (SystemSettings) TableName() string {
	return "system_settings"
}

// Value 实现 driver.Valuer 接口
func (b BasicSettings) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan 实现 sql.Scanner 接口
func (b *BasicSettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("无法扫描为 []byte 类型")
	}
	return json.Unmarshal(bytes, b)
}

// Value 实现 driver.Valuer 接口
func (m ModelSettings) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan 实现 sql.Scanner 接口
func (m *ModelSettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("无法扫描为 []byte 类型")
	}
	return json.Unmarshal(bytes, m)
}

// Value 实现 driver.Valuer 接口
func (s SecuritySettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口
func (s *SecuritySettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("无法扫描为 []byte 类型")
	}
	return json.Unmarshal(bytes, s)
}

// Value 实现 driver.Valuer 接口
func (n NotificationSettings) Value() (driver.Value, error) {
	return json.Marshal(n)
}

// Scan 实现 sql.Scanner 接口
func (n *NotificationSettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("无法扫描为 []byte 类型")
	}
	return json.Unmarshal(bytes, n)
}

// Value 实现 driver.Valuer 接口
func (s StorageSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口
func (s *StorageSettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("无法扫描为 []byte 类型")
	}
	return json.Unmarshal(bytes, s)
}
