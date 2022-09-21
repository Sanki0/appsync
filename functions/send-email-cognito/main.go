package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Request RequestType `json:"request"`
}

type RequestType struct {
	CodeParameter  string             `json:"codeParameter"`
	LinkParameter  string             `json:"linkParameter"`
	UserAttributes UserAttributesType `json:"userAttributes"`
}
type UserAttributesType struct {
	CognitoUser       string `json:"cognito:user_status"`
	Nombre            string `json:"custom:nombre"`
	Email             string `json:"email"`
	Sub               string `json:"sub"`
	UsernameParameter string `json:"usernameParameter"`
}

type UserAttributes struct {
	Email string `json:"email"`
}

func handler(ctx context.Context, event Event) string {
	fmt.Print(event.Request.UserAttributes.Email)
	return "printing"
}

func main() {
	lambda.Start(handler)
}
