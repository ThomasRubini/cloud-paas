# Build and deploy the paas-backend image from source code to the k3d cluster
# Must be run from this directory
set -xe
docker build . -t paas-backend:latest
kubectl delete deploy -n paas paas-deployment || true
k3d image import paas-backend:latest

cd deployment
helm upgrade -i --create-namespace paas ./ -n paas -f values.prod.yaml
