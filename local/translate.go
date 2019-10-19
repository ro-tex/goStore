package local

import (
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
)

type GinResponse struct {
	// TODO Maybe headers as well?
	StatusCode int
	Response   gin.H
}

// TransGin2AwsReq converts a Gin Request into an APIGatewayProxyRequest
func TransGin2AwsReq(ctx *gin.Context) (events.APIGatewayProxyRequest, error) {
	r := events.APIGatewayProxyRequest{
		Path:                            ctx.Request.URL.Path,
		HTTPMethod:                      ctx.Request.Method,
		Headers:                         map[string]string{},
		MultiValueHeaders:               ctx.Request.Header,
		QueryStringParameters:           map[string]string{},
		MultiValueQueryStringParameters: ctx.Request.URL.Query(),
		RequestContext:                  events.APIGatewayProxyRequestContext{HTTPMethod: ctx.Request.Method},
		Body:                            "",
		// Resource:                        "",
		// PathParameters:                  map[string]string{},
		// StageVariables:                  map[string]string{},
		// IsBase64Encoded:                 false,
	}

	// get the normal headers
	for k, v := range r.MultiValueHeaders {
		if len(v) == 1 {
			r.Headers[k] = v[0]
			delete(r.MultiValueHeaders, k)
		}
	}

	// get the body
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return r, err
	}
	ctx.Request.Body.Close()
	r.Body = string(body)

	// get the normal query params
	for k, v := range r.MultiValueQueryStringParameters {
		if len(v) == 1 {
			r.QueryStringParameters[k] = v[0]
			delete(r.MultiValueQueryStringParameters, k)
		}
	}

	return r, nil
}

// TransAwsRes2Gin converts an APIGatewayProxyResponse into an object that Gin can return as response
func TransAwsRes2Gin(r *events.APIGatewayProxyResponse) GinResponse {
	res := GinResponse{
		StatusCode: r.StatusCode,
		Response: gin.H{
			"StatusCode": r.StatusCode,
			"Body":       r.Body,
			// Headers           map[string]string   `json:"headers"`
			// MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
			// IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
		},
	}

	return res
}
