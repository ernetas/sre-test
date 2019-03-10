package main

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "log"
  "net/http"
  "os"
)

type Item struct {
  User string`json:"user"`
  Bdate string`json:"bdate"`
}

func getAwsRegion() string {
  region := os.Getenv("AWS_REGION")
  if len(region) == 0 {
    // Port 9000 is an easy choice for local development as it wouldn't require root permissions
    return "eu-central-1"
  }
  return region
}

func writeBirthday(username, date string, w http.ResponseWriter) {
  awsRegion := getAwsRegion() // this is probably handled by AWS SDK as well, but I want to pass a sensible default value for local development where I don't have AWS_REGION exported for some reasons beyond this task
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String(awsRegion)},
  )
  // Create DynamoDB client
  svc := dynamodb.New(sess)

  input := &dynamodb.UpdateItemInput{
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
      ":s": {
        S: aws.String(date),
      },
    },
    Key: map[string]*dynamodb.AttributeValue{
      "user": {
        S: aws.String(username),
      },
    },
    ReturnValues: aws.String("UPDATED_NEW"),
    UpdateExpression: aws.String("set bdate = :s"),
    TableName: aws.String("birthday"),
  }

  _, err = svc.UpdateItem(input)
  if err != nil {
    log.Print("Got error calling PutItem: ", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprint(w, "Got error calling PutItem: ", err.Error())
  } else {
    // Give HTTP 204 status code on success
    w.WriteHeader(http.StatusNoContent)
    log.Printf("PUT: %s %s", username, date)
  }
}

func getBirthday(username string, w http.ResponseWriter) {
  awsRegion := getAwsRegion() // this is probably handled by AWS SDK as well, but I want to pass a sensible default value for local development where I don't have AWS_REGION exported for some reasons beyond this task
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String(awsRegion)},
  )
  // Create DynamoDB client
  svc := dynamodb.New(sess)

  var queryInput = &dynamodb.QueryInput{
    Limit:     aws.Int64(1),
    TableName: aws.String("birthday"),
    KeyConditions: map[string]*dynamodb.Condition{
      "user": {
        ComparisonOperator: aws.String("EQ"),
          AttributeValueList: []*dynamodb.AttributeValue{
            {
              S: aws.String(username),
            },
          },
      },
    },
  }
  resp, err := svc.Query(queryInput)
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintln(w, "Failed to retrieve birthday (ERR1)")
  } else {
    ItemObj := []Item{}
    err = dynamodbattribute.UnmarshalListOfMaps(resp.Items, &ItemObj)
    if len(ItemObj) != 1 {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintln(w, "Failed to retrieve birthday (ERR2)")
    } else {
      date := ItemObj[0].Bdate
      log.Printf("GET: %s %s", username, date)

      fmt.Fprintln(w, fmt.Sprintf("{ \"message\": \"Hello, %s! %s\"}", username, templateMessage(date)))
    }
  }
}