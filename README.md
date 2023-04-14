# cloud-army-admission-controller

![](/img/2023-04-13_19-04.png) 
This is a [Kubernetes admission controller] to be used as a mutating admission webhook to add a container-init with a custom binary to extract secrets from GCP Secret Manager and to push this secret to the container entrypoint. This solution can be used to compliance with the CIS Kubernetes Benchmark v1.5.1 with the control id: 5.4.1 (no-secrets-as-env-vars). It has been developed as a simple Go web service.

## Installation

### Requirements
* Docker
* kubectl
* cert-manager

## Usage
### Deploy Admission Webhook
To configure the cluster to use the admission webhook and to deploy said webhook, simply run:
```

⚙️  Applying cluster config...
kubectl apply -f manifests/cluster-config/
namespace/apps created
issuer/admission-issuer created
certificate/admission-tls-secret created
issuer/admission-issuer created
mutatingwebhookconfiguration/carmy-kubernetes-webhook created

```
### _🚨 IMPORTANT NOTE: cert-manager is necesary to create Self-Sign certificates🚨_


```
🚀 Deploying carmor-kubernetes-webhook...
kubectl apply -f manifests/webhook/
deployment.apps/carmor-kubernetes-webhook created
service/carmor-kubernetes-webhook created
```

Then, make sure the admission webhook pod is running (in the `mutator` namespace):
```
❯ kubectl get pods
NAME                                        READY   STATUS    RESTARTS   AGE
carmor-kubernetes-webhook-77444566b7-wzwmx   1/1     Running   0          2m21s
```

### Deploying pods
Deploy a test pod that gets secret from secret manager and and print in the console:
```

🚀 Building and Deploying a test pod...
kubectl apply -f manifests/pods/pod-example.yaml
pod/envserver created

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: envserver
spec:
  replicas: 1
  selector:
    matchLabels:
      app: envserver
  template:
    metadata:
      labels:
        app: envserver
    spec:
      containers:
      - name: envserver
        image: xxxxxxxxxx
        imagePullPolicy: Always
        command: ["entrypoint.sh"]

```

### _🚨 NOTE: For test, you should create a docker image with a simple entrypoint that use printenv & sleep with time in seconds and the secret-manager-file🚨_

About the secret-manager-file this is his estructure:

```json
{
    "secrets":[
        {
            "env":"",
            "name":"projects/86303628392/secrets/camunda-compiled-dev/versions/latest"
        }
    ],
    "config":
        {
            "convert_to_uppercase_var_names": true
        }
}

```
### K8S references:

- https://github.com/GoogleCloudPlatform/berglas/tree/main/examples/kubernetes

- https://www.sobyte.net/post/2023-01/cert-manager-admission-webhook/

- https://cert-manager.io/docs/troubleshooting/webhook/

- https://cert-manager.io/docs/installation/helm/

- https://cloud.google.com/anthos-config-management/docs/how-to/using-cis-k8s-benchmark
