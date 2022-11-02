package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

func main() {
	githubOutputFilename := os.Getenv("GITHUB_OUTPUT")
	credData := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_DATA")
	var credentials *google.Credentials
	var err error
	// get credentials
	if len(credData) != 0 {
		credentials, err = google.CredentialsFromJSON(context.Background(), []byte(credData), compute.ComputeReadonlyScope)
	} else {
		credentials, err = google.FindDefaultCredentials(context.Background(), compute.ComputeReadonlyScope)
	}
	if err != nil {
		log.Fatalln(err)
	}

	// support OIDC
	if credentials.ProjectID == "" {
		// 	Coerce Credentials to handle GCP OIDC auth
		//	Common ProjectID ENVs:
		//		https://github.com/google-github-actions/auth/blob/b05f71482f54380997bcc43a29ef5007de7789b1/src/main.ts#L187-L191
		//		https://github.com/hashicorp/terraform-provider-google/blob/d6734812e2c6a679334dcb46932f4b92729fa98c/google/provider.go#L64-L73
		coercedProjectID := multiEnvLoadFirst([]string{
			"CLOUDSDK_CORE_PROJECT",
			"CLOUDSDK_PROJECT",
			"GCLOUD_PROJECT",
			"GCP_PROJECT",
			"GOOGLE_CLOUD_PROJECT",
			"GOOGLE_PROJECT",
		})
		if coercedProjectID == "" {
			// last effort to load
			fromCredentialsID, err := coerceOIDCCredentials(credentials.JSON)
			if err != nil {
				log.Fatalln(fmt.Errorf("couldn't extract the project identifier from the given credentials!: [%w]", err))
			}
			coercedProjectID = fromCredentialsID
		}
		credentials.ProjectID = coercedProjectID
	}

	// create serivce
	project := credentials.ProjectID
	service, err := compute.New(oauth2.NewClient(context.Background(), credentials.TokenSource))
	if err != nil {
		log.Fatalln(err)
	}
	instanceCount := 0

	// count instances
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
	if githubOutputFilename != "" {
		file, _ := os.OpenFile(githubOutputFilename, os.O_APPEND|os.O_WRONLY, 0644)
		defer file.Close()
		file.WriteString(fmt.Sprintf("total=%v\n", instanceCount))
	} else {
		fmt.Printf("::set-output name=total::%v", instanceCount)
	}
}

// https://github.com/iterative/terraform-provider-iterative/blob/b9cd04a981df2b1426a67c58d506bdf9669eca5e/iterative/gcp/provider.go#L355-L370
func coerceOIDCCredentials(credentialsJSON []byte) (string, error) {
	var credentials map[string]interface{}
	if err := json.Unmarshal(credentialsJSON, &credentials); err != nil {
		return "", err
	}

	if url, ok := credentials["service_account_impersonation_url"].(string); ok {
		re := regexp.MustCompile(`^https://iamcredentials\.googleapis\.com/v1/projects/-/serviceAccounts/.+?@(?P<project>.+)\.iam\.gserviceaccount\.com:generateAccessToken$`)
		if match := re.FindStringSubmatch(url); match != nil {
			return match[1], nil
		}
		return "", errors.New("failed to get project identifier from service_account_impersonation_url")
	}

	return "", errors.New("unable to load service_account_impersonation_url")
}

func multiEnvLoadFirst(envs []string) string {
	for _, val := range envs {
		if env_value := os.Getenv(val); env_value != "" {
			return env_value
		}
	}
	return ""
}
