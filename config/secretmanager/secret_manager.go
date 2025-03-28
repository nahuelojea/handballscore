package secretmanager

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/nahuelojea/handballscore/config/awsgo"
	"github.com/nahuelojea/handballscore/dto"
)

func GetSecret(secretName string) (dto.Secret, error) {
	fmt.Println("Getting secret from AWS Secret Manager")
	var secretData dto.Secret

	svc := secretsmanager.NewFromConfig(awsgo.Cfg)
	key, err := svc.GetSecretValue(awsgo.Ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		fmt.Println(err.Error())
		return secretData, err
	}

	json.Unmarshal([]byte(*key.SecretString), &secretData)
	return secretData, nil
}
