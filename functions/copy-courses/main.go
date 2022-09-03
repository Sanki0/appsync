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

type Course struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type StudentInput struct {
	Dni1 string `json:"dni1"`
	Dni2 string `json:"dni2"`
}

type Student struct {
	Id      string   `json:"id"`
	Sort    string   `json:"sort"`
	Name    string   `json:"name"`
	Dni     string   `json:"dni"`
	Courses []Course `json:"courses"`
}

func handler(ctx context.Context, event StudentInput) (string, error) {
	TABLE_NAME := os.Getenv("DB")

	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return "", err
	}

	svc := dynamodb.New(sess)

	query, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("STUDENT"),
			},
			"sort": {
				S: aws.String(event.Dni1),
			},
		},
	})

	student1 := Student{}

	err1 := dynamodbattribute.UnmarshalMap(query.Item, &student1)
	if err1 != nil {
		return "", err
	}

	student2 := Student{
		Id:      "STUDENT",
		Sort:    event.Dni2,
		Name:    student1.Name,
		Dni:     event.Dni2,
		Courses: student1.Courses,
	}

	item, err := dynamodbattribute.MarshalMap(student2)
	if err != nil {
		fmt.Println("error on marshal")

		return "Error on marshal", err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(TABLE_NAME),
	}

	svc2 := dynamodb.New(sess)

	_, err = svc2.PutItem(input)
	if err != nil {
		fmt.Println("error on putitem")
		return "error on putitem", err
	}

	str := student1.Name
	return str, nil
}

func main() {
	lambda.Start(handler)
}
