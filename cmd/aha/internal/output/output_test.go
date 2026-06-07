package output

import (
	"bytes"
	"testing"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
		wantErr  bool
	}{
		{"table", FormatTable, false},
		{"TABLE", FormatTable, false},
		{"", FormatTable, false},
		{"json", FormatJSON, false},
		{"JSON", FormatJSON, false},
		{"yaml", FormatYAML, false},
		{"YAML", FormatYAML, false},
		{"yml", FormatYAML, false},
		{"invalid", "", true},
		{"xml", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseFormat(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatIsStructured(t *testing.T) {
	tests := []struct {
		format   Format
		expected bool
	}{
		{FormatTable, false},
		{FormatJSON, true},
		{FormatYAML, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.IsStructured(); got != tt.expected {
				t.Errorf("Format(%q).IsStructured() = %v, want %v", tt.format, got, tt.expected)
			}
		})
	}
}

func TestPrinterJSON(t *testing.T) {
	data := map[string]string{"key": "value"}
	var buf bytes.Buffer

	printer := &Printer{
		format: FormatJSON,
		writer: &buf,
	}

	if err := printer.Print(data); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	expected := "{\n  \"key\": \"value\"\n}\n"
	if got := buf.String(); got != expected {
		t.Errorf("Print() output = %q, want %q", got, expected)
	}
}

func TestPrinterYAML(t *testing.T) {
	data := map[string]string{"key": "value"}
	var buf bytes.Buffer

	printer := &Printer{
		format: FormatYAML,
		writer: &buf,
	}

	if err := printer.Print(data); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	expected := "key: value\n"
	if got := buf.String(); got != expected {
		t.Errorf("Print() output = %q, want %q", got, expected)
	}
}

func TestTable(t *testing.T) {
	// Just verify it doesn't panic
	table := NewTable()
	table.AddRow("a", "b")
	if err := table.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}

	// Test with headers
	table2 := NewTable("Col1", "Col2")
	table2.AddRow("val1", "val2")
	if err := table2.Flush(); err != nil {
		t.Errorf("Flush() error = %v", err)
	}
}
