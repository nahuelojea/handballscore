package dto

import "github.com/aws/aws-lambda-go/events"

type RestResponse struct {
	Status     int
	Message    string
	CustomResp *events.APIGatewayProxyResponse
}
