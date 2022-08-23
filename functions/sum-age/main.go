package main

import (
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Age string `json:"age"`
}

func sumAge(event Event) (int, error) {

	intAge, _ := strconv.Atoi(event.Age)

	ageDiff := intAge + 10

	return ageDiff, nil

}

func main() {
	lambda.Start(sumAge)
}
