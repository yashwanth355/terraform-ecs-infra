package Config

import (
	"fmt"
	SecretManagerService "usermvc/service/awssecretmanager"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

const (
	APP_ENV     = "APP_ENV"
	PRODUCTION  = "production"
	DEVELOPMENT = "development"
	QA          = "qa"
	configDir   = "../config"
)

var (
	configuration *Configuration
)

type SecretsData struct {
	DbHost             string `json:"DB_HOST"`
	DbPassword         string `json:"DB_PASSWORD"`
	DbName             string `json:"DB_NAME"`
	CognitoUserPoolId  string `json:"COGNITO_USER_POOL_ID"`
	CognitoAppClientId string `json:"COGNITO_APP_CLIENT_ID"`
}

type Configuration struct {
	AppEnv             string `mapstructure:"APP_ENV"`
	DbPort             int    `mapstructure:"DB_PORT"`
	DbName             string `mapstructure:"DB_NAME"`
	DbUserName         string `mapstructure:"DB_USERNAME"`
	DbPassword         string `mapstructure:"DB_PASSWORD"`
	DbHost             string `mapstructure:"DB_HOST"`
	AwsCognitoRegion   string `mapstructure:"AWS_COGNITO_REGION"`
	CognitoUserPoolId  string `mapstructure:"COGNITO_USER_POOL_ID"`
	CognitoAppClientId string `mapstructure:"COGNITO_APP_CLIENT_ID"`
	FileName           string `mapstructure:"FILE_NAME"`
	MaxSize            int    `mapstructure:"MAX_SIZE"`
	MaxAge             int    `mapstructure:"MAX_AGE"`
	MaxBackups         int    `mapstructure:"MAX_BACKUPS"`
	LocalTime          bool   `mapstructure:"LOCAL_TIME"`
	Compress           bool   `mapstructure:"COMPRESS"`
}

func LoadConfig() *Configuration {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	env := viper.GetString(APP_ENV)
	if env == DEVELOPMENT || env == QA || env == PRODUCTION {
		secretResponse := SecretManagerService.GetSecret(env, configuration.AwsCognitoRegion, "AWSCURRENT")
		configuration.DbHost = secretResponse.DbHost
		configuration.DbName = secretResponse.DbName
		configuration.DbPassword = secretResponse.DbPassword
		configuration.CognitoAppClientId = secretResponse.CognitoAppClientId
		configuration.CognitoUserPoolId = secretResponse.CognitoUserPoolId
	}

	fmt.Println(configuration)
	return configuration
}
