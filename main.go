package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

func main() {
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
	credentials, err := google.CredentialsFromJSON(oauth2.NoContext, []byte(credData), compute.ComputeScope)
	if err != nil {
		log.Fatalln(err)
	}
	project := credentials.ProjectID
	service, err := compute.New(oauth2.NewClient(oauth2.NoContext, credentials.TokenSource))
	if err != nil {
		log.Fatalln(err)
	}
	instanceCount := 0

	zones, err := service.Zones.List(project).Do()
	if err != nil {
		log.Fatalln(err)
	}
	for _, zone := range zones.Items {
		instances, err := service.Instances.List(project, zone.Name).Do()
		if err != nil {
			log.Fatalln(err)
		}
		instanceCount += len(instances.Items)
	}

	/*
		regions, err := service.Regions.List(project).Do()
		if err != nil {
			log.Fatalln(err)
		}
		for _, item := range regions.Items {
			for _, zone := range item.Zones {
				instances, err := service.Instances.List(project, zone).Do()
				if err != nil {
					log.Fatalln(err)
				}
				instanceCount += len(instances.Items)
			}
		}
	*/
	fmt.Println(instanceCount)
}
