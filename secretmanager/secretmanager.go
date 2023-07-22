package secretmanager

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/nahuelojea/handballscore/awsgo"
	"github.com/nahuelojea/handballscore/models"
)

func GetSecret(secretName string) (models.Secret, error) {
	var secretData models.Secret
	fmt.Println("> Getting Secret " + secretName)

	svc := secretsmanager.NewFromConfig(awsgo.Cfg)
	key, err := svc.GetSecretValue(awsgo.Ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		fmt.Println(err.Error())
		return secretData, err
	}

	json.Unmarshal([]byte(*key.SecretString), &secretData)
	fmt.Println(" > Secret readed OK " + secretName)
	return secretData, nil
}
