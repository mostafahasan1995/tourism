package main

import (

	"larsa-tourism-microservices/pkg/api"
	"log"


	_ "github.com/joho/godotenv/autoload"
)

func main() {

	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

}
