package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nahuelojea/handballscore/awsgo"
)

func main() {
	lambda.Start(executeLambda)
}

func executeLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.Init()

	if !ValidateParameters() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error to get environment variables. Must include 'SecretName', 'BucketName', 'UrlPrefix'",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}
	return res, nil
}

func ValidateParameters() bool {
	_, parameter := os.LookupEnv("SecretName")

	if !parameter {
		return false
	}

	_, parameter = os.LookupEnv("BucketName")

	if !parameter {
		return false
	}

	_, parameter = os.LookupEnv("UrlPrefix")

	if !parameter {
		return false
	}
	return true
}
