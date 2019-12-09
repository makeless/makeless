<img width="180" src="https://raw.githubusercontent.com/loeffel-io/serve-server/master/serve-logo.png" alt="logo">

# Serve Client - Painless Docker Deployments

[![Build Status](https://travis-ci.com/loeffel-io/serve-server.svg?token=diwUYjrdo8kHiwiMCFuq&branch=master)](https://travis-ci.com/loeffel-io/serve-server)
![Project Status](https://www.repostatus.org/badges/latest/wip.svg)

## Examples

### Shared Resources
*.serve.yml*

```yaml
host: 'localhost:8080'
name: 'test-shared'

shared:
  networks:
    test-network:
```

### Apache httpd service
*.serve.yml*

```yaml
host: 'localhost:8080'
name: 'test-apache'

files:
  - index.html
  - images
  - Dockerfile

service:
  build:
    context: '%build_dir%'
    dockerfile: Dockerfile
  ports:
    - 80:80
  networks:
    - test-network
```

### MySQL Service
*.serve.yml*

```yaml
host: 'localhost:8080'
name: 'test-mysql'

service:
  image: 'mysql'
  command: --default-authentication-plugin=mysql_native_password
  restart: always
  environment:
    - MYSQL_ROOT_PASSWORD=example
  networks:
    - test-network
```

## Run

### MacOS

```bash
curl -sL -o serve https://github.com/loeffel-io/serve/releases/download/v0.2.0/serve-darwin && TOKEN="RANDOM-TOKEN-HERE" ./serve
```

### Linux

```bash
curl -sL -o serve https://github.com/loeffel-io/serve/releases/download/v0.2.0/serve-linux && TOKEN="RANDOM-TOKEN-HERE" ./serve
```