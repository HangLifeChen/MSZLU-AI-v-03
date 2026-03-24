package settings

import (
	"context"
	"model"
	"time"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
)

type service struct {
	repo repository
}

func (s *service) getSettings(ctx context.Context) (*model.SystemSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	settings, err := s.repo.getSettings(ctx)
	if err != nil {
		logs.Errorf("get settings error: %v", err)
		return nil, errs.DBError
	}
	return settings, nil
}

func (s *service) saveSettings(ctx context.Context, req SaveSettingsReq) (*model.SystemSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	settings := &model.SystemSettings{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		Basic:        req.Basic,
		Model:        req.Model,
		Security:     req.Security,
		Notification: req.Notification,
		Storage:      req.Storage,
	}
	err := s.repo.createSettings(ctx, settings)
	if err != nil {
		logs.Errorf("save settings error: %v", err)
		return nil, errs.DBError
	}
	return settings, nil
}

func (s *service) updateSettings(ctx context.Context, req UpdateSettingsReq) (*model.SystemSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	// 获取现有设置
	settings, err := s.repo.getSettings(ctx)
	if err != nil {
		logs.Errorf("get settings error: %v", err)
		return nil, errs.DBError
	}
	if settings == nil {
		// 如果不存在则创建新的
		settings = &model.SystemSettings{
			BaseModel: model.BaseModel{
				ID: uuid.New(),
			},
			Basic:        req.Basic,
			Model:        req.Model,
			Security:     req.Security,
			Notification: req.Notification,
			Storage:      req.Storage,
		}
		err = s.repo.createSettings(ctx, settings)
		if err != nil {
			logs.Errorf("create settings error: %v", err)
			return nil, errs.DBError
		}
	} else {
		// 更新现有设置
		settings.Basic = req.Basic
		settings.Model = req.Model
		settings.Security = req.Security
		settings.Notification = req.Notification
		settings.Storage = req.Storage
		err = s.repo.updateSettings(ctx, settings)
		if err != nil {
			logs.Errorf("update settings error: %v", err)
			return nil, errs.DBError
		}
	}
	return settings, nil
}

func (s *service) getSettingsModule(ctx context.Context, module string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	settings, err := s.repo.getSettings(ctx)
	if err != nil {
		logs.Errorf("get settings error: %v", err)
		return nil, errs.DBError
	}
	if settings == nil {
		return nil, nil
	}
	switch module {
	case "basic":
		return settings.Basic, nil
	case "model":
		return settings.Model, nil
	case "security":
		return settings.Security, nil
	case "notification":
		return settings.Notification, nil
	case "storage":
		return settings.Storage, nil
	default:
		return nil, errs.NewError(400, "invalid module")
	}
}

func (s *service) updateSettingsModule(ctx context.Context, module string, data map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	settings, err := s.repo.getSettings(ctx)
	if err != nil {
		logs.Errorf("get settings error: %v", err)
		return errs.DBError
	}
	if settings == nil {
		settings = &model.SystemSettings{
			BaseModel: model.BaseModel{
				ID: uuid.New(),
			},
		}
	}
	switch module {
	case "basic":
		if systemName, ok := data["systemName"]; ok {
			settings.Basic.SystemName = systemName.(string)
		}
		if systemDescription, ok := data["systemDescription"]; ok {
			settings.Basic.SystemDescription = systemDescription.(string)
		}
		if language, ok := data["language"]; ok {
			settings.Basic.Language = language.(string)
		}
		if theme, ok := data["theme"]; ok {
			settings.Basic.Theme = theme.(string)
		}
		if enableNotifications, ok := data["enableNotifications"]; ok {
			settings.Basic.EnableNotifications = enableNotifications.(bool)
		}
	case "model":
		if defaultProvider, ok := data["defaultProvider"]; ok {
			settings.Model.DefaultProvider = defaultProvider.(string)
		}
		if defaultModel, ok := data["defaultModel"]; ok {
			settings.Model.DefaultModel = defaultModel.(string)
		}
		if temperature, ok := data["temperature"]; ok {
			settings.Model.Temperature = temperature.(float64)
		}
		if maxTokens, ok := data["maxTokens"]; ok {
			settings.Model.MaxTokens = int(maxTokens.(float64))
		}
		if topP, ok := data["topP"]; ok {
			settings.Model.TopP = topP.(float64)
		}
	case "security":
		if passwordPolicy, ok := data["passwordPolicy"]; ok {
			settings.Security.PasswordPolicy = passwordPolicy.(string)
		}
		if sessionTimeout, ok := data["sessionTimeout"]; ok {
			settings.Security.SessionTimeout = int(sessionTimeout.(float64))
		}
		if enable2FA, ok := data["enable2FA"]; ok {
			settings.Security.Enable2FA = enable2FA.(bool)
		}
		if maxLoginAttempts, ok := data["maxLoginAttempts"]; ok {
			settings.Security.MaxLoginAttempts = int(maxLoginAttempts.(float64))
		}
	case "notification":
		if emailEnabled, ok := data["emailEnabled"]; ok {
			settings.Notification.EmailEnabled = emailEnabled.(bool)
		}
		if smsEnabled, ok := data["smsEnabled"]; ok {
			settings.Notification.SMSEnabled = smsEnabled.(bool)
		}
		if inAppEnabled, ok := data["inAppEnabled"]; ok {
			settings.Notification.InAppEnabled = inAppEnabled.(bool)
		}
		if systemTemplate, ok := data["systemTemplate"]; ok {
			settings.Notification.SystemTemplate = systemTemplate.(string)
		}
		if securityTemplate, ok := data["securityTemplate"]; ok {
			settings.Notification.SecurityTemplate = securityTemplate.(string)
		}
	case "storage":
		if defaultProvider, ok := data["defaultProvider"]; ok {
			settings.Storage.DefaultProvider = defaultProvider.(string)
		}
		if aliyunData, ok := data["aliyun"]; ok {
			aliyunMap := aliyunData.(map[string]interface{})
			if accessKeyId, ok := aliyunMap["accessKeyId"]; ok {
				settings.Storage.Aliyun.AccessKeyID = accessKeyId.(string)
			}
			if accessKeySecret, ok := aliyunMap["accessKeySecret"]; ok {
				settings.Storage.Aliyun.AccessKeySecret = accessKeySecret.(string)
			}
			if bucket, ok := aliyunMap["bucket"]; ok {
				settings.Storage.Aliyun.Bucket = bucket.(string)
			}
			if endpoint, ok := aliyunMap["endpoint"]; ok {
				settings.Storage.Aliyun.Endpoint = endpoint.(string)
			}
			if pathPrefix, ok := aliyunMap["pathPrefix"]; ok {
				settings.Storage.Aliyun.PathPrefix = pathPrefix.(string)
			}
		}
		if qiniuData, ok := data["qiniu"]; ok {
			qiniuMap := qiniuData.(map[string]interface{})
			if accessKey, ok := qiniuMap["accessKey"]; ok {
				settings.Storage.Qiniu.AccessKey = accessKey.(string)
			}
			if secretKey, ok := qiniuMap["secretKey"]; ok {
				settings.Storage.Qiniu.SecretKey = secretKey.(string)
			}
			if bucket, ok := qiniuMap["bucket"]; ok {
				settings.Storage.Qiniu.Bucket = bucket.(string)
			}
			if zone, ok := qiniuMap["zone"]; ok {
				settings.Storage.Qiniu.Zone = zone.(string)
			}
			if pathPrefix, ok := qiniuMap["pathPrefix"]; ok {
				settings.Storage.Qiniu.PathPrefix = pathPrefix.(string)
			}
			if domain, ok := qiniuMap["domain"]; ok {
				settings.Storage.Qiniu.Domain = domain.(string)
			}
		}
	default:
		return errs.NewError(400, "invalid module")
	}
	err = s.repo.updateSettings(ctx, settings)
	if err != nil {
		logs.Errorf("update settings error: %v", err)
		return errs.DBError
	}
	return nil
}

func newService() *service {
	return &service{
		repo: newModels(database.GetPostgresDB().GormDB),
	}
}
