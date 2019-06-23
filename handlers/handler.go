package handlers

import (
  "encoding/json"
  "fmt"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "strings"

  "goStore/lib"
  "goStore/middlewares"
)

/*
TODO:
  This should be the main handler function that takes care of auth, executing middlewares, etc.
  The actual handling of requests should be done in separate, per-endpoint handlers.
  Those other handlers should return to this one, so it can execute response middlewares and also check for errors,
  so it can run the error middlewares. Middlewares rule eveything! :D
*/
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

  // Execute request middlewares
  for _, m := range middlewares.Request.GetAll() {
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

      var it lib.Item
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
            Body:       "Failed to create Item. Error: " + string(err.Error()),
          }, err
        }
      }
      fmt.Println("> All good!")
      return events.APIGatewayProxyResponse{
        StatusCode: 201,
        Body:       "Successfully created Item.",
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
      fmt.Println("Trying ot get Item: " + name)
      it, err := getItem(name)
      if err != nil {
        fmt.Println("> Error getting Item: " + err.Error())
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
          Body:       "Couldn't marshal the Item! Error: " + string(err.Error()),
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

func getItem(name string) (lib.Item, error) {
  ddb, err := lib.GetDynamoClient()
  if err != nil {
    fmt.Println("> Error while getting ddb client: " + err.Error())
    return lib.Item{}, err
  }

  res, err := ddb.GetItem(&dynamodb.GetItemInput{
    TableName: aws.String(lib.TableName),
    Key: map[string]*dynamodb.AttributeValue{
      "Name": {
        S: aws.String(name),
      },
    },
  })
  if err != nil {
    fmt.Println("> Error while reading Item: " + err.Error())
    return lib.Item{}, err
  }

  var it lib.Item
  err = dynamodbattribute.UnmarshalMap(res.Item, &it)
  return it, err
}

func createItem(it lib.Item) error {
  ddb, err := lib.GetDynamoClient()
  if err != nil {
    return err
  }

  ItemFields, err := dynamodbattribute.MarshalMap(it)
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
    Item:      ItemFields,
    TableName: aws.String(lib.TableName),
  }

  _, err = ddb.PutItem(input)
  return err
}
