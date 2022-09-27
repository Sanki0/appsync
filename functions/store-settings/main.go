package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Option struct {
	Title  string `json:"title"`
	Url    string `json:"url"`
	Icon   string `json:"icon"`
	Active bool   `json:"active"`
}

type UserObject struct {
	Id                  string   `json:"id"`
	Sort                string   `json:"sort"`
	Name                string   `json:"name"`
	DocType             string   `json:"docType"`
	Dni                 string   `json:"dni"`
	Gender              string   `json:"gender"`
	BirthDate           string   `json:"birthDate"`
	CountryOfBirth      string   `json:"countryOfBirth"`
	PersonalEmail       string   `json:"personalEmail"`
	MaritalStatus       string   `json:"maritalStatus"`
	PersonalPhone       string   `json:"personalPhone"`
	CountryOfResidence  string   `json:"countryOfResidence"`
	ResidenceDepartment string   `json:"residenceDepartment"`
	Address             string   `json:"address"`
	Area                string   `json:"area"`
	SubArea             string   `json:"subArea"`
	WorkerType          string   `json:"workerType"`
	Email               string   `json:"email"`
	CreationDate        string   `json:"creationDate"`
	EntryDate           string   `json:"entryDate"`
	Phone               string   `json:"phone"`
	Apps                []Option `json:"apps"`
	Menu                []Option `json:"menu"`
	Processes           []Option `json:"processes"`
	UserType            string   `json:"userType"`
	UserStatus          string   `json:"userStatus"`
	Role                string   `json:"role"`
	OfficeRole          string   `json:"officeRole"`
	Days                int      `json:"days"`
	HomeOffice          int      `json:"homeOffice"`
	Boss                string   `json:"boss,omitempty"`
	BossName            string   `json:"bossName,omitempty"`
	User                string   `json:"user"`
	Backup              string   `json:"backup"`
	BackupName          string   `json:"backupName"`
}

type Event struct {
	Settings UserObject `json:"settings"`
}

func handler(ctx context.Context, event Event) (string, error) {
	TABLE_NAME := os.Getenv("TableName")

	timeNow := time.Now()
	now := strings.Split(timeNow.Format(time.RFC3339), "Z")[0] + "Z"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		return "", err
	}

	svc := dynamodb.New(sess)
	in := UserObject{
		Id:                  event.Settings.Email,
		Sort:                "SETTINGS",
		Name:                event.Settings.Name,
		DocType:             event.Settings.DocType,
		Dni:                 event.Settings.Dni,
		Gender:              event.Settings.Gender,
		BirthDate:           event.Settings.BirthDate,
		CountryOfBirth:      event.Settings.CountryOfBirth,
		PersonalEmail:       event.Settings.PersonalEmail,
		MaritalStatus:       event.Settings.MaritalStatus,
		PersonalPhone:       event.Settings.PersonalPhone,
		CountryOfResidence:  event.Settings.CountryOfResidence,
		ResidenceDepartment: event.Settings.ResidenceDepartment,
		Address:             event.Settings.Address,
		Area:                event.Settings.Area,
		SubArea:             event.Settings.SubArea,
		WorkerType:          event.Settings.WorkerType,
		Email:               event.Settings.Email,
		CreationDate:        now,
		EntryDate:           event.Settings.EntryDate,
		Phone:               event.Settings.Phone,
		Apps:                event.Settings.Apps,
		Menu:                event.Settings.Menu,
		Processes:           event.Settings.Processes,
		UserType:            event.Settings.UserType,
		UserStatus:          "UNCONFIRMED",
		Role:                event.Settings.Role,
		OfficeRole:          event.Settings.OfficeRole,
		Days:                event.Settings.Days,
		HomeOffice:          event.Settings.HomeOffice,
		Boss:                event.Settings.Boss,
		BossName:            event.Settings.BossName,
		User:                event.Settings.Email,
		Backup:              "",
		BackupName:          "",
	}

	item, err := MarshalMap(in)
	if err != nil {
		return "", err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(TABLE_NAME),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return "", err
	}

	return "Success", nil
}

func main() {
	lambda.Start(handler)
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
