<p align="center">
    <img src="/img/logo.png" width="30%" align="center">
</p>

# cloud-army-secret-admission-controller


![](/img/2023-04-13_19-04.png) 
This is a [Kubernetes admission controller] to be used as a mutating admission webhook to add a container-init with a custom binary that extract secrets from GCP Secret Manager and to push this secrets to the container entrypoint sub-process. This solution can be used to compliance with the CIS Kubernetes Benchmark v1.5.1 specially with the control id: 5.4.1 (no-secrets-as-env-vars).

## Installation

### Deploy Admission Webhook
To configure the cluster to use the admission webhook and to deploy said webhook, simply run:
```
git clone https://github.com/cloud-army/secret-mutator-admission-webhook.git

helm install cloud-army-secret-injector secret-mutator

```
### _🚨 IMPORTANT NOTE: Cert-manager controller should be installed in your cluster 🚨_

Then, make sure the admission webhook pod is running (in the `mutator` namespace):
```
❯ k get all -n mutator
NAME                                           READY   STATUS    RESTARTS   AGE
pod/carmy-kubernetes-webhook-86f44f55d-qrshk   1/1     Running   0          34h
pod/test-deployment-5dbbd7b8b8-d2n6l           1/1     Running   0          2d7h

NAME                               TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
service/carmy-kubernetes-webhook   ClusterIP   10.192.48.14   <none>        443/TCP   4d2h

NAME                                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/carmy-kubernetes-webhook   1/1     1            1           2d13h
deployment.apps/test-deployment            1/1     1            1           9d

NAME                                                 DESIRED   CURRENT   READY   AGE
replicaset.apps/carmy-kubernetes-webhook-86f44f55d   1         1         1       2d13h
replicaset.apps/test-deployment-5dbbd7b8b8           1         1         1       9d

```
## Usage
### Deploying pods
Build and Deploy a test pod that gets secrets from GCP Secret Manager and print its in the pod console, remember that: The namespace where running the applications should be labeled with 'admission-webhook: enabled':
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
        command: ["entrypoint.sh"] <<<<< Use entrypoint.sh command as a standard name

```

### _🚨 IMPORTANT NOTE: For test, you should create a docker image with a simple entrypoint that use printenv & sleep with time in seconds, a envsecrets-config.json file, and running the pods using Workload Identity🚨_

About the envsecrets-config.json file, it is the place were declaring the GCP Secrets resources that you need consume, and it's his estructure is:

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
For more information about the envsecrets-config.json file, check this repo https://github.com/cloud-army/envsecrets

### K8S references:

- https://github.com/GoogleCloudPlatform/berglas/tree/main/examples/kubernetes

- https://www.sobyte.net/post/2023-01/cert-manager-admission-webhook/

- https://cert-manager.io/docs/troubleshooting/webhook/

- https://cert-manager.io/docs/installation/helm/

- https://cloud.google.com/anthos-config-management/docs/how-to/using-cis-k8s-benchmark

- https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity?hl=es-419
