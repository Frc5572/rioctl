
package googledrive

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed credentials.json
var embeddedCredentials []byte

// Credentials returns the embedded Google OAuth client configuration.
func Credentials() ([]byte, error) {
	var parsed struct {
		Installed json.RawMessage `json:"installed"`
		Web       json.RawMessage `json:"web"`
	}

	if err := json.Unmarshal(embeddedCredentials, &parsed); err != nil {
		return nil, fmt.Errorf("parse embedded OAuth credentials: %w", err)
	}

	if len(parsed.Installed) == 0 && len(parsed.Web) == 0 {
		return nil, fmt.Errorf("embedded credentials must contain an installed or web OAuth client")
	}

	return embeddedCredentials, nil
}
