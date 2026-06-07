// Package prism provides integration between aha-go and the PRISM ecosystem.
//
// This package implements the export.CanvasProvider interface from prism-roadmap,
// allowing canvases to be exported to Aha.io Strategic Models.
//
// # Usage
//
//	import (
//	    "github.com/grokify/aha-go"
//	    "github.com/grokify/aha-go/prism"
//	    "github.com/grokify/prism-roadmap/canvas/export"
//	)
//
//	// Create Aha client
//	client, _ := aha.NewClient()
//
//	// Create provider
//	provider := prism.NewProvider(client, "PROD")
//
//	// Register with export registry
//	registry := export.NewRegistry()
//	registry.Register(provider)
//
//	// Export a canvas
//	externalID, _ := registry.Export(ctx, "aha", myCanvas)
package prism
