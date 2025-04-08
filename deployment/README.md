# How to deploy:

## Setup k3s cluster
**To deploy this software, you need a k3s cluster with modified configuration !** 
If you do not have a cluster, you can use [k3d](https://k3d.io) to create one with this command `k3d cluster create --config ./config.yaml`

If you already have a cluster, you need to add the following to its configuration:
```
mirrors:
    paas-docker-registry.paas.svc.cluster.local:30005:
    endpoint:
        - "http://localhost:30005"
```
This allows the k3s cluster to access the images we build and push to our local registry.

# Deploy the software in k3s
- Copy `values.yaml.example` to `values.prod.yaml` (for example)
- Edit `values.prod.yaml` to match your environment
- Run `kubectl create namespace paas`
- Run `helm install paas ./ -n paas -f values.prod.yaml`