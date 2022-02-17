package awssecretmanager

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretsData struct {
	DbHost             string `json:"DB_HOST"`
	DbPassword         string `json:"DB_PASSWORD"`
	DbName             string `json:"DB_NAME"`
	CognitoUserPoolId  string `json:"COGNITO_USER_POOL_ID"`
	CognitoAppClientId string `json:"COGNITO_APP_CLIENT_ID"`
}

func GetSecret(secretName, region, versionStage string) SecretsData {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials("AKIAW4SF47I55BEXZTEB", "U+qiw2GlMK+mH7Kom4sE/iRniOpwftASDq2ktplM", ""),
	})
	svc := secretsmanager.New(
		sess,
		aws.NewConfig().WithRegion(region),
	)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName + "_configuration"),
		VersionStage: aws.String(versionStage),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		fmt.Printf("Error in getting response from secret manager, %s", err)
	}

	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	}

	var secretData SecretsData
	err = json.Unmarshal([]byte(secretString), &secretData)
	if err != nil {
		fmt.Printf("Error unmarshal in secret data, %s", err)
	}

	return secretData
}
