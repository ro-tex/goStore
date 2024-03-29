package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"goStore/handlers"
	"goStore/middlewares"
)

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

func multiLog(event interface{}) {
	var err error
	defer func() {
		if err != nil {
			fmt.Println("[Error][multiLog]:", err.Error())
		}
	}()

	switch reflect.TypeOf(event).Name() {
	case "*events.APIGatewayProxyRequest":
		r := event.(*events.APIGatewayProxyRequest)
		j, err := json.Marshal(r)
		if err != nil {
			return
		}
		fmt.Println("[Request] " + string(j))
	case "*events.APIGatewayProxyResponse":
		res := event.(*events.APIGatewayProxyResponse)
		j, err := json.Marshal(res)
		if err != nil {
			return
		}
		fmt.Println("[Response] " + string(j))
	case "*errors.errorString":
		r := event.(*error)
		j, err := json.Marshal(r)
		if err != nil {
			return
		}
		fmt.Println("[Error] " + string(j))
	default:
		fmt.Println("WARN: Trying to log unknown object of type " + reflect.TypeOf(event).Name())
	}
}

func main() {
	/*
		TODO: check if API Gateway's error is errorstring or an actual *error

		TODO: Switch to gin-gonic

		TODO: Add a local execution method in order to test without deploying to AWS
					https://djhworld.github.io/post/2018/01/27/running-go-aws-lambda-functions-locally/

		TODO:
			    Ideas:
			    	* try https://github.com/appleboy/gin-lambda or just gin
			    Features:
			    	* auth0
					* swagger
					* automated execution of middlewares?
						* error class that can output nice JSON errors
						* decent logging
	*/

	godotenv.Load() // load .env variables

	middlewares.RegisterRequestMW(middlewares.CleanRequest)
	middlewares.RegisterResponseMW(middlewares.GnuTerryPratchett)
	middlewares.RegisterResponseMW(logResponse)
	middlewares.RegisterErrorMW(logError)

	var api = new(handlers.API)
	api.AddHandler("/v0/doc", handlers.V0DocHandler)

	lambda.Start(api.MasterHandler)
}
