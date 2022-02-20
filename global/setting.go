package global

import (
	"blogservice/pkg/logger"
	"blogservice/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	JWTSetting      *setting.JWTSettinS
	Logger          *logger.Logger
	EmailSetting    *setting.EmailSettingS
)
