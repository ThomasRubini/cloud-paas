How to deploy:

- Copy `values.yaml.example` to `values.prod.yaml` (for example)
- Edit `values.prod.yaml` to match your environment
- Run `kubectl create namespace paas`
- Run `helm install paas ./ -n paas  -f values.yaml -f values.prod.yaml`
