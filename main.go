package main

import (
	"context"
	"fmt"
	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/viper"
	"log"
	"time"
)

func main() {

	fmt.Printf("DYNAMODB_ENDPOINT %s", viper.GetString("DYNAMODB_ENDPOINT"))
	fmt.Printf("PRIME_APIS_TABLE_NAME is %s", viper.GetString("PRIME_APIS_TABLE_NAME"))

	//Mini demo

	primeAPI := PrimeAPI{
		Id:      1,
		ApiName: "CoverageEligibilityRequest API",
		Enabled: true,
	}

	marshalledAPI, err := av.MarshalMap(primeAPI)

	var tableName = "prime-apis"
	client, err := NewDynamoDBClient()
	if err != nil {
		log.Printf("Error creating client: %s", err.Error())
	}
	log.Println("Table Name: ", tableName)
	log.Println("Client: ", client)
	_, err = client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		//TableName:   &tableName,
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("prime_api_id"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},

		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("prime_api_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
	})

	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}
	log.Printf("Created table %s", tableName)
	start := time.Now()
	input := &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      marshalledAPI,
	}
	_, err = client.PutItem(context.Background(), input)

	if err != nil {
		log.Println("Failed on put item initially: ", err)
		log.Println("Putting in item again...")
	}

	//Looping through until we can PutItem
	/*for err != nil {
		_, err = client.PutItem(context.Background(), input)
	}*/

	log.Println("Total time for the table: ", time.Since(start).Seconds(), "seconds")

	/*err := seedThoseAPIs()
	if err != nil {
		fmt.Println("Unable to load seed data")
	}*/
}
