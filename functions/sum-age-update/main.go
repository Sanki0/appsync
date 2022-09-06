package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoString struct {
	S string `json:"S"`
}
type DynamoInt struct {
	N string `json:"N"`
}

type StudentInput struct {
	Name DynamoString `json:"name"`
	Dni  DynamoString `json:"dni"`
	Age  DynamoInt    `json:"age"`
}

type Student struct {
	Id   string `json:"id"`
	Sort string `json:"sort"`
	Name string `json:"name"`
	Dni  string `json:"dni"`
	Age  int    `json:"age"`
	// Courses []Course `json:"courses"`
}

func sumAge(age int) int {

	ageDiff := age + 10
	return ageDiff

}

func handler(ctx context.Context, event StudentInput) (string, error) {
	TABLE_NAME := os.Getenv("DB")

	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return "", err
	}

	svc := dynamodb.New(sess)
	intAge, _ := strconv.Atoi(event.Age.N)

	student := Student{
		Id:   "STUDENT",
		Sort: event.Dni.S,
		Name: event.Name.S,
		Dni:  event.Dni.S,
		Age:  sumAge(intAge),
	}

	item, err := dynamodbattribute.MarshalMap(student)
	if err != nil {
		fmt.Println("error on marshal")

		return "Error on marshal", err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(TABLE_NAME),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("error on putitem")
		return "error on putitem", err
	}

	str := student.Name
	return str, nil
}

func main() {
	lambda.Start(handler)
}
