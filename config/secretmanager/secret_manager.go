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
	// fmt.Println("Getting secret from AWS Secret Manager") // Removed
	var secretData dto.Secret

	svc := secretsmanager.NewFromConfig(awsgo.Cfg)
	key, err := svc.GetSecretValue(awsgo.Ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		// fmt.Println(err.Error()) // Removed
		return secretData, err
	}

	json.Unmarshal([]byte(*key.SecretString), &secretData)
	return secretData, nil
}

func GetFirebaseSecret(secretName string) (dto.FirebaseSecret, error) {
	// fmt.Println("Getting Firebase secret from AWS Secret Manager") // Removed
	var secretData dto.FirebaseSecret

	svc := secretsmanager.NewFromConfig(awsgo.Cfg)
	key, err := svc.GetSecretValue(awsgo.Ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		// fmt.Printf("Error getting secret value: %s\n", err.Error()) // Removed
		return secretData, err
	}

	if key.SecretString == nil {
		return secretData, fmt.Errorf("secret string is nil for secret: %s", secretName)
	}

	err = json.Unmarshal([]byte(*key.SecretString), &secretData)
	if err != nil {
		// fmt.Printf("Error unmarshalling secret string: %s\n", err.Error()) // Removed
		return secretData, err
	}

	return secretData, nil
}
