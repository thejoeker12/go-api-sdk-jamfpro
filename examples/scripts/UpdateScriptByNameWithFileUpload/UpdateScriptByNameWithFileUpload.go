package main

import (
	"encoding/xml"
	"fmt"
	"io"
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
	scriptNameToUpdate           = "Embedded Sample Script" // The name of the script to update.
)

func main() {
	configFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/clientauth.json"
	scriptFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/examples/support_files/scriptfile.sh"

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

	file, err := os.Open(scriptFilePath)
	if err != nil {
		log.Fatalf("Error opening script file: %v", err)
	}
	defer file.Close()

	scriptContents, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading script file: %v", err)
	}

	sampleScript := &jamfpro.ResponseScript{
		Name:     scriptNameToUpdate,
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
		ScriptContents: string(scriptContents),
	}

	updatedScript, err := client.UpdateScriptByName(sampleScript)
	if err != nil {
		log.Fatalf("Error updating script by name: %v", err)
	}

	updatedScriptXML, err := xml.MarshalIndent(updatedScript, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling updated script data: %v", err)
	}
	fmt.Println("Updated Script Details:\n", string(updatedScriptXML))
}