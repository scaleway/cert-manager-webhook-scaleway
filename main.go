package main

import (
	"os"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"

	"github.com/arbreagile/cert-manager-webhook-bunny/pkg/dns"
)

// GroupName is the name under which the webhook will be available
var GroupName = os.Getenv("GROUP_NAME")

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	cmd.RunWebhookServer(GroupName,
		&dns.ProviderSolver{},
	)
}
