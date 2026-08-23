package graphql

import (
	"context"
	"fmt"
	"strings"

	genql "github.com/Khan/genqlient/graphql"
	"github.com/grokify/aha-go/graphql/generated"
)

// CustomFieldValue is a custom field value on a record, as returned by
// SetCustomFieldValues.
type CustomFieldValue struct {
	ID         string
	Key        string
	Value      any
	HumanValue string
}

// SetCustomFieldValues sets one or more custom field values on a record
// (Feature, Initiative, Release, and other Aha entity types - see
// generated.CustomFieldableTypeEnum for the full list).
//
// values maps a custom field's API key to its new value. Values are
// usually plain scalars (string/number/bool) matching the custom field's
// type in Aha, not JSON objects.
//
// The underlying GraphQL mutation can report field-level validation
// errors (e.g. an unknown key, or a value that doesn't match the custom
// field's type) inside a successful response payload rather than as a
// transport-level error - this function checks for that and returns a
// non-nil error in either case, so callers never see a silent no-op
// success.
func SetCustomFieldValues(
	ctx context.Context,
	client genql.Client,
	recordID string,
	recordType generated.CustomFieldableTypeEnum,
	values map[string]any,
) ([]CustomFieldValue, error) {
	input := make([]generated.CustomFieldValueInput, 0, len(values))
	for key, value := range values {
		input = append(input, generated.CustomFieldValueInput{Key: key, Value: value})
	}

	resp, err := generated.SetCustomFieldValues(ctx, client, recordID, recordType, input)
	if err != nil {
		return nil, fmt.Errorf("setting custom field values: %w", err)
	}

	if attrs := resp.SetCustomFieldValues.Errors.Attributes; len(attrs) > 0 {
		msgs := make([]string, 0, len(attrs))
		for _, a := range attrs {
			msgs = append(msgs, fmt.Sprintf("%s: %s", a.Name, strings.Join(a.FullMessages, "; ")))
		}
		return nil, fmt.Errorf("setting custom field values: %s", strings.Join(msgs, " | "))
	}

	out := make([]CustomFieldValue, 0, len(resp.SetCustomFieldValues.CustomFieldValues))
	for _, v := range resp.SetCustomFieldValues.CustomFieldValues {
		humanValue := ""
		if v.HumanValue != nil {
			humanValue = *v.HumanValue
		}
		out = append(out, CustomFieldValue{
			ID:         v.Id,
			Key:        v.Key,
			Value:      v.Value,
			HumanValue: humanValue,
		})
	}
	return out, nil
}
