package dns

import (
	"encoding/json"
	"fmt"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// loadConfig is a small helper function that decodes JSON configuration into
// the typed config struct.
func loadConfig(cfgJSON *extapi.JSON) (ProviderConfig, error) {
	cfg := ProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	return cfg, nil
}
