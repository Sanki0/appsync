package main

import (
	"context"
	"fmt"
	"os"

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
	Name string `json:"name"`
	Dni  string `json:"dni"`
	Age  int    `json:"age"`
}

type Student struct {
	Id   string `json:"id"`
	Sort string `json:"sort"`
	Name string `json:"name"`
	Dni  string `json:"dni"`
	Age  int    `json:"age"`
	// Courses []Course `json:"courses"`
}

func handler(ctx context.Context, event StudentInput) (string, error) {
	TABLE_NAME := os.Getenv("DB")

	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return "", err
	}

	svc := dynamodb.New(sess)
	// intAge, _ := strconv.Atoi(event.Age)

	student := &Student{
		Id:   "STUDENT",
		Sort: event.Dni,
		Name: event.Name,
		Dni:  event.Dni,
		Age:  event.Age,
	}

	item, err := MarshalMap(student)
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

func MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av, err := getEncoder().Encode(in)
	if err != nil || av == nil || av.M == nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	return av.M, nil
}

func getEncoder() *dynamodbattribute.Encoder {
	encoder := dynamodbattribute.NewEncoder()
	encoder.NullEmptyString = false
	return encoder
}

func main() {
	lambda.Start(handler)
}
