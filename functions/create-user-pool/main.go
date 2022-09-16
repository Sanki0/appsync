package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type deps struct {
}

type CognitoClient interface {
	SignUp(email string, password string) (string, error)
	AdminCreateUser(email string) (string, error)
}

type awsCognitoClient struct {
	cognitoClient *cognito.CognitoIdentityProvider
	appClientId   string
	userPoolId    string
}

type Event struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Case     int    `json:"case"`
}

func (d *deps) handler(ctx context.Context, event Event) (string, error) {
	// CONECTAR SESSION CON AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(*aws.String("us-east-1"))},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect session, %v", err))
	}
	// INICIAR SESSION EN COGNITO
	svc := cognito.New(sess)

	client := awsCognitoClient{
		cognitoClient: svc,
		appClientId:   "7hhotj7f76mlpi9vbpn9u5hic",
		userPoolId:    "us-east-1_6ExOuLULG",
	}
	fmt.Printf("Email :%s Password: %s \n", event.Email, event.Password)
	fmt.Println("cliente: ", client)

	switch event.Case {
	case 0: // SignUp
		client.SignUp(event.Email, event.Password)
	case 1: // AdminCreateUser
		client.AdminCreateUser(event.Email, event.Name)
	}

	fmt.Print(client)
	return "", nil
}

func main() {
	d := deps{}
	lambda.Start(d.handler)
}

func (ctx *awsCognitoClient) SignUp(email string, password string) (string, error) {

	user := &cognito.SignUpInput{
		ClientId: aws.String(ctx.appClientId),
		Username: aws.String(email),
		Password: aws.String(password),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}
	fmt.Println("USER: ", user)

	result, err := ctx.cognitoClient.SignUp(user)
	if err != nil {
		fmt.Println("Error :", err)
		return "", err
	}
	return result.String(), nil
}

func (ctx *awsCognitoClient) AdminCreateUser(email string, name string) (string, error) {

	user := &cognito.AdminCreateUserInput{
		UserPoolId: aws.String(ctx.userPoolId),
		Username:   aws.String(email),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}
	fmt.Println("USER: aaaa ", user)

	result, err := ctx.cognitoClient.AdminCreateUser(user)
	if err != nil {
		fmt.Println("Error :", err)
		return "", err
	}
	return result.String(), nil
}
