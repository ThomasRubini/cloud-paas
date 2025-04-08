# How to deploy:

## Setup k3s cluster
**To deploy this software, you need a k3s cluster with modified configuration !** 
If you do not have a cluster, you can use [k3d](https://k3d.io) to create one, by running this command in this directory: `k3d cluster create --config ./k3d_config.yaml`

If you already have a cluster, you need to add the following to its configuration:
```
mirrors:
    paas-docker-registry.paas.svc.cluster.local:30005:
    endpoint:
        - "http://localhost:30005"
```
This allows the k3s cluster to access the images we build and push to our local registry.

## Deploy the software in k3s

### Deploying from release
- Copy `values.example.yaml` to `values.prod.yaml` (for example)
- Edit `values.prod.yaml` to match your environment (if you are deploying locally, do not change the server name, keep localhost !)
- Run `helm dependency update` to load helm dependencies that we use
- Run `helm upgrade -i --create-namespace paas ./ -n paas -f values.prod.yaml`

### Deploying source code
If you want to deploy from the source code instead of using the release image, you need to do the following just before installing the `helm upgrade` command:
- `docker build . -t ghcr.io/thomasrubini/cloud-paas:latest`
- `k3d image import ghcr.io/thomasrubini/cloud-paas:latest`
(Adapt the tag if you changed it in your values)

**If you had already deployed the software**, you need to delete the previous deployment before importing the image:
- `kubectl delete deploy -n paas paas-deployment`

## Use in the CLI
At the root of this repo, modify `paas_cli_config.yml` (there is a template named `paas_cli_config.example.yml`) to set `backend_url: http://localhost`. (or the URL you used to deploy the software). This will make the CLI use your fresh new deployment

