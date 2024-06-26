package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {
	// Define the path to the JSON configuration file
	configFilePath := "/Users/joseph/github/terraform-sandbox/clientconfig.json"

	// Initialize the Jamf Pro client with the HTTP client configuration
	client, err := jamfpro.BuildClientWithConfigFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize Jamf Pro client: %v", err)
	}

	// Define a new policy with all required fields
	newPolicy := &jamfpro.ResourcePolicy{
		// General
		General: jamfpro.PolicySubsetGeneral{
			Name:    fmt.Sprintf("Policy Test-%v", rand.Intn(10000)),
			Enabled: false,
		},
		Scope: &jamfpro.PolicySubsetScope{
			AllComputers: true,
			Computers: &[]jamfpro.PolicyDataSubsetComputer{
				{
					ID: 16,
				},
				{
					ID: 15,
				},
			},
		},
		SelfService: &jamfpro.PolicySubsetSelfService{
			SelfServiceCategories: &[]jamfpro.PolicySubsetSelfServiceCategory{
				{
					Name:      "Applications",
					ID:        5,
					DisplayIn: true,
					FeatureIn: false,
				},
			},
		},
	}

	policyXML, err := xml.MarshalIndent(newPolicy, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling policy data: %v", err)
	}
	fmt.Println("Policy Details to be Sent:\n", string(policyXML))

	// Call CreatePolicy function
	createdPolicy, err := client.CreatePolicy(newPolicy)
	if err != nil {
		log.Fatalf("Error creating policy: %v", err)
	}

	// Pretty print the created policy details in XML
	policyXML, err = xml.MarshalIndent(createdPolicy, "", "    ") // Indent with 4 spaces and use '='
	if err != nil {
		log.Fatalf("Error marshaling policy details data: %v", err)
	}
	fmt.Println("Created Policy Details:\n", string(policyXML))
}
