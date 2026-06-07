package canvas

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateOptionsValidation(t *testing.T) {
	tests := []struct {
		name    string
		opts    CreateOptions
		wantErr bool
	}{
		{
			name:    "empty options",
			opts:    CreateOptions{},
			wantErr: true,
		},
		{
			name: "missing product ID",
			opts: CreateOptions{
				Name: "Test Canvas",
				Kind: KindOpportunity,
			},
			wantErr: true,
		},
		{
			name: "missing name",
			opts: CreateOptions{
				ProductID: "PROD",
				Kind:      KindOpportunity,
			},
			wantErr: true,
		},
		{
			name: "missing kind",
			opts: CreateOptions{
				ProductID: "PROD",
				Name:      "Test Canvas",
			},
			wantErr: true,
		},
		{
			name: "valid options",
			opts: CreateOptions{
				ProductID: "PROD",
				Name:      "Test Canvas",
				Kind:      KindOpportunity,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateOptions(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCreateOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func validateCreateOptions(opts CreateOptions) error {
	if opts.ProductID == "" {
		return errProductIDRequired
	}
	if opts.Name == "" {
		return errNameRequired
	}
	if opts.Kind == "" {
		return errKindRequired
	}
	return nil
}

var (
	errProductIDRequired = errorString("product ID is required")
	errNameRequired      = errorString("canvas name is required")
	errKindRequired      = errorString("canvas kind is required")
)

type errorString string

func (e errorString) Error() string { return string(e) }

func TestBlocksForKind(t *testing.T) {
	tests := []struct {
		kind          Kind
		expectedLen   int
		expectedFirst string
	}{
		{KindOpportunity, 10, "Users & Customers"},
		{KindLeanUX, 8, "Business Problem"},
		{KindBMC, 9, "Key Partners"},
		{"Unknown", 0, ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.kind), func(t *testing.T) {
			blocks := BlocksForKind(tt.kind)
			if len(blocks) != tt.expectedLen {
				t.Errorf("BlocksForKind(%s) returned %d blocks, want %d", tt.kind, len(blocks), tt.expectedLen)
			}
			if tt.expectedLen > 0 && blocks[0] != tt.expectedFirst {
				t.Errorf("BlocksForKind(%s) first block = %s, want %s", tt.kind, blocks[0], tt.expectedFirst)
			}
		})
	}
}

func TestLoadBlocksFromFile(t *testing.T) {
	// Create temp file with valid JSON
	tmpDir := t.TempDir()
	validPath := filepath.Join(tmpDir, "valid.json")
	validContent := `{
		"Users & Customers": "<p>Test users</p>",
		"Problems": "<p>Test problems</p>"
	}`
	if err := os.WriteFile(validPath, []byte(validContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create temp file with invalid JSON
	invalidPath := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(invalidPath, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantLen int
		wantErr bool
	}{
		{
			name:    "valid JSON file",
			path:    validPath,
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "invalid JSON file",
			path:    invalidPath,
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "non-existent file",
			path:    filepath.Join(tmpDir, "missing.json"),
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks, err := LoadBlocksFromFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBlocksFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(blocks) != tt.wantLen {
				t.Errorf("LoadBlocksFromFile() returned %d blocks, want %d", len(blocks), tt.wantLen)
			}
		})
	}
}

func TestOpportunityBlocksComplete(t *testing.T) {
	expected := []string{
		"Users & Customers",
		"Problems",
		"Solution Ideas",
		"Solutions Today",
		"User Value",
		"Adoption Strategy",
		"User Metrics",
		"Business Problem",
		"Business Metrics",
		"Budget",
	}

	if len(OpportunityBlocks) != len(expected) {
		t.Errorf("OpportunityBlocks has %d items, want %d", len(OpportunityBlocks), len(expected))
	}

	for i, block := range expected {
		if OpportunityBlocks[i] != block {
			t.Errorf("OpportunityBlocks[%d] = %s, want %s", i, OpportunityBlocks[i], block)
		}
	}
}

func TestLeanUXBlocksComplete(t *testing.T) {
	expected := []string{
		"Business Problem",
		"Business Outcomes",
		"Users",
		"Benefits",
		"Solutions",
		"Hypotheses",
		"Riskiest Assumption",
		"Smallest Experiment",
	}

	if len(LeanUXBlocks) != len(expected) {
		t.Errorf("LeanUXBlocks has %d items, want %d", len(LeanUXBlocks), len(expected))
	}

	for i, block := range expected {
		if LeanUXBlocks[i] != block {
			t.Errorf("LeanUXBlocks[%d] = %s, want %s", i, LeanUXBlocks[i], block)
		}
	}
}

func TestBMCBlocksComplete(t *testing.T) {
	expected := []string{
		"Key Partners",
		"Key Activities",
		"Key Resources",
		"Value Propositions",
		"Customer Relationships",
		"Channels",
		"Customer Segments",
		"Cost Structure",
		"Revenue Streams",
	}

	if len(BMCBlocks) != len(expected) {
		t.Errorf("BMCBlocks has %d items, want %d", len(BMCBlocks), len(expected))
	}

	for i, block := range expected {
		if BMCBlocks[i] != block {
			t.Errorf("BMCBlocks[%d] = %s, want %s", i, BMCBlocks[i], block)
		}
	}
}
