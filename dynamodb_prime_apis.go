package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

var (
	client    *dynamodb.Client
	tableName string
	//DYNAMODB_ENDPOINT string
)

func init() {
	//DYNAMODB_ENDPOINT = "http://janus-localstack:4566"
	client, _ = NewDynamoDBClient()
	tableName = "prime-apis"
}

type PrimeAPI struct {
	Id      int    `yaml:"id" json:"id" dynamodbav:"prime_api_id"`
	ApiName string `yaml:"api_name" json:"api_name" dynamodbav:"prime_api_name"`
	Enabled bool   `yaml:"enabled" json:"enabled" dynamodbav:"enabled"`
}

func NewDynamoDBClient() (*dynamodb.Client, error) {
	DynamoDBEndpoint := viper.GetString("DYNAMODB_ENDPOINT")
	//DynamoDBEndpoint := DYNAMODB_ENDPOINT
	//log.Println(DynamoDBEndpoint)
	var client *dynamodb.Client
	/*cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}*/

	//Specifying AWS Region
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		return nil, err
	}
	log.Println("cfg: ", cfg)
	fmt.Println("in LOCAL mode")
	client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolver = dynamodb.EndpointResolverFromURL(DynamoDBEndpoint)
	})

	//Experiment (Error)
	/*
		sess, err := session.NewSession(&aws.Config{
				Region: aws.String("us-west-2")},
		)
		client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.EndpointResolver = dynamodb.EndpointResolverFromURL("https://test.us-west-2.amazonaws.com")
		})*/
	//client = dynamodb.New(sess, &aws.Config{Endpoint: aws.String("https://test.us-west-2.amazonaws.com")})
	return client, nil
}

func AllPrimeAPIs(ctx context.Context) (*[]PrimeAPI, error) {
	out, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName:      &tableName,
		ConsistentRead: aws.Bool(true),
	},
	)

	if err != nil {
		return nil, err
	}

	var primeAPIS []PrimeAPI
	err = av.UnmarshalListOfMaps(out.Items, &primeAPIS)
	if err != nil {
		return nil, err
	}

	return &primeAPIS, nil
}

// SeedPrimeAPIs inserts (upserts, really) these Prime APIs into DynamoDB
// so we can toggle their status as needed
func seedThoseAPIs() error {
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
	table, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String("prime-apis"),
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
	log.Println("Table: ", table)
	log.Println("Error: ", err)

	for _, api := range primeAPIs["prime_apis"] {
		err := insertPrimeApiItem(&api)
		if err != nil {
			log.Printf("Unable to insert prime API configuration: %s", err.Error())
			return err
		}
	}

	return nil
}

func insertPrimeApiItem(primeAPI *PrimeAPI) error {
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
	if err != nil {
		return err
	}

	fmt.Printf("Seeded PrimeAPI %d %s: (enabled %t)", primeAPI.Id, primeAPI.ApiName,
		primeAPI.Enabled)
	return nil
}

func SetEnabledTo(ctx context.Context, primeApiId int, enabled bool) error {

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"prime_api_id": &types.AttributeValueMemberN{Value: strconv.Itoa(primeApiId)},
		},
		UpdateExpression: aws.String("set enabled = :newVal"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newVal": &types.AttributeValueMemberBOOL{
				Value: enabled,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
