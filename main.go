package main

import (
	"os"

	"github.com/jetstack/cert-manager/pkg/acme/webhook/cmd"

	"github.com/scaleway/cert-manager-webhook-scaleway/pkg/dns"
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
