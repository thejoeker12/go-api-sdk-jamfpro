package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

const (
	maxConcurrentRequestsAllowed = 5 // Maximum allowed concurrent requests.
	defaultTokenLifespan         = 30 * time.Minute
	defaultBufferPeriod          = 5 * time.Minute
	scriptFilePath               = "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/examples/support_files/scriptfile.sh" // Replace with your script file path
)

func main() {
	// Define the path to the JSON configuration file inside the main function
	configFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/clientauth.json"

	// Load the client OAuth credentials from the configuration file
	authConfig, err := http_client.LoadClientAuthConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load client OAuth configuration: %v", err)
	}

	// Configuration for the jamfpro
	config := jamfpro.Config{
		InstanceName:          authConfig.InstanceName,
		DebugMode:             true,
		Logger:                jamfpro.NewDefaultLogger(),
		MaxConcurrentRequests: maxConcurrentRequestsAllowed,
		TokenLifespan:         defaultTokenLifespan,
		BufferPeriod:          defaultBufferPeriod,
		ClientID:              authConfig.ClientID,
		ClientSecret:          authConfig.ClientSecret,
	}

	// Create a new jamfpro client instanceclient,
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	// Read script contents from a file
	scriptFile, err := os.ReadFile(scriptFilePath)
	if err != nil {
		log.Fatalf("Error reading script file: %v", err)
	}
	scriptContents := string(scriptFile)

	// Define a sample script for testing
	sampleScript := &jamfpro.ResponseScript{
		Name:     "Sample Script",
		Category: "None",
		Filename: "string",
		Info:     "Script information",
		Notes:    "Sample Script",
		Priority: "Before",
		Parameters: jamfpro.Parameters{
			Parameter4:  "string",
			Parameter5:  "string",
			Parameter6:  "string",
			Parameter7:  "string",
			Parameter8:  "string",
			Parameter9:  "string",
			Parameter10: "string",
			Parameter11: "string",
		},
		OSRequirements: "string",
		ScriptContents: scriptContents,
	}

	// Call CreateScriptByID function
	createdScript, err := client.CreateScriptByID(sampleScript)
	if err != nil {
		log.Fatalf("Error creating script: %v", err)
	}

	// Pretty print the created script details in XML
	createdScriptXML, err := xml.MarshalIndent(createdScript, "", "    ") // Indent with 4 spaces
	if err != nil {
		log.Fatalf("Error marshaling created script data: %v", err)
	}
	fmt.Println("Created Script Details:\n", string(createdScriptXML))
}