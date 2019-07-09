package handlers

import (
	"encoding/json"
	"fmt"
	"goStore/middlewares"

	"github.com/aws/aws-lambda-go/events"
)

/*
TODO:
  This should be the main handler function that takes care of auth, executing middlewares, etc.
  The actual handling of requests should be done in separate, per-endpoint handlers.
  Those other handlers should return to this one, so it can execute response middlewares and also check for errors,
  so it can run the error middlewares. Middlewares rule eveything! :D
*/
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Execute request middlewares:
	for _, m := range middlewares.GlobalMWs.GetRequestMWs() {
		m(&req)
	}

	// Delegate requests to different endpoints to different handlers:
	// if path starts with v0/doc...
	res, err := v0DocHandler(req)

	// Execute error middlewares:
	if err != nil {
		for _, m := range middlewares.GlobalMWs.GetErrorMWs() {
			m(&err)
		}
	}

	// Execute response middlewares:
	for _, m := range middlewares.GlobalMWs.GetResponseMWs() {
		m(&res)
	}

	return res, err
}

func inspect(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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
