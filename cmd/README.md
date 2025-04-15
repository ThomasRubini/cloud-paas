# How to use Cloud-PaaS as an enduser

## Prerequisites
- A functional _Cloud-PaaS_ instance (see [Project Deployment](./deployment/README.md))
- A project that has a ``Dockerfile`` with the following configuration:
    - An ``EXPOSE`` directive that indicates the port on which your application listens
    - An ``HEALTHCHECK``directive (for seamless updates)
- The project's git repository must be accessible by the _Cloud-PaaS_ instance (authentication is not yet supported).

## Quick start

We will be using a exemple app (https://github.com/ThomasRubini/cloud-paas-test) for this tutorial.
Clone it: ``git clone https://github.com/ThomasRubini/cloud-paas-test && cd cloud-paas-test``

### Connecting to the _Cloud-PaaS_ instance
Create `paas_cli.config.yml` in the project's repository:
```yaml
backend_url: <Your instance's URL>
```

### Creating an application
Let's create an application:
```bash
pass-cli app create exempleApp --desc 'An example app' --source-url 'https://github.com/ThomasRubini/cloud-paas-test'
```
You can check that the application was correctly created:
```bash
pass-cli app list
pass-cli app info myNewApp
```

### Creating an environement
An environement is an instance of an application. It automatically updates as the underlaying git branch gets updated.

Create an environement for the "dev" branch:
```bash
pass-cli env myNewApp create dev --branch "dev" --domain "dev-myapp.<your instance's domain>"
```

**TIP:** For local testing of _Cloud-PaaS_ plateform, you can use the ``localhost`` top-level domain:
```bash
pass-cli env myNewApp create dev --branch "dev" --domain "dev-myapp.localhost"
```

You can check the environment state at any time:
```bash
pass-cli env myNewApp info dev
```

You can also `edit`, `delete` or `list` your environements:
```bash
pass-cli env myNewApp list
pass-cli env myNewApp edit dev
pass-cli env myNewApp delete dev
```

## Cleanup
When you are done, delete the environement as well as the application;
```bash
pass-cli env myNewApp delete dev
pass-cli app delete myNewApp
```

