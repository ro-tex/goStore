package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "innoIvo"

type item struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

var middlewares = struct {
	Request  []func(req *events.APIGatewayProxyRequest)
	Response []func(res *events.APIGatewayProxyResponse)
	Error    []func()
}{}

func main() {
	middlewares.Request = append(middlewares.Request, cleanRequest)
	lambda.Start(Handler)
}

func cleanRequest(req *events.APIGatewayProxyRequest) {
	pathLen := len(req.Path)

	// strip trailing '/'
	if req.Path[pathLen-1] == '/' && len(req.Path) > 1 {
		req.Path = req.Path[:pathLen-1]
	}
}

// TODO try https://github.com/appleboy/gin-lambda

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	/* This is some test code for inspecting the request. */
	//jsonReq, err := json.Marshal(req)
	//if err != nil {
	//	return events.APIGatewayProxyResponse{
	//		StatusCode: 555,
	//		Body:       err.Error(),
	//	}, nil
	//}
	//fmt.Println("Request: " + string(jsonReq))
	//return events.APIGatewayProxyResponse{
	//	StatusCode: 222,
	//	Body:       string(jsonReq),
	//}, nil

	/*
		TODO:
			* error class that can output nice JSON errors
			* decent logging
	*/

	// Execute request middlewares
	for _, m := range middlewares.Request {
		m(&req)
	}

	switch req.HTTPMethod {
	case "POST":
		{
			fmt.Println("> Processing a POST")

			if req.Path != "/v0/doc" {
				return events.APIGatewayProxyResponse{
					StatusCode: 405,
					Body:       "Method not allowed (POST) Path: " + req.Path,
				}, nil
			}

			var it item
			err := json.Unmarshal([]byte(req.Body), &it)
			if err != nil {
				fmt.Println("> Error: " + err.Error())

				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       "Bad Request: The request should have 'Name' and 'Value' string fields.",
				}, err
			}

			err = createItem(it)
			if err != nil {
				fmt.Println("> Error: " + err.Error())

				if err.Error() == "Item already exists!" {
					return events.APIGatewayProxyResponse{
						StatusCode: 400,
						Body:       "Item already exists!",
					}, err
				} else {
					return events.APIGatewayProxyResponse{
						StatusCode: 500,
						Body:       "Failed to create item. Error: " + string(err.Error()),
					}, err
				}
			}
			fmt.Println("> All good!")
			return events.APIGatewayProxyResponse{
				StatusCode: 201,
				Body:       "Successfully created item.",
			}, nil
		}
	case "GET":
		{
			fmt.Println("> Processing a GET")
			if !strings.HasPrefix(req.Path, "/v0/doc/") {
				return events.APIGatewayProxyResponse{
					StatusCode: 404,
					Body:       "Path not found: invalid path.",
				}, nil
			}

			name := strings.Split(req.Path, "/")[3] // yeah, yeah, I know...
			fmt.Println("Trying ot get item: " + name)
			it, err := getItem(name)
			if err != nil {
				fmt.Println("> Error getting item: " + err.Error())
				if strings.HasPrefix(err.Error(), "ResourceNotFoundException") { // TODO What's the real string for this?
					return events.APIGatewayProxyResponse{
						StatusCode: 404,
						Body:       "Not found (not in DB)",
					}, nil
				} else {
					fmt.Println("> Error4: " + err.Error())
					return events.APIGatewayProxyResponse{
						StatusCode: 500,
						Body:       "Oops! Error: " + string(err.Error()),
					}, err
				}
			}
			jsonItem, err := json.Marshal(it)
			if err != nil {
				fmt.Println("> Error5: " + err.Error())
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:       "Couldn't marshal the item! Error: " + string(err.Error()),
				}, err
			}
			fmt.Println("> All good on  processing GET!")
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       string(jsonItem),
			}, nil
		}
	}

	fmt.Println("> Method not allowed.")
	return events.APIGatewayProxyResponse{
		StatusCode: 405,
		Body:       "Method not allowed: " + req.HTTPMethod,
	}, nil
}

func getDynamoClient() (*dynamodb.DynamoDB, error) {
	// TODO Cache this
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		fmt.Println("> Err while getting DynamoDB client: " + err.Error())
		return nil, err
	}

	return dynamodb.New(sess), nil
}

func createItem(it item) error {
	ddb, err := getDynamoClient()
	if err != nil {
		return err
	}

	itemFields, err := dynamodbattribute.MarshalMap(it)
	if err != nil {
		return err
	}

	// check for existence
	// TODO this can be better - check the error, it should be a specific one
	e, _ := getItem(it.Name)
	if e.Name != "" {
		return fmt.Errorf("Item already exists!")
	}

	input := &dynamodb.PutItemInput{
		Item:      itemFields,
		TableName: aws.String(tableName),
	}

	_, err = ddb.PutItem(input)
	return err
}

func getItem(name string) (item, error) {
	ddb, err := getDynamoClient()
	if err != nil {
		fmt.Println("> Error while getting ddb client: " + err.Error())
		return item{}, err
	}

	res, err := ddb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(name),
			},
		},
	})
	if err != nil {
		fmt.Println("> Error while reading item: " + err.Error())
		return item{}, err
	}

	var it item
	err = dynamodbattribute.UnmarshalMap(res.Item, &it)
	return it, err
}
