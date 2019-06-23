package lib

import (
  "fmt"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

const TableName = "innoIvo"

type Item struct {
  Name  string `json:"Name"`
  Value string `json:"Value"`
}

var Middlewares = struct {
  Request  []func(req *events.APIGatewayProxyRequest)
  Response []func(res *events.APIGatewayProxyResponse)
  Error    []func()
}{}

// GetDynamoClient returns what you'd expect
func GetDynamoClient() (*dynamodb.DynamoDB, error) {
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
