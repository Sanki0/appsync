package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type deps struct {
}

type CognitoClient interface {
	ResendAdminCreateUser(email string) (string error)
	GetUser(email string) ([]Response, error)
}

type awsCognitoClient struct {
	cognitoClient *cognito.CognitoIdentityProvider
	appClientId   string
	userPoolId    string
}

type Response struct {
	Username            string    `json:"username"`
	Enabled             bool      `json:"enabled"`
	AccountStatus       string    `json:"accountStatus"`
	Email               string    `json:"email"`
	EmailVerified       string    `json:"emailVerified"`
	PhoneNumberVerified string    `json:"phoneNumberVerified"`
	Updated             time.Time `json:"updated"`
	Created             time.Time `json:"created"`
}

type Event struct {
	Email string `json:"email"`
}

func (d *deps) handler(ctx context.Context, event Event) (string, error) {

	email := event.Email
	domain := strings.Split(email, "@")[1]

	if !(domain == "securitygrupo.pe") || !(domain == "protectasecurity.pe") {
		var result string
		var response []Response
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
			appClientId:   os.Getenv("APP_CLIENT_ID"),
			userPoolId:    os.Getenv("USER_POOL_ID"),
		}

		result, err = client.ResendAdminCreateUser(event.Email)

		if err != nil {
			fmt.Println("Error :", err)
			return "", err
		}
		fmt.Println("CLIENTE :", client)
		fmt.Println("Response :", response)

		fmt.Println("result :", result)

		return "OK", nil
	}

	return "Fail", nil

}

func main() {
	d := deps{}
	lambda.Start(d.handler)
}

func (ctx *awsCognitoClient) ResendAdminCreateUser(email string) (string, error) {

	user := &cognito.AdminCreateUserInput{
		UserPoolId:    aws.String(ctx.userPoolId),
		Username:      aws.String(email),
		MessageAction: aws.String("RESEND"),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}
	fmt.Println("USER: ", user)

	result, err := ctx.cognitoClient.AdminCreateUser(user)
	if err != nil {
		fmt.Println("Error : AdminCreateUser", err)
		return "", err
	}
	return result.String(), nil
}
