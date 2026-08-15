package reporter

import (
	"encoding/json"
)

// ExportToJSON serializes suite execution summaries to formatted JSON bytes.
func ExportToJSON(result *SuiteResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}
