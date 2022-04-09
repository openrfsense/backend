# OpenRFSense Backend
The backend is responsible for:
- controlling and querying the sensors remotely via MQTT
- managing encryption keys and TLS certificates for secure communication with the sensors
- exposing a RESTful API to fetch data and provide an abstraction layer over the sensors control

## Table of contents <!-- omit in toc -->
- [OpenRFSense Backend](#openrfsense-backend)
    - [Usage and deployment](#usage-and-deployment)
    - [Configuration](#configuration)
      - [YAML](#yaml)
      - [Environment variables](#environment-variables)
    - [API](#api)

### Usage and deployment
Docker has become the industry standard nowadays, so an official Docker image and an example `docker-compose.yml` are provided. The basic setup requires:
- the OpenRFSense backend, of course
- the [Emitter MQTT broker](https://emitter.io/download/)

A reverse proxy is also recommended: something like [Caddy](https://caddyserver.com/) is simple to use and has automatic TLS among other nice features.

To run the Docker stack, first get yourself a license and secret key by running Emitter without any configuration:
```shell
$ docker run --rm emitter/server
```

which will output something along the lines of:
```shell
[service] unable to find a license, make sure 'license' value is set in the config file or EMITTER_LICENSE environment variable
[license] generated new license: <license>
[license] generated new secret key: <secret key>
```

Finally, copy the license and secret key and run the following command:

```shell
$ env EMITTER_LICENSE=<license> EMITTER_SECRET=<secret key> docker-compose up -d
```

You can optionally set the `EMITTER_TLS_HOST=<your domain where Emitter is reachable>` and `EMITTER_TLS_EMAIL=<your email>` variables to have Emitter fetch certificates and use TLS automatically for all connections. For anything, you can modify `emitter.conf` or load your own Emitter configuration file in the broker container (see Emitter documentation for that).

### Configuration
> ⚠️ The configuration is still WIP: keys may change in the future

The `config` module from [`openrfsense/common`](https://github.com/openrfsense/common) is used. As such, configuration values are loaded from a YAML file first, then from environment variables.

#### YAML
See the example [`config.yml`](./config.yml) file for now, as the configuration is very prone to change.

#### Environment variables
Environment variables are defined as follows: `ORFS_SECTION_SUBSECTION_KEY=value`. They are loaded after any other configuration file, so they cam be used to overwrite any configuration value.

For example, it's recommended to set `mqtt.secret` (the MQTT broker secret key, see [Emitter documentation](https://emitter.io/develop/getting-started/) under **Channel Security**) as an envionment variable: `ORFS_MQTT_SECRET=key`.

### API
The backend will automatically generate its own [Swagger](https://swagger.io/) documentation and serve a webpage with [Swagger UI](https://swagger.io/tools/swagger-ui/) at `https://$DOMAIN/swagger`. The common JSON objects are defined as Golang structs in [`openrfsense/common.types`](https://github.com/openrfsense/common).
