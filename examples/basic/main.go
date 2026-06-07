// Package main demonstrates basic usage of the aha-go client.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/grokify/aha-go"
)

func main() {
	// Create client from environment variables
	// Set AHA_SUBDOMAIN and AHA_API_KEY before running
	client, err := aha.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Connected to Aha: %s\n", client.BaseURL())

	ctx := context.Background()

	// Example: Get a feature
	if len(os.Args) > 1 {
		featureID := os.Args[1]
		fmt.Printf("\nFetching feature: %s\n", featureID)

		// TODO: Uncomment once wrapper methods are implemented
		// feature, err := client.GetFeature(ctx, featureID)
		// if err != nil {
		//     if aha.IsNotFound(err) {
		//         fmt.Println("Feature not found")
		//         return
		//     }
		//     log.Fatalf("Failed to get feature: %v", err)
		// }
		// fmt.Printf("Feature: %s - %s\n", feature.ReferenceNum, feature.Name)
		_ = ctx
		fmt.Println("Feature retrieval not yet implemented")
	}

	fmt.Println("\nClient initialized successfully!")
}
