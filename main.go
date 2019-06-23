package main

import (
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "goStore/handlers"
  "goStore/lib"
)

// TODO try https://github.com/appleboy/gin-lambda or just gin

func main() {
  lib.Middlewares.Request = append(lib.Middlewares.Request, cleanRequest)
  lambda.Start(handlers.Handler)
}

func cleanRequest(req *events.APIGatewayProxyRequest) {
  pathLen := len(req.Path)

  // strip trailing '/'
  if req.Path[pathLen-1] == '/' && len(req.Path) > 1 {
    req.Path = req.Path[:pathLen-1]
  }
}
