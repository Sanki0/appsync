package main

import (
	"context"
	"reflect"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	UserType string `json:"userType"`
	Email    string `json:"email"`
}

func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func handler(ctx context.Context, event Event) (bool, error) {
	userType := event.UserType
	email := event.Email
	domain := strings.Split(email, "@")[1]
	domainsAvailable := [2]string{"securitygrupo.pe", "protectasecurity.pe"}

	switch userType {
	case "Usuario Protecta", "Administrador Protecta":
		if itemExists(domainsAvailable, domain) {
			return true, nil
		} else {
			return false, nil
		}
	case "Usuario Externo":
		return true, nil
	default:
		return false, nil

	}

}

func main() {
	lambda.Start(handler)
}
