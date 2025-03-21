package dns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arbreagile/cert-manager-webhook-bunny/pkg/dns/internal/ptr"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/nrdcg/bunny-go"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (p *ProviderSolver) Client(ch *v1alpha1.ChallengeRequest) (*bunny.Client, error) {
	if p.client != nil {
		return p.client, nil
	}

	config, err := loadConfig(ch.Config)
	if err != nil {
		return nil, err
	}

	if config.ApiKey != nil {
		apiKeySecret, err := p.k8sClient.CoreV1().Secrets(ch.ResourceNamespace).Get(context.Background(), config.ApiKey.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not get secret %s: %w", config.ApiKey.Name, err)
		}
		apiKeyData, ok := apiKeySecret.Data[config.ApiKey.Key]
		if !ok {
			return nil, fmt.Errorf("could not get key %s in secret %s", config.ApiKey.Key, config.ApiKey.Name)
		}
		apiKey := string(apiKeyData)

		p.client = bunny.NewClient(apiKey)
		return p.client, nil
	} else {
		return nil, fmt.Errorf("no api key provided in secrets")
	}
}

func findZoneID(client *bunny.Client, ctx context.Context, domainName string) (int64, error) {
	if domainName == "" {
		return 0, fmt.Errorf("empty domain name")
	}

	zones, err := client.DNSZone.List(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to list DNS zones: %w", err)
	}

	zoneMap := make(map[string]int64)
	for _, zone := range zones.Items {
		zoneMap[ptr.Deref(zone.Domain)] = ptr.Deref(zone.ID)
	}

	zoneID, found := zoneMap[domainName]
	if !found {
		return 0, fmt.Errorf("DNS zone not found for domain: %s", domainName)
	}

	return zoneID, nil
}
