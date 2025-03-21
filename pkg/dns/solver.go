package dns

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/arbreagile/cert-manager-webhook-bunny/pkg/dns/internal/ptr"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/nrdcg/bunny-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	providerName = "bunny"
)

// ProviderSolver is the struct implementing the webhook.Solver interface
// for Bunny DNS
type ProviderSolver struct {
	k8sClient kubernetes.Interface
	client    *bunny.Client
	recordID  *int64
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
	ctx := context.Background()
	domainName := ptr.Pointer(strings.TrimSuffix(ch.ResolvedZone, "."))
	recordOptions := &bunny.AddOrUpdateDNSRecordOptions{
		Type:  ptr.Pointer(bunny.DNSRecordTypeTXT),
		Name:  ptr.Pointer(strings.TrimSuffix(ch.ResolvedFQDN, ptr.Deref(domainName))),
		Value: ptr.Pointer(strconv.Quote(ch.Key)),
		TTL:   ptr.Pointer(int32(60)),
	}
	client, err := p.Client(ch)
	if err != nil {
		return fmt.Errorf("failed to get bunny client: %w", err)
	}
	zoneID, err := findZoneID(client, ctx, ptr.Deref(domainName))
	if err != nil {
		return fmt.Errorf("bunny: failed to get zone ID: %w", err)
	}
	record, err := client.DNSZone.AddDNSRecord(ctx, zoneID, recordOptions)
	if err != nil {
		return fmt.Errorf("bunny: failed to add TXT record: fqdn=%s, zoneID=%d: %w", ch.ResolvedFQDN, zoneID, err)
	}
	p.recordID = record.ID
	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (p *ProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	ctx := context.Background()
	domainName := ptr.Pointer(strings.TrimSuffix(ch.ResolvedZone, "."))
	client, err := p.Client(ch)
	if err != nil {
		return fmt.Errorf("failed to get bunny client: %w", err)
	}
	zoneID, err := findZoneID(client, ctx, ptr.Deref(domainName))
	if err != nil {
		return fmt.Errorf("bunny: failed to get zone ID: %w", err)
	}
	err = client.DNSZone.DeleteDNSRecord(ctx, zoneID, *p.recordID)
	if err != nil {
		return fmt.Errorf("bunny: failed to delete TXT record: fqdn=%s, zoneID=%d: %w", ch.ResolvedFQDN, zoneID, err)
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
	p.k8sClient = cl
	return nil
}
