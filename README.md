# cloud-army-kubernetes-webhook (mutator)
WIPPPPPPPPPP
This is a [Kubernetes mutator webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/). It is meant to be used as a mutating admission webhook to add container-init and custom binary to extract and add secret environment to the entrypoint. It has been developed as a simple Go web service without using any framework or boilerplate such as kubebuilder.


## Installation

### Requirements
* Docker
* kubectl

## Usage
### Deploy Admission Webhook
To configure the cluster to use the admission webhook and to deploy said webhook, simply run:
```

⚙️  Applying cluster config...
kubectl apply -f manifests/cluster-config/
namespace/apps created
x
x

🚀 Deploying carmor-kubernetes-webhook...
kubectl apply -f manifests/webhook/
deployment.apps/carmor-kubernetes-webhook created
service/carmor-kubernetes-webhook created
```

Then, make sure the admission webhook pod is running (in the `default` namespace):
```
❯ kubectl get pods
NAME                                        READY   STATUS    RESTARTS   AGE
carmor-kubernetes-webhook-77444566b7-wzwmx   1/1     Running   0          2m21s
```

And hit it's health endpoint from your local machine:
```
❯ curl -k https://localhost:8443/health
OK
```

### Deploying pods
Deploy a valid test pod that gets succesfully created:
```

🚀 Deploying test pod...
kubectl apply -f manifests/pods/xxxxx.yaml
pod/lifespan-seven created
```
You should see in the admission webhook logs that the pod got mutated and validated.

```
