package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (

	// HTMLBody ...
	HTMLBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	//The email body for recipients with non-HTML email clients.

	// The character encoding for the email.
	charSet = "UTF-8"
)

type deps struct {
}

type Event struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	HTMLBody string `json:"template"`
	Stage    string
}

func (d *deps) handler(ctx context.Context, event Event) (string, error) {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	body := `¡Hola %s!
	¡Bienvenido/a a la Oficina Virtual de Protecta Security!
	Es necesario que inicies sesión. Para ello debes realizar lo siguiente:
		 1. Ingresa aquí {{url}}.
		 2. Clic en el botón "Iniciar sesión con Google".
		 3. En la nueva ventana abierta, clic en su cuenta corporativa de Protecta.
	¡Listo! Para volver a iniciar sesión deberás seguir los pasos anteriores.`

	//  // The HTML body for the email.
	//  testBody :=  "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
	//  "<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
	//  "<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	// //The email body for recipients with non-HTML email clients.
	// testText := "This email was sent with Amazon SES using the AWS SDK for Go."

	event.Stage = os.Getenv("Stage")

	message := fmt.Sprintf(body, event.Name)
	t, err := template.New("mailhtml").Parse(event.HTMLBody)
	if err != nil {
		fmt.Println(err.Error())
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, event); err != nil {
		fmt.Println(err.Error())
	}

	resultT := tpl.String()

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(event.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(resultT),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(message),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String("Invitación a Oficina Virtual: Cuenta Nueva"),
			},
		},
		Source: aws.String(os.Getenv("SENDER")),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "error", err

	}

	fmt.Println("Email Sent to address: " + event.Email)

	fmt.Println(result)
	return "success", nil
}

func main() {
	d := deps{}
	lambda.Start(d.handler)
}
