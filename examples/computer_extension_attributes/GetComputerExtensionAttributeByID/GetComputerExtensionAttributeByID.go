package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {
	// Define the path to the JSON configuration file inside the main function
	configFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/clientauth.json"

	// Load the client OAuth credentials from the configuration file
	authConfig, err := jamfpro.LoadClientAuthConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load client OAuth configuration: %v", err)
	}

	// Configuration for the jamfpro
	config := jamfpro.Config{
		InstanceName: authConfig.InstanceName,
		DebugMode:    true,
		Logger:       jamfpro.NewDefaultLogger(),
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
	}

	// Create a new jamfpro client instance
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	// Provide the ID of the computer extension attribute you want to fetch
	attributeID := 11 // You can change this ID as required

	// Call GetComputerExtensionAttributeByID function
	attribute, err := client.GetComputerExtensionAttributeByID(attributeID)
	if err != nil {
		log.Fatalf("Error fetching Computer Extension Attribute by ID: %v", err)
	}

	// Pretty print the attribute in XML
	attributeXML, err := xml.MarshalIndent(attribute, "", "    ") // Indent with 4 spaces
	if err != nil {
		log.Fatalf("Error marshaling Computer Extension Attribute data: %v", err)
	}
	fmt.Println("Fetched Computer Extension Attribute by ID:\n", string(attributeXML))
}
