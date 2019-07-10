package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"goStore/handlers"
	"goStore/middlewares"
)

/*
	TODO:
    Ideas:
    * try https://github.com/appleboy/gin-lambda or just gin
    * set up .gitlab-ci.yml
    Features:
    * auth0
    * automated execution of middlewares?
		* error class that can output nice JSON errors
		* decent logging
*/

// just a test response middleware to see if it works correctly
func logResponse(res *events.APIGatewayProxyResponse) {
	jsonReq, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("[logResponse] Response: " + string(jsonReq))
}

func logError(e *error) {
	jsonE, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("[logError] Response: " + string(jsonE))
}

func main() {
	middlewares.RegisterRequestMW(middlewares.CleanRequest)
	middlewares.RegisterResponseMW(logResponse)
	middlewares.RegisterErrorMW(logError)
	lambda.Start(handlers.Handler)
}
