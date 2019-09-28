package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"goStore/middlewares"
)

// Handler is a function that can process a request
type Handler func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func NotFoundHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Body:       "Endpoint not supported: " + req.Path,
	}, nil
}

// API defines the exposed endpoints and their handlers
type API struct {
	routes map[string]Handler
}

func (api *API) AddHandler(route string, handler Handler) {
	if api.routes == nil {
		api.routes = map[string]Handler{}
	}
	api.routes[route] = handler
}

func (api *API) GetHandler(path string) Handler {
	if api.routes == nil {
		api.routes = map[string]Handler{}
	}

	/*
		TODO Well, this is obviously naive and terrible. But I wanted to get here on my own before involving frameworks.
			Now that I'm here and it's obvious that I need an actual multiplexor, it's time to go and get gin into the mix.
	*/
	for route, handler := range api.routes {
		if strings.HasPrefix(path, route) {
			return handler
		}
	}

	return NotFoundHandler
}

func (api *API) GetRoutes() map[string]Handler {
	if api.routes == nil {
		api.routes = map[string]Handler{}
	}
	return api.routes
}

/*
TODO: https://gitlab.com/ro-tex/gostore/issues/2
  This should be the main handler function that takes care of auth, executing middlewares, etc.
  The actual handling of requests should be done in separate, per-endpoint handlers.
  Those other handlers should return to this one, so it can execute response middlewares and also check for errors,
  so it can run the error middlewares. Middlewares rule everything! :D
*/
// MasterHandler takes care of all global, route-agnostic tasks, such as running middlewares, system checks,
// initial logging, etc. It then selects the right handler for the request's path (or returns an error).
// With that, its work is done and it delegates to the respective handler by returning it.
//
// MasterHandler does NOT care about authentication/authorisation, as those are route-dependent.
func (api API) MasterHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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
	res, err = api.GetHandler(req.Path)(req)

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
