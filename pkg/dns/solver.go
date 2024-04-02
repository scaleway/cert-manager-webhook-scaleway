package dns

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	providerName = "scaleway"
)

// ProviderSolver is the struct implementing the webhook.Solver interface
// for Scaleway DNS
type ProviderSolver struct {
	client kubernetes.Interface
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource
func (p *ProviderSolver) Name() string {
	return providerName
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (p *ProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	domainAPI, err := p.getDomainAPI(ch)
	if err != nil {
		return err
	}

	request := &domain.UpdateDNSZoneRecordsRequest{
		DNSZone: strings.TrimSuffix(ch.ResolvedZone, "."),
		Changes: []*domain.RecordChange{
			{
				Set: &domain.RecordChangeSet{
					IDFields: &domain.RecordIdentifier{
						Name: strings.TrimSuffix(strings.TrimSuffix(ch.ResolvedFQDN, ch.ResolvedZone), "."),
						Type: domain.RecordTypeTXT,
						Data: scw.StringPtr(strconv.Quote(ch.Key)),
					},
					Records: []*domain.Record{
						{
							Name: strings.TrimSuffix(strings.TrimSuffix(ch.ResolvedFQDN, ch.ResolvedZone), "."),
							Data: strconv.Quote(ch.Key),
							Type: domain.RecordTypeTXT,
							TTL:  60,
						},
					},
				},
			},
		},
	}

	_, err = domainAPI.UpdateDNSZoneRecords(request)
	if err != nil {
		return fmt.Errorf("failed to update DNS zone recrds: %w", err)
	}

	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (p *ProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	domainAPI, err := p.getDomainAPI(ch)
	if err != nil {
		return err
	}

	request := &domain.UpdateDNSZoneRecordsRequest{
		DNSZone: strings.TrimSuffix(ch.ResolvedZone, "."),
		Changes: []*domain.RecordChange{
			{
				Delete: &domain.RecordChangeDelete{
					IDFields: &domain.RecordIdentifier{
						Name: strings.TrimSuffix(strings.TrimSuffix(ch.ResolvedFQDN, ch.ResolvedZone), "."),
						Data: scw.StringPtr(strconv.Quote(ch.Key)),
						Type: domain.RecordTypeTXT,
					},
				},
			},
		},
	}

	_, err = domainAPI.UpdateDNSZoneRecords(request)
	if err != nil {
		return fmt.Errorf("failed to update DNS zone recrds: %w", err)
	}

	return nil
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (p *ProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {

	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return fmt.Errorf("failed to get kubernetes client: %w", err)
	}

	p.client = cl

	return nil
}
