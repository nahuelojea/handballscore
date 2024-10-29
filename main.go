package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nahuelojea/handballscore/api/handlers"
	"github.com/nahuelojea/handballscore/config/awsgo"
	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/config/secretmanager"
	"github.com/nahuelojea/handballscore/dto"
)

const (
	APP_DOMAINS = "https://handballscore.onrender.com, https://www.handballscore.com"
)

var (
	secretModel     dto.Secret
	loadSecretOnce  sync.Once
	secretLoadError error
)

func main() {
	lambda.Start(executeLambda)
}

func executeLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
		"Access-Control-Allow-Methods": "OPTIONS, POST, GET, PUT, PATCH, DELETE",
		"Content-Type":                 "application/json",
	}

	if request.HTTPMethod == "OPTIONS" {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    headers,
		}, nil
	}

	awsgo.Init()

	if !ValidEnvironmentVariables() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error to get environment variables. Must include 'SecretName', 'BucketName', 'UrlPrefix'",
			Headers:    headers,
		}
		return res, nil
	}

	loadSecretOnce.Do(func() {
		secretModel, secretLoadError = secretmanager.GetSecret(os.Getenv("SecretName"))
	})

	if secretLoadError != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error to read Secret " + secretLoadError.Error(),
			Headers:    headers,
		}
		return res, nil
	}

	path := strings.Replace(request.PathParameters["handballscore"], os.Getenv("UrlPrefix"), "", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("user"), secretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("password"), secretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("host"), secretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("database"), secretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("jwtSign"), secretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("bucketName"), os.Getenv("BucketName"))

	err := db.Connect(awsgo.Ctx)

	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error to connect with database " + err.Error(),
			Headers:    headers,
		}
		return res, nil
	}

	respAPI := handlers.ProcessRequest(awsgo.Ctx, request)

	fmt.Println("API Response: " + respAPI.Message)

	if respAPI.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       respAPI.Message,
			Headers:    headers,
		}
		return res, nil
	} else {
		return respAPI.CustomResp, nil
	}
}

func ValidEnvironmentVariables() bool {
	_, parameter := os.LookupEnv("SecretName")

	if !parameter {
		return false
	}

	_, parameter = os.LookupEnv("BucketName")

	if !parameter {
		return false
	}

	_, parameter = os.LookupEnv("UrlPrefix")

	return parameter
}