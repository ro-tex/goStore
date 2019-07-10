package middlewares

import "github.com/aws/aws-lambda-go/events"

/*
TODO
  OK, obviously a global object is not the perfect approach. I need middlewares per handler and all that. For that the
  handler needs to be a constructable object with a Handle() function exposed and that function should be passed to the
  lambda.Start(). But until then this should do.
*/

var globalMWs = Middlewares{}

type Middlewares struct {
	requestMWs  []func(req *events.APIGatewayProxyRequest)
	responseMWs []func(req *events.APIGatewayProxyResponse)
	errorMWs    []func(req *error)
}

func RegisterRequestMW(fn func(req *events.APIGatewayProxyRequest)) {
	globalMWs.requestMWs = append(globalMWs.requestMWs, fn)
}

func RegisterResponseMW(fn func(req *events.APIGatewayProxyResponse)) {
	globalMWs.responseMWs = append(globalMWs.responseMWs, fn)
}

func RegisterErrorMW(fn func(req *error)) {
	globalMWs.errorMWs = append(globalMWs.errorMWs, fn)
}

func GetRequestMWs() []func(req *events.APIGatewayProxyRequest) {
	return globalMWs.requestMWs
}

func GetResponseMWs() []func(req *events.APIGatewayProxyResponse) {
	return globalMWs.responseMWs
}

func GetErrorMWs() []func(req *error) {
	return globalMWs.errorMWs
}

// CleanRequest strips the trailing '/' if it exists
func CleanRequest(req *events.APIGatewayProxyRequest) {
	pathLen := len(req.Path)
	if req.Path[pathLen-1] == '/' && len(req.Path) > 1 {
		req.Path = req.Path[:pathLen-1]
	}
}
