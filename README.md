# cloud-army-secret-admission-controller
![Cloud Army Spartan Logo](/img/logo.png) 

![](/img/2023-04-13_19-04.png) 
This is a [Kubernetes admission controller] to be used as a mutating admission webhook to add a container-init with a custom binary that extract secrets from GCP Secret Manager and to push this secrets to the container entrypoint sub-process. This solution can be used to compliance with the CIS Kubernetes Benchmark v1.5.1 specially with the control id: 5.4.1 (no-secrets-as-env-vars).

## Installation

### Requirements
* Docker
* kubectl
* cert-manager
* golang

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
### _🚨 IMPORTANT NOTE: cert-manager controller is necesary to create the Admission Controller Self-Signed certificates, and the namespace where running the applications should be labeled with 'admission-webhook: enabled'🚨_
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
## Usage
### Deploying pods
Deploy a test pod that gets secrets from GCP Secret Manager and print its in the pod console:
```

🚀 Building and Deploying a test pod...
kubectl apply -f manifests/pods-example/pod-example.yaml
pod/envserver created

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: envprinter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: envprinter
  template:
    metadata:
      labels:
        app: envprinter
    spec:
      containers:
      - name: envprinter
        image: xxxxxxxxxx
        imagePullPolicy: Always
        command: ["entrypoint.sh"]

```

### _🚨 IMPORTANT NOTE: For test, you should create a docker image with a simple entrypoint that use printenv & sleep with time in seconds and a envsecrets-config.json🚨_

About the envsecrets-config.json, it is the place were declaring the GCP Secrets resources that you need consume, and it's his estructure is:

```json
{
    "secrets":[
        {
            "env":"",
            "name":"projects/PROJECT_NUMBER/secrets/YOUR_SECRET_NAME/versions/latest"
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
