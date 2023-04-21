package main

// Code is adapted from python code found here: https://youtu.be/7m_q1ldzw0U

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type MathEvent struct {
	Base     string `json:"base"`
	Exponent string `json:"exponent"`
}

type Answer struct {
	StatusCode int    `json:"code"`
	Body       string `json:"body"`
}

// Table item for DynamoDB table
type Item struct {
	ID                 string `json:"id"`
	LatestGreetingTime string `json:"latest_greeting_time"`
}

// store the current time in a human readable format in a variable
var now = time.Now().String()

func HandleRequest(ctx context.Context, event MathEvent) (Answer, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	base, err := strconv.Atoi(event.Base)
	if err != nil {
		log.Fatal(err)
	}

	exponent, err := strconv.Atoi(event.Exponent)
	if err != nil {
		log.Fatal(err)
	}

	result := math.Pow(float64(base), float64(exponent))

	item := Item{
		ID:                 fmt.Sprint(result),
		LatestGreetingTime: now,
	}

	attribute, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new math item: %s", err)
	}

	tableName := "MathDatabase"

	input := &dynamodb.PutItemInput{
		Item:      attribute,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	return Answer{StatusCode: 200, Body: fmt.Sprint("Your result is " + fmt.Sprint(result))}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
