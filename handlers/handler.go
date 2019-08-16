package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"goStore/middlewares"
)

/*
TODO: https://gitlab.com/ro-tex/gostore/issues/2
  This should be the main handler function that takes care of auth, executing middlewares, etc.
  The actual handling of requests should be done in separate, per-endpoint handlers.
  Those other handlers should return to this one, so it can execute response middlewares and also check for errors,
  so it can run the error middlewares. Middlewares rule everything! :D
*/
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ver := os.Getenv("VERSION")
	sha := os.Getenv("SHA")
	if len(ver) > 0 || len(sha) > 0 {
		fmt.Printf("Version: %s, SHA: %s\n", ver, sha)
	}

	// Execute request middlewares:
	for _, m := range middlewares.GetRequestMWs() {
		m(&req)
	}

	var res events.APIGatewayProxyResponse
	var err error

	// Delegate requests to different endpoints to different handlers:
	if strings.HasPrefix(req.Path, "/v0/doc") {
		res, err = v0DocHandler(req)
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Endpoint not found: " + req.Path,
		}, nil
	}

	// Execute error middlewares:
	if err != nil {
		for _, m := range middlewares.GetErrorMWs() {
			m(&err)
		}
	}

	// Execute response middlewares:
	for _, m := range middlewares.GetResponseMWs() {
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
