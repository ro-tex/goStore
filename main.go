package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"

	"goStore/handlers"
	"goStore/local"
	"goStore/middlewares"
)

func main() {

	// TODO: check if API Gateway's error is errorstring or an actual *error
	//
	// TODO: Switch to gin-gonic
	//
	// TODO: Add a local execution method in order to test without deploying to AWS
	// 			https://djhworld.github.io/post/2018/01/27/running-go-aws-lambda-functions-locally/
	//
	// TODO:
	// 	    Ideas:
	// 	    	* try https://github.com/appleboy/gin-lambda or just gin
	// 	    Features:
	// 	    	* auth0
	// 			* swagger
	// 			* automated execution of middlewares?
	// 				* error class that can output nice JSON errors
	// 				* decent logging

	if runtime.GOOS == "linux" && os.Getenv("_HANDLER") != "" { // we're running on AWS Lambda

		middlewares.RegisterRequestMW(middlewares.CleanRequest)
		middlewares.RegisterResponseMW(middlewares.GnuTerryPratchett)
		middlewares.RegisterResponseMW(logResponse)
		middlewares.RegisterErrorMW(logError)

		var api = new(handlers.API)
		api.AddHandler("/v0/doc", handlers.V0DocHandler)

		lambda.Start(api.MasterHandler)

	} else { // we're running locally

		fmt.Println("WE are local")

		g := gin.Default()

		/**
		One way to define handlers is directly here. It's a bit messy but it's fine for simple handlers.
		*/
		g.GET("/about", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"about": "Send a GET request to `/add/:a/:b` and you'll get their sum."})
		})

		g.GET("/v0/doc", func(ctx *gin.Context) {
			req, err := local.TransGin2AwsReq(ctx)
			if err != nil {
				// TODO Change this to an InternalServerError with GIn
				log.Fatalf("Failed to parse request: %v", err)
			}
			res, err := handlers.V0DocHandler(req)
			if err != nil {
				// TODO Change this to an InternalServerError with GIn
				log.Fatalf("Handler returned an error: %v", err)
			}

			r := local.TransAwsRes2Gin(&res)

			ctx.JSON(r.StatusCode, r.Response)
		})

		// /**
		// Another way is to have a generator function that returns a func(ctx *gin.Context).
		// This allows us to move the handler implementation somewhere else and have this part nice and clean.
		// */
		// g.GET("/add/:a/:b", addHandler(addService))
		// g.GET("/mult/:a/:b", multHandler(addService))

		if err := g.Run(":3000"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}

	}
}

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
