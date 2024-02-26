package main

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {
	// Define the path to the JSON configuration file
	configFilePath := "/Users/dafyddwatkins/localtesting/jamfpro/clientconfig.json"

	// Initialize the Jamf Pro client with the HTTP client configuration
	client, err := jamfpro.BuildClientWithConfigFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize Jamf Pro client: %v", err)
	}

	profileName := "Corporate Wireless" // Replace with the actual profile name
	subset := "desired_subset"          // Replace with the desired subset
	profile, err := client.GetMobileDeviceConfigurationProfileByNameWithSubset(profileName, subset)
	if err != nil {
		log.Fatalf("Error fetching mobile device configuration profile by name and subset: %v", err)
	}

	fmt.Printf("Profile: %+v\n", profile)
}
