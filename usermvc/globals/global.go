package globals

import (
	Config "usermvc/config"
)

var (
	appConfig *Config.Configuration
)

func GetConfig() *Config.Configuration {
	if appConfig == nil {
		return Config.LoadConfig()
	}
	return appConfig
}
