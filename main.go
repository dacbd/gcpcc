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
	instances, err := service.Instances.List(project, "all").Do()
	fmt.Println(len(instances.Items))
}
