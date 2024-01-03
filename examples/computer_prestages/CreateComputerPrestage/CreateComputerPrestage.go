package main

import (
	"encoding/json"
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

	// Define the payload for creating a new computer prestage
	var prestage jamfpro.ResourceComputerPrestage
	payload := `{
		"mandatory": false,
		"mdmRemovable": true,
		"defaultPrestage": false,
		"keepExistingSiteMembership": true,
		"keepExistingLocationInformation": true,
		"requireAuthentication": true,
		"preventActivationLock": true,
		"enableDeviceBasedActivationLock": false,
		"skipSetupItems": {
			"Biometric": true
		},
		"locationInformation": {
			"username": "name",
			"realname": "realName",
			"phone": "123-456-7890",
			"email": "test@jamf.com",
			"room": "room",
			"position": "postion",
			"departmentId": "-1",
			"buildingId": "-1",
			"id": "-1",
			"versionLock": 1
		},
		"purchasingInformation": {
			"leased": true,
			"purchased": true,
			"id": "-1",
			"appleCareId": "abcd",
			"poNumber": "53-1",
			"vendor": "Example Vendor",
			"purchasePrice": "$500",
			"lifeExpectancy": 5,
			"purchasingAccount": "admin",
			"purchasingContact": "true",
			"leaseDate": "2019-01-01",
			"poDate": "2019-01-01",
			"warrantyDate": "2019-01-01",
			"versionLock": 1
		},
		"autoAdvanceSetup": true,
		"installProfilesDuringSetup": true,
		"accountSettings": {
			"payloadConfigured": true,
			"localAdminAccountEnabled": false,
			"hiddenAdminAccount": false,
			"localUserManaged": false,
			"userAccountType": "STANDARD",
			"versionLock": 0,
			"prefillPrimaryAccountInfoFeatureEnabled": false,
			"prefillType": "CUSTOM",
			"preventPrefillInfoFromModification": false,
			"id": "1",
			"adminUsername": "admin",
			"adminPassword": "password",
			"prefillAccountFullName": "TestUser FullName",
			"prefillAccountUserName": "UserName"
		},
		"displayName": "Example Mobile Prestage Name",
		"supportPhoneNumber": "5555555555",
		"supportEmailAddress": "example@example.com",
		"department": "",
		"enrollmentSiteId": "-1",
		"authenticationPrompt": "LDAP authentication prompt",
		"deviceEnrollmentProgramInstanceId": "1",
		"anchorCertificates": [],
		"enrollmentCustomizationId": "0",
		"language": "en",
		"region": "US",
		"customPackageDistributionPointId": "-1",
    "enableRecoveryLock": true,
    "recoveryLockPasswordType": "MANUAL",
    "rotateRecoveryLockPassword": true,
		"prestageInstalledProfileIds": [],
		"recoveryLockPassword": "password123",
		"customPackageIds": []
	}
	`

	// Unmarshal the JSON payload into the prestage struct
	err = json.Unmarshal([]byte(payload), &prestage)
	if err != nil {
		log.Fatalf("Error unmarshaling prestage payload: %v", err)
	}

	// Call the CreateComputerPrestage function
	createdPrestage, err := client.CreateComputerPrestage(&prestage)
	if err != nil {
		log.Fatalf("Error creating computer prestage: %v", err)
	}

	// Print the response
	fmt.Printf("Created Computer Prestage: %+v\n", createdPrestage)
}
