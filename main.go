package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {

	fmt.Printf("DYNAMODB_ENDPOINT %s", viper.GetString("DYNAMODB_ENDPOINT"))
	fmt.Printf("PRIME_APIS_TABLE_NAME is %s", viper.GetString("PRIME_APIS_TABLE_NAME"))

	err := seedThoseAPIs()
	if err != nil {
		fmt.Println("Unable to load seed data")
	}
}

