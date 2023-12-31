package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nahuelojea/handballscore/api/handlers"
	"github.com/nahuelojea/handballscore/config/awsgo"
	"github.com/nahuelojea/handballscore/config/db"
	"github.com/nahuelojea/handballscore/config/secretmanager"
	"github.com/nahuelojea/handballscore/dto"
)

func main() {
	lambda.Start(executeLambda)
}

func executeLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.Init()

	if !ValidEnvironmentVariables() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error to get environment variables. Must include 'SecretName', 'BucketName', 'UrlPrefix'",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))

	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error to read Secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	path := strings.Replace(request.PathParameters["handballscore"], os.Getenv("UrlPrefix"), "", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, dto.Key("bucketName"), os.Getenv("BucketName"))

	err = db.Connect(awsgo.Ctx)

	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error to connect with database " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	respAPI := handlers.ProcessRequest(awsgo.Ctx, request)
	if respAPI.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       respAPI.Message,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
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

	if !parameter {
		return false
	}
	return true
}
