package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// Define the ID of the advanced computer search
const advancedComputerSearchID = 123 // Replace 123 with the actual ID

func main() {
	// Define the path to the JSON configuration file
	configFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/clientauth.json"

	// Load the client OAuth credentials from the configuration file
	authConfig, err := jamfpro.LoadClientAuthConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load client OAuth configuration: %v", err)
	}

	// Configuration for Jamf Pro
	config := jamfpro.Config{
		InstanceName: authConfig.InstanceName,
		DebugMode:    true,
		Logger:       jamfpro.NewDefaultLogger(),
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
	}

	// Create a new Jamf Pro client instance
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	// Call GetAdvancedComputerSearchByID function using the constant ID
	advancedComputerSearch, err := client.GetAdvancedComputerSearchByID(advancedComputerSearchID)
	if err != nil {
		log.Fatalf("Error fetching advanced computer search by ID: %v", err)
	}

	// Pretty print the advanced computer search in XML
	advancedComputerSearchXML, err := xml.MarshalIndent(advancedComputerSearch, "", "    ") // Indent with 4 spaces
	if err != nil {
		log.Fatalf("Error marshaling advanced computer search data: %v", err)
	}
	fmt.Println("Fetched Advanced Computer Search by ID:\n", string(advancedComputerSearchXML))
}
