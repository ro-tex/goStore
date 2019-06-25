package middlewares

import "github.com/aws/aws-lambda-go/events"

/*
TODO
  OK, obviously a global object is not the perfect approach. I need middlewares per handler and all that. For that the
  handler needs to be a constructable object with a Handle() function exposed and that function should be passed to the
  lambda.Start(). But until then this should do.
*/

var GlobalMWs = Middlewares{}

type Middlewares struct {
  requestMWs  []func(req *events.APIGatewayProxyRequest)
  responseMWs []func(req *events.APIGatewayProxyResponse)
  errorMWs    []func(req *error)
}

func (mws *Middlewares) RegisterRequestMW(fn func(req *events.APIGatewayProxyRequest)) {
  mws.requestMWs = append(mws.requestMWs, fn)
}

func (mws *Middlewares) RegisterResponseMW(fn func(req *events.APIGatewayProxyResponse)) {
  mws.responseMWs = append(mws.responseMWs, fn)
}

func (mws *Middlewares) RegisterErrorMW(fn func(req *error)) {
  mws.errorMWs = append(mws.errorMWs, fn)
}

func (mws *Middlewares) GetRequestMWs() []func(req *events.APIGatewayProxyRequest) {
  return mws.requestMWs
}

func (mws *Middlewares) GetResponseMWs() []func(req *events.APIGatewayProxyResponse) {
  return mws.responseMWs
}

func (mws *Middlewares) GetErrorMWs() []func(req *error) {
  return mws.errorMWs
}

// CleanRequest strips the trailing '/' if it exists
func CleanRequest(req *events.APIGatewayProxyRequest) {
  pathLen := len(req.Path)
  if req.Path[pathLen-1] == '/' && len(req.Path) > 1 {
    req.Path = req.Path[:pathLen-1]
  }
}
