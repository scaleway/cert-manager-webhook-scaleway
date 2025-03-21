package dns

import (
	v1 "k8s.io/api/core/v1"
)

// ProviderConfig represents the config used for Bunny DNS
type ProviderConfig struct {
	ApiKey *v1.SecretKeySelector `json:"apiKeySecretRef,omitempty"`
}
