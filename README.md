# cert-manager Webhook for Bunny DNS

cert-manager Webhook for Bunny DNS is a ACME [webhook](https://cert-manager.io/docs/configuration/acme/dns01/webhook/) for [cert-manager](https://cert-manager.io/) allowing users to use [Bunny DNS](https://bunny.net/dns/) for DNS01 challenge.

## Getting started

### Prerequisites

- A [Bunny API Key](https://docs.bunny.net/reference/bunnynet-api-overview)
- A valid domain configured on [Bunny DNS](https://bunny.net/dns/)
- A Kubernetes cluster (v1.29+ recommended)
- [Helm 3](https://helm.sh/) [installed](https://helm.sh/docs/intro/install/) on your computer
- cert-manager [deployed](https://cert-manager.io/docs/installation/) on the cluster

### Installing

> Attention: starting from `0.1.0` the chart's name is now named `cert-manager-webhook-bunny`. 

- Add arbreagile's helm chart repository:

```bash
helm repo add arbreagile https://helm.arbreagile.eu/
helm repo update
```

- Install the chart

```bash
helm install cert-manager-webhook-bunny arbreagile/cert-manager-webhook-bunny
```

The Bunny Webhook is now installed! :tada:

> Refer to the chart's [documentation](https://github.com/arbreagile/helm-charts/blob/master/charts/cert-manager-webhook-bunny/README.md) for more configuration options.

> Alternatively, you may use the provided bundle for a basic install in the cert-manager namespace:
> `kubectl apply -f https://raw.githubusercontent.com/arbreagile/cert-manager-webhook-arbreagile/main/deploy/bundle.yaml`

### How to use it

**Note**: It uses the [cert-manager webhook system](https://cert-manager.io/docs/configuration/acme/dns01/webhook/). Everything after the issuer is configured is just cert-manager. You can find out more in [their documentation](https://cert-manager.io/docs/usage/).

Now that the webhook is installed, here is how to use it.
Let's say you need a certificate for `example.com` (should be registered in Bunny DNS).

First step is to create a secret containing the Bunny Access and Secret keys. Create the `bunny-secret.yaml` file with the following content:
(Only needed if you don't have default credentials as seen above).
```yaml
apiVersion: v1
data:
  api-key: <BUNNY-API-KEY>
kind: Secret
metadata:
  name: bunny-secret
type: Opaque
```

And run:
```bash
kubectl create -f bunny-secret.yaml
```

Next step is to create a cert-manager `Issuer`. Create a `issuer.yaml` file with the following content:
```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: my-bunny-issuer
spec:
  acme:
    email: my-user@example.com
    # this is the acme staging URL
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    # for production use this URL instead
    # server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: my-bunny-private-key-secret
    solvers:
    - dns01:
        webhook:
          groupName: acme.arbreagile.eu
          solverName: bunny
          config:
            # Only needed if you don't have default credentials as seen above.
            apiKeySecretRef:
              key: BUNNY-API-KEY
              name: bunny-secret
```

And run:
```bash
kubectl create -f issuer.yaml
```

Finally, you can now create the `Certificate` object for `example.com`. Create a `certificate.yaml` file with the following content:
```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com
spec:
  dnsNames:
  - example.com
  issuerRef:
    name: my-bunny-issuer
  secretName: example-com-tls
```

And run:
```bash
kubectl create -f certificate.yaml
```

After some seconds, you should see the certificate as ready:
```bash
$ kubectl get certificate example-com
NAME          READY   SECRET            AGE
example-com   True    example-com-tls   1m12s
```

Your certificate is now available in the `example-com-tls` secret!

## Integration testing

Before running the test, you need:
- A valid domain on Bunny DNS (here `example.com`)
- The variable `BUNNY-API-KEY` valid and in the environment

In order to run the integration tests, run:
```bash
TEST_ZONE_NAME=example.com make test
```
