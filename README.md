# WIP: Serve Client - Painless Docker Deployments

[![Build Status](https://travis-ci.com/loeffel-io/serve.svg?token=diwUYjrdo8kHiwiMCFuq&branch=master)](https://travis-ci.com/loeffel-io/serve)

## Run

```bash
curl ..
```

## Example

### .serve.yml

```yaml
host: 'localhost:8080'
name: 'serve-project'

files:
  - index.html
  - Dockerfile

service:
  build:
    context: "%build_dir%"
    dockerfile: Dockerfile
  ports:
    - 3000:80
```

