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
	Email            string `json:"email"`
	Password         string `json:"password"`
	Name             string `json:"name"`
	Case             int    `json:"case"`
	ConfirmationCode string `json:"confirmation"`
	Username         string `json:"username"`
	NewPassword      string `json:"newpassword"`
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
		appClientId:   "1q1ima4dt8821ehja023i7uljh",
		userPoolId:    "us-east-1_FxTKZgHWQ",
	}
	fmt.Printf("Email :%s Password: %s \n", event.Email, event.Password)
	fmt.Println("cliente: ", client)

	switch event.Case {
	case 0: // SignUp
		client.SignUp(event.Email, event.Password, event.Name)
	case 1: // AdminCreateUser
		client.AdminCreateUser(event.Email, event.Name)
	case 2:
		client.ConfirmSignUp(event.Email, event.Username, event.ConfirmationCode)
	case 3:
		client.ResendConfirmationCode(event.Email, event.Username)
	case 4:
		client.SignIn(event.Email, event.Password)
	case 5:
		client.ChangePasswordUser(event.Email, event.Password, event.NewPassword)
	case 6:
		client.ForgotPassword(event.Email, event.Password, event.Username)
	case 7:
		client.ConfirmForgotPassword(event.Email, event.Password, event.Username, event.ConfirmationCode)
	}

	fmt.Print(client)
	return "", nil
}

func main() {
	d := deps{}
	lambda.Start(d.handler)
}

func (ctx *awsCognitoClient) SignUp(email string, password string, name string) (string, error) {

	user := &cognito.SignUpInput{
		ClientId: aws.String(ctx.appClientId),
		Username: aws.String(email),
		Password: aws.String(password),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(name),
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

func (ctx *awsCognitoClient) ConfirmSignUp(email string, username string, confirmationCode string) (string, error) {

	user := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String("1q1ima4dt8821ehja023i7uljh"),
		ConfirmationCode: aws.String(confirmationCode),
		Username:         aws.String(username),
	}

	result, err := ctx.cognitoClient.ConfirmSignUp(user)
	if err != nil {
		fmt.Println("Error :", err)
		return "", err
	}
	return result.String(), nil
}

func (ctx *awsCognitoClient) ResendConfirmationCode(email string, username string) (string, error) {

	user := &cognito.ResendConfirmationCodeInput{
		ClientId: aws.String("1q1ima4dt8821ehja023i7uljh"),
		Username: aws.String(username),
	}
	fmt.Println("USER: aaaa ", user.Username)

	result, err := ctx.cognitoClient.ResendConfirmationCode(user)
	if err != nil {
		fmt.Println("Error :", err)
		return "", err
	}
	return result.String(), nil
}

func (ctx *awsCognitoClient) SignIn(email string, password string) (string, error) {
	initiateAuthInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: aws.StringMap(map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		}),
		ClientId: aws.String(ctx.appClientId),
	}

	result, err := ctx.cognitoClient.InitiateAuth(initiateAuthInput)

	if err != nil {
		fmt.Println("Error  : InitiateAuth", err)
		return "", err
	}

	return result.String(), nil
}

func (ctx *awsCognitoClient) ChangePasswordUser(email string, password string, newpassword string) (string, error) {
	fmt.Println("Error 1")
	initiateAuthInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: aws.StringMap(map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		}),
		ClientId: aws.String(ctx.appClientId),
	}

	fmt.Println("Error 2")

	result, err := ctx.cognitoClient.InitiateAuth(initiateAuthInput)

	if err != nil {
		fmt.Println("Error  : InitiateAuth", err)
		return "", err
	}
	fmt.Println(result)
	fmt.Println("Error 2.1")

	accessToken := result.AuthenticationResult.AccessToken

	fmt.Println("Token expires in ", result.AuthenticationResult.ExpiresIn)
	fmt.Println("AccessToken: ", accessToken)
	fmt.Println("Error 3")

	changePasswordInput := &cognito.ChangePasswordInput{
		AccessToken:      aws.String(*result.AuthenticationResult.AccessToken),
		PreviousPassword: aws.String(password),
		ProposedPassword: aws.String(newpassword),
	}

	fmt.Println("Error 4")

	result2, err2 := ctx.cognitoClient.ChangePassword(changePasswordInput)

	if err2 != nil {
		fmt.Println("Error  : ChangePassword", err2)
		return "", err2
	}

	fmt.Println("Error 5")

	return result2.String(), nil
}

func (ctx *awsCognitoClient) ForgotPassword(email string, password string, username string) (string, error) {

	forgotPasswordInput := &cognito.ForgotPasswordInput{
		ClientId: aws.String("1q1ima4dt8821ehja023i7uljh"),
		Username: aws.String(username),
	}

	result2, err2 := ctx.cognitoClient.ForgotPassword(forgotPasswordInput)

	println(result2.CodeDeliveryDetails.DeliveryMedium)

	if err2 != nil {
		fmt.Println("Error  : ChangePassword", err2)
		return "", err2
	}

	return result2.String(), nil
}

func (ctx *awsCognitoClient) ConfirmForgotPassword(email string, newPassword string, username string, confirmationCode string) (string, error) {

	confirmForgotPasswordInput := &cognito.ConfirmForgotPasswordInput{
		ClientId:         aws.String("1q1ima4dt8821ehja023i7uljh"),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(confirmationCode),
		Password:         aws.String(newPassword),
	}

	result2, err2 := ctx.cognitoClient.ConfirmForgotPassword(confirmForgotPasswordInput)

	if err2 != nil {
		fmt.Println("Error  : ChangePassword", err2)
		return "", err2
	}

	return result2.String(), nil
}
