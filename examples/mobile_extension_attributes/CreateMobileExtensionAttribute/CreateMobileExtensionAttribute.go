package main

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client" // Import http_client for logging
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {
	// Define the path to the JSON configuration file
	configFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/clientauth.json"

	// Load the client OAuth credentials from the configuration file
	authConfig, err := jamfpro.LoadClientAuthConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load client OAuth configuration: %v", err)
	}

	// Instantiate the default logger and set the desired log level
	logger := http_client.NewDefaultLogger()
	logLevel := http_client.LogLevelDebug // LogLevelNone // LogLevelWarning // LogLevelInfo  // LogLevelDebug

	// Configuration for the jamfpro
	config := jamfpro.Config{
		InstanceName: authConfig.InstanceName,
		LogLevel:     logLevel,
		Logger:       logger,
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
	}

	// Create a new jamfpro client instance
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	// Define a new mobile device extension attribute
	newAttribute := &jamfpro.ResponseMobileExtensionAttributes{
		Name:             "New Attribute",
		Description:      "This is a test attribute",
		DataType:         "String",
		InventoryDisplay: "General",
		InputType: jamfpro.MobileExtensionAttributeInputType{
			Type: "Text Field",
		},
	}

	// Create the new attribute
	createdAttribute, err := client.CreateMobileExtensionAttribute(newAttribute)
	if err != nil {
		log.Fatalf("Error creating mobile device extension attribute: %v", err)
	}

	fmt.Printf("Created Attribute: %+v\n", createdAttribute)
}
