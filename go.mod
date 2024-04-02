module github.com/scaleway/cert-manager-webhook-scaleway

go 1.15

require (
	github.com/cert-manager/cert-manager v1.8.1
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.7.0.20201117145121-3abc1efd92f7
	k8s.io/api v0.23.4
	k8s.io/apiextensions-apiserver v0.23.4
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
)
