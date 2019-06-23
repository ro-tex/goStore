package main

import (
  "github.com/aws/aws-lambda-go/lambda"

  "goStore/handlers"
  "goStore/middlewares"
)

/*
	TODO:
    Ideas:
    * try https://github.com/appleboy/gin-lambda or just gin
    Features:
    * auth0
    * automated execution of middlewares?
		* error class that can output nice JSON errors
		* decent logging
*/

func main() {
  middlewares.Request.Register(middlewares.CleanRequest)
  lambda.Start(handlers.Handler)
}
