package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

var (
	tableName = "prime-apis"
)

type PrimeAPI struct {
	Id      int    `yaml:"id" json:"id" dynamodbav:"prime_api_id"`
	ApiName string `yaml:"api_name" json:"api_name" dynamodbav:"prime_api_name"`
	Enabled bool   `yaml:"enabled" json:"enabled" dynamodbav:"enabled"`
}

func NewDynamoDBClient() (*dynamodb.Client, error) {
	DynamoDBEndpoint := "http://localstack:4566"
	log.Printf("Setting dynamo endpoint to  %s", DynamoDBEndpoint)
	var client *dynamodb.Client

	//Specifying AWS Region
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		return nil, err
	}
	client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolver = dynamodb.EndpointResolverFromURL(DynamoDBEndpoint)
	})
	return client, nil
}

// SeedPrimeAPIs inserts (upserts, really) these Prime APIs into DynamoDB
// so we can toggle their status as needed
func seedThoseAPIs() error {
	client, err := NewDynamoDBClient()
	if err != nil {
		log.Printf("Error creating client: %s", err.Error())
		return err
	}
	seedsYml := "prime_api_seeds.yml"

	var primeAPIs = make(map[string][]PrimeAPI)
	yf, err := os.ReadFile(seedsYml)

	if err != nil {
		log.Printf("Unable to read prime API configuration: %s", err.Error())
		return err
	}

	err = yaml.Unmarshal(yf, &primeAPIs)
	if err != nil {
		log.Printf("Unable to unmarshal prime API configuration: %s", err.Error())
		return err
	}
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
		return err
	}
	log.Printf("Created table %s", tableName)

	start := time.Now()
	for _, api := range primeAPIs["prime_apis"] {
		err := insertPrimeApiItem(&api, client)

		if err != nil {
			log.Printf("Unable to insert prime API configuration: %s", err.Error())
			return err
		}
	}
	log.Println("The table is ready in: ", time.Since(start).Seconds(), "seconds")

	return nil
}

func insertPrimeApiItem(primeAPI *PrimeAPI, client *dynamodb.Client) error {
	marhsalledAPI, err := av.MarshalMap(primeAPI)
	if err != nil {
		return err
	}

	input := dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      marhsalledAPI,
	}

	_, err = client.PutItem(context.Background(), &input)
	enabled, exists := input.Item["enabled"]
	log.Printf("Enabled Exists: %t, value: %v\n", exists, enabled)
	api_id, id_exists := input.Item["prime_api_id"]
	log.Printf("ID Exists: %t, value: %v\n", id_exists, api_id)
	api_name, name_exists := input.Item["prime_api_name"]
	log.Printf("Name Exists: %t, value: %v\n", name_exists, api_name)
	for err != nil {
		_, err = client.PutItem(context.Background(), &input)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Seeded PrimeAPI %d %s: (enabled %t)", primeAPI.Id, primeAPI.ApiName,
		primeAPI.Enabled)
	return nil
}
