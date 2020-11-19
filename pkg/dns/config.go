package dns

import (
	"k8s.io/api/core/v1"
)

// ProviderConfig represents the config used for Scaleway DNS
type ProviderConfig struct {
	AccessKey *v1.SecretKeySelector `json:"accessKeySecretRef,omitempty"`
	SecretKey *v1.SecretKeySelector `json:"secretKeySecretRef,omitempty"`
}
