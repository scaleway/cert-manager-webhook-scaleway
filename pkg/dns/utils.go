package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/scaleway/cert-manager-webhook-scaleway/pkg/util"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
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

func (p *ProviderSolver) getDomainAPI(ch *v1alpha1.ChallengeRequest) (*domain.API, error) {
	config, err := loadConfig(ch.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	var accessKey string
	var secretKey string

	if ch.AllowAmbientCredentials {
		accessKey = os.Getenv(scw.ScwAccessKeyEnv)
		secretKey = os.Getenv(scw.ScwSecretKeyEnv)
	}

	if config.AccessKey != nil && config.SecretKey != nil {
		accessKeySecret, err := p.client.CoreV1().Secrets(ch.ResourceNamespace).Get(context.Background(), config.AccessKey.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not get secret %s: %w", config.AccessKey.Name, err)
		}
		secretKeySecret, err := p.client.CoreV1().Secrets(ch.ResourceNamespace).Get(context.Background(), config.SecretKey.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not get secret %s: %w", config.SecretKey.Name, err)
		}

		accessKeyData, ok := accessKeySecret.Data[config.AccessKey.Key]
		if !ok {
			return nil, fmt.Errorf("could not get key %s in secret %s", config.AccessKey.Key, config.AccessKey.Name)
		}

		secretKeyData, ok := secretKeySecret.Data[config.SecretKey.Key]
		if !ok {
			return nil, fmt.Errorf("could not get key %s in secret %s", config.SecretKey.Key, config.SecretKey.Name)
		}

		accessKey = string(accessKeyData)
		secretKey = string(secretKeyData)
	}

	scwClient, err := scw.NewClient(
		scw.WithEnv(),
		scw.WithAuth(accessKey, secretKey),
		scw.WithUserAgent("cert-manager-webhook-scaleway/"+util.GetVersion().Version),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize scaleway client: %w", err)
	}

	domainAPI := domain.NewAPI(scwClient)

	return domainAPI, nil
}
