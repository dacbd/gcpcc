package main

import (
	"context"
	"fmt"
	"log"
	"os"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	credData := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_DATA")
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(credData) == 0 && len(credPath) == 0 {
		log.Fatalln("no GCP credentials provided")
	}
	if len(credData) == 0 {
		jsondata, err := os.ReadFile(credPath)
		if err != nil {
			log.Fatalln(err)
		}
		credData = string(jsondata)
	}
	client, err := compute.NewClient(ctx, option.WithCredentialsJSON([]byte(credData)))
	if err != nil {
		log.Fatalln(err)
	}
	client.InstanceList()
	fmt.Println("1")
}
