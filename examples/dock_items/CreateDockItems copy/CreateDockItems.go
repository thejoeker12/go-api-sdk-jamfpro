package main

import (
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {

	// Initialize the Jamf Pro client with the HTTP client configuration
	client, err := jamfpro.BuildClientWithEnv()
	if err != nil {
		log.Fatalf("Failed to initialize Jamf Pro client: %v", err)
	}

	// Define the details of the dock item to create
	newDockItem := &jamfpro.ResourceDockItem{
		Name:     "Safari6",
		Type:     "App",
		Path:     "file://localhost/Applications/Safari.app/",
		Contents: "string",
	}

	// Call the CreateDockItem function
	result, err := client.CreateDockItem(newDockItem)
	if err != nil {
		log.Fatalf("Error creating dock item: %v", err)
	}

	// Output the result
	fmt.Printf("Created Dock Item: %+v\n", result)
}