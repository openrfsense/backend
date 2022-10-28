# OpenRFSense Backend
The backend is responsible for:
- controlling and querying the sensors remotely via NATS
- managing encryption keys and TLS certificates for secure communication with the sensors
- exposing a RESTful API to fetch data and provide an abstraction layer over the sensors control

## Table of contents <!-- omit in toc -->
- [OpenRFSense Backend](#openrfsense-backend)
    - [Usage and deployment](#usage-and-deployment)
    - [Configuration](#configuration)
    - [API](#api)
    - [Metrics](#metrics)

### Usage and deployment
Docker has become the industry standard nowadays, so an official Docker image is provided as part of this project. The basic setup requires only the one image, but an external NATS server can be provided (`WIP`, see configuration).

An external database is required for data storage (measurement campaigns and general sample retrieval). Only PostgreSQL is supported and it's easy enough to deploy (see example below). See also default configuration file for more information on how to configure the database connection.

A reverse proxy is also recommended: something like [Caddy](https://caddyserver.com/) is simple to use and has automatic TLS among other nice features.

Docker Compose can also be used to manage your infrastructure with Caddy. An example `docker-compose.yml` would look like this:
> Please not that this is just **an example** of configuration for an external project and, as such, not guaranteed to be future-proof or maintained as-is.
<!-- TODO: add TLS docs for embedded NATS -->
```yaml
services:
  caddy:
    image: caddy:2-alpine
    container_name: caddy
    restart: always
    ports:
      - "80:80"
      - "443:443"
    networks:
      - proxy
    volumes:
      - ./caddy_data:/data
      - ./caddy_config:/config
    command: ["caddy", "reverse-proxy", "--to", "orfs_backend:8081"]

  postgres:
    image: postgres:alpine
    container_name: postgres
    restart: always
    ports:
      - "5432:5432"
    environment:
      POSTGRES_PASSWORD: postgres
  
  openrfsense_backend:
    image: openrfsense-backend:latest
    container_name: orfs_backend
    depends_on:
      - caddy
    networks:
      - db
      - proxy
    environment:
      ORFS_POSTGRES_HOST: postgres
      ORFS_POSTGRES_USERNAME: postgres
      ORFS_POSTGRES_PASSWORD: postgres
      ORFS_POSTGRES_DBNAME: postgres

networks:
  proxy:
  db:
```

### Configuration
> ⚠️ The configuration is still WIP: keys may change in the future

Configuration values are loaded from a YAML file first, then from environment variables.

#### YAML
See the example [`config.yml`](./config.yml) file for now, as the configuration is very prone to change.

#### Environment variables
Environment variables are defined as follows: `ORFS_SECTION_SUBSECTION_KEY=value`. They are loaded after any other configuration file, so they cam be used to overwrite any configuration value.

For example, it's recommended to set `nats.token` (the NATS server access token for API calls, see [NATS documentation](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/tokens)) as an envionment variable: `ORFS_NATS_TOKEN=token`.

### API
The backend will automatically generate its own [Swagger](https://swagger.io/) documentation and serve a webpage with [Swagger UI](https://swagger.io/tools/swagger-ui/) at `https://$DOMAIN/api/docs`. The common JSON objects are defined as Golang structs in [`openrfsense/common.types`](https://github.com/openrfsense/common).

### Metrics
Metrics are served at `https://$DOMAIN/metrics` if enabled in the configuration (defined by the value of `backend.metrics`, see default configuration). A simple, dynamic web page is shown by default but the metrics can also be retrieved in JSON format by sending `Accept: application/json` along with the request. For more information, see the [Monitor middleware for Fiber](https://docs.gofiber.io/api/middleware/monitor).
