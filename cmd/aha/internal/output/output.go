package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Format represents the output format type.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// ParseFormat parses a format string into a Format type.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "table", "":
		return FormatTable, nil
	case "json":
		return FormatJSON, nil
	case "yaml", "yml":
		return FormatYAML, nil
	default:
		return "", fmt.Errorf("unknown format: %s (valid: table, json, yaml)", s)
	}
}

// Printer handles output formatting.
type Printer struct {
	format Format
	writer io.Writer
}

// NewPrinter creates a new Printer with the specified format.
func NewPrinter(format Format) *Printer {
	return &Printer{
		format: format,
		writer: os.Stdout,
	}
}

// Print outputs data in the configured format.
func (p *Printer) Print(data any) error {
	switch p.format {
	case FormatJSON:
		return p.printJSON(data)
	case FormatYAML:
		return p.printYAML(data)
	default:
		return fmt.Errorf("use PrintTable for table format")
	}
}

func (p *Printer) printJSON(data any) error {
	enc := json.NewEncoder(p.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func (p *Printer) printYAML(data any) error {
	enc := yaml.NewEncoder(p.writer)
	enc.SetIndent(2)
	return enc.Encode(data)
}

// Table helps build tabular output.
type Table struct {
	writer  *tabwriter.Writer
	headers []string
}

// NewTable creates a new table printer.
func NewTable(headers ...string) *Table {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	t := &Table{
		writer:  w,
		headers: headers,
	}
	if len(headers) > 0 {
		t.printRow(headers...)
	}
	return t
}

// AddRow adds a row to the table.
func (t *Table) AddRow(values ...string) {
	t.printRow(values...)
}

func (t *Table) printRow(values ...string) {
	_, _ = fmt.Fprintln(t.writer, strings.Join(values, "\t"))
}

// Flush writes the table to output.
func (t *Table) Flush() error {
	return t.writer.Flush()
}

// IsStructured returns true if the format is JSON or YAML.
func (f Format) IsStructured() bool {
	return f == FormatJSON || f == FormatYAML
}
