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
		InstanceName:       authConfig.InstanceName,
		OverrideBaseDomain: authConfig.OverrideBaseDomain,
		LogLevel:           logLevel,
		Logger:             logger,
		ClientID:           authConfig.ClientID,
		ClientSecret:       authConfig.ClientSecret,
	}

	// Create a new jamfpro client instance
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	// Example cache settings to update
	newSettings := &jamfpro.ResourceCacheSettings{
		CacheType:                  "ehcache",
		TimeToLiveSeconds:          180,
		TimeToIdleSeconds:          180,
		DirectoryTimeToLiveSeconds: 180,
		EhcacheMaxBytesLocalHeap:   "1GB",
		CacheUniqueId:              "24864549-94ea-4cc1-bb80-d7fb392c6556",
		Elasticache:                false,
		MemcachedEndpoints: []jamfpro.CacheSettingsSubsetMemcachedEndpoint{
			{
				HostName: "localhost",
				Port:     11211,
				Enabled:  true,
			},
		},
	}

	updatedSettings, err := client.UpdateCacheSettings(newSettings)
	if err != nil {
		log.Fatalf("Error updating cache settings: %s", err)
	}

	fmt.Printf("Updated Cache Settings: %+v\n", updatedSettings)
}
