package lib

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Item struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

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
