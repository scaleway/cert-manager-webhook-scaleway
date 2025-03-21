//go:build integration
// +build integration

package tests

import (
	"os"
	"testing"

	bunnyDns "github.com/arbreagile/cert-manager-webhook-bunny/pkg/dns"
	dns "github.com/cert-manager/cert-manager/test/acme"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting current working dir: %s", err.Error())
	}

	fixture := dns.NewFixture(&bunnyDns.ProviderSolver{},
		dns.SetResolvedZone(zone),
		dns.SetAllowAmbientCredentials(true),
		dns.SetBinariesPath(currentDir+"/kubebuilder/bin"),
		dns.SetManifestPath(currentDir+"/testdata"),
		dns.SetStrict(true),
	)

	fixture.RunConformance(t)
}
