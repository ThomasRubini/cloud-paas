# Cloud PaaS CLI Documentation

## App Commands

The `app` command allows you to interact with your applications in the Cloud PaaS platform.

### Usage
```
cli app [command] [arguments]
```

### Commands

#### Create an application
```
cli app create <app_name> --source-url <url> --source-username <username> --source-password <password>
```

Create a new application with the specified name and source repository.

**Arguments:**
- `<app_name>`: Name of the application (required)

**Flags:**
- `--desc`: Description of the application (optional)
- `--source-url`: Source URL of the application (required)
- `--source-username`: Username for repository authentication (optional)
- `--source-password`: Password for repository authentication (optional) - For GitHub, this can be a personal access token

**Example:**
```
cli app create myapp --source-url https://github.com/username/repo --source-username myuser --source-password mytoken
```

#### List applications
```
cli app list
```

Lists all applications associated with your account.

#### Get application information
```
cli app info <app_name>
```

Retrieves detailed information about a specific application.

**Arguments:**
- `<app_name>`: Name of the application to retrieve information about (required)

Output includes Application ID, name, description, and repos source URL

#### Delete an application
```
cli app delete <app_name>
```

Removes an application from your account.

**Arguments:**
- `<app_name>`: Name of the application to delete (required)

## Environment Commands

The `env` command allows you to manage environments for your applications.

### Usage
```
cli env <app_name> <command>
```

### Commands

#### Create an environment
```
cli env <app_name> create <env_name> --branch <branch> --domain <domain>
```

Creates a new environment for the specified application.

**Arguments:**
- `<app_name>`: Name of the application (required)
- `<env_name>`: Name of the environment to create (required)

**Flags:**
- `--branch`, `-b`: Branch to use for the environment (required)
- `--domain`, `-d`: Domain to use for the environment (ie. `main.localhost`) (required)

#### List environments
```
cli env <app_name> list
```

Lists all environments for the specified application.

**Arguments:**
- `<app_name>`: Name of the application (required)

#### Get environment information
```
cli env <app_name> info <env_name>
```

Retrieves detailed information about a specific environment.

**Arguments:**
- `<app_name>`: Name of the application (required)
- `<env_name>`: Name of the environment (required)

Output includes environment name, branch it is based on, domain name and environment variables.

#### Edit environment
```
cli env <app_name> edit <env_name> --branch <branch> --domain <domain>
```

Updates branch or domain for a specific environment.

**Arguments:**
- `<app_name>`: Name of the application (required)
- `<env_name>`: Name of the environment to edit (required)

**Flags:**
- `--branch`, `-b`: Branch name that the environment will be based on. (optional)
- `--domain`, `-d`: Domain for the environment (ie. `main.localhost`) (optional)

#### Edit environment variables
```
cli env <app_name> vars <env_name>
```

Opens an editor to modify environment variables for a specific environment.

**IMPORTANT: Variables need to be written in the YAML format !**

**Arguments:**
- `<app_name>`: Name of the application (required)
- `<env_name>`: Name of the environment (required)

#### Delete an environment
```
cli env <app_name> delete <env_name>
```

Removes an environment from an application.

**Arguments:**
- `<app_name>`: Name of the application (required)
- `<env_name>`: Name of the environment to delete (required)

#### Redeploy an environment
```
cli env <app_name> redeploy <env_name>
```

Triggers a redeployment of a specific environment, In case you wish to restart your application at any time.

**Arguments:**
- `<app_name>`: Name of the application (required)
- `<env_name>`: Name of the environment to redeploy (required)