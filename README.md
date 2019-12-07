# WIP: Serve Client - Painless Docker Deployments

[![Build Status](https://travis-ci.com/loeffel-io/serve.svg?token=diwUYjrdo8kHiwiMCFuq&branch=master)](https://travis-ci.com/loeffel-io/serve)

## Run

```bash
curl ..
```

## Example

### .serve.yml - Basic Apache httpd service

```yaml
host: 'localhost:8080'
name: 'test-project'

files:
  - index.html
  - images
  - Dockerfile

service:
  build:
    context: "%build_dir%"
    dockerfile: Dockerfile
  ports:
    - 3000:80
```

### .serve.yml - Basic MySQL Service

```yaml
host: 'localhost:8080'
name: 'test-mysql'

service:
  image: 'mysql'
  command: --default-authentication-plugin=mysql_native_password
  restart: always
  environment:
    - MYSQL_ROOT_PASSWORD=example
```
  