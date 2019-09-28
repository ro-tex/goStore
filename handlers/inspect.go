package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

// Inspect is a debug handler that just responds with the request in the body
func Inspect(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 555,
			Body:       err.Error(),
		}, nil
	}

	fmt.Println("Request: " + string(jsonReq))

	return events.APIGatewayProxyResponse{
		StatusCode: 222,
		Body:       string(jsonReq),
	}, nil
}
