package middlewares

import "github.com/aws/aws-lambda-go/events"

type Middleware = func(req *events.APIGatewayProxyRequest)

type MiddlewareGroup struct {
  middlewares []Middleware
}

func (mg *MiddlewareGroup) Register(fn Middleware) {
  mg.middlewares = append(mg.middlewares, fn)
}

func (mg *MiddlewareGroup) GetAll() []Middleware {
  return mg.middlewares
}

var Request MiddlewareGroup
var Response MiddlewareGroup
var Error MiddlewareGroup

// CleanRequest strips the trailing '/' if it exists
func CleanRequest(req *events.APIGatewayProxyRequest) {
  pathLen := len(req.Path)
  if req.Path[pathLen-1] == '/' && len(req.Path) > 1 {
    req.Path = req.Path[:pathLen-1]
  }
}
