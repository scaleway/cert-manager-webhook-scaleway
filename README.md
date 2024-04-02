# cert-manager Webhook for Scaleway DNS

cert-manager Webhook for Scaleway DNS is a ACME [webhook](https://cert-manager.io/docs/configuration/acme/dns01/webhook/) for [cert-manager](https://cert-manager.io/) allowing users to use [Scaleway DNS](https://www.scaleway.com/en/docs/scaleway-dns/) for DNS01 challenge.

## Getting started

### Prerequisites

- A [Scaleway Access Key and a Scaleway Secret Key](https://www.scaleway.com/en/docs/generate-api-keys/)
- A valid domain configured on [Scaleway DNS](https://www.scaleway.com/en/docs/scaleway-dns/)
- A Kubernetes cluster (v1.22+ recommended)
- [Helm 3](https://helm.sh/) [installed](https://helm.sh/docs/intro/install/) on your computer
- cert-manager [deployed](https://cert-manager.io/docs/installation/) on the cluster

### Installing

Once everything is set up, you can now install the Scaleway Webhook:
- Clone this repository: 
```bash
git clone https://github.com/scaleway/cert-manager-webhook-scaleway.git
```

- Run:
```bash
helm install scaleway-webhook deploy/scaleway-webhook
```
- Alternatively, you can install the webhook with default credentials with: 
```bash
helm install scaleway-webhook deploy/scaleway-webhook --set secret.accessKey=<YOUR-ACCESS-KEY> --set secret.secretKey=<YOUR-SECRET_KEY>
```

The Scaleway Webhook is now installed! :tada:

### How to use it

**Note**: It uses the [cert-manager webhook system](https://cert-manager.io/docs/configuration/acme/dns01/webhook/). Everything after the issuer is configured is just cert-manager. You can find out more in [their documentation](https://cert-manager.io/docs/usage/).

Now that the webhook is installed, here is how to use it.
Let's say you need a certificate for `example.com` (should be registered in Scaleway DNS).

First step is to create a secret containing the Scaleway Access and Secret keys. Create the `scaleway-secret.yaml` file with the following content:
(Only needed if you don't have default credentials as seen above).
```yaml
apiVersion: v1
stringData:
  SCW_ACCESS_KEY: <YOUR-SCALEWAY-ACCESS-KEY>
  SCW_SECRET_KEY: <YOUR-SCALEWAY-SECRET-KEY>
kind: Secret
metadata:
  name: scaleway-secret
type: Opaque
```

And run:
```bash
kubectl create -f scaleway-secret.yaml
```

Next step is to create a cert-manager `Issuer`. Create a `issuer.yaml` file with the following content:
```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: my-scaleway-issuer
spec:
  acme:
    email: my-user@example.com
    # this is the acme staging URL
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    # for production use this URL instead
    # server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: my-scaleway-private-key-secret
    solvers:
    - dns01:
        webhook:
          groupName: acme.scaleway.com
          solverName: scaleway
          config:
            # Only needed if you don't have default credentials as seen above.
            accessKeySecretRef:
              key: SCW_ACCESS_KEY
              name: scaleway-secret
            secretKeySecretRef:
              key: SCW_SECRET_KEY
              name: scaleway-secret
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
    name: my-scaleway-issuer
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
- A valid domain on Scaleway DNS (here `example.com`)
- The variables `SCW_ACCESS_KEY` and `SCW_SECRET_KEY` valid and in the environment

In order to run the integration tests, run:
```bash
TEST_ZONE_NAME=example.com make test
```
