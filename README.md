# logspout-unomaly
Unomaly adapter for Logspout. 

The adapter forwards all your docker logs to unomaly, using the container name as system.

# Building

To build the Unomaly Logspout Docker image, run:
* `make docker`

# Configuration and invocation

This module can be configured either by setting environment variables in
Docker, or by using the [Logspout routesapi](https://github.com/gliderlabs/logspout/tree/master/routesapi). The following variables are available:

Env. Variable | routesapi key | Type | Required? | Description |
| --- | --- | --- | --- | -----|
| `UNOMALY_INGESTION` | `ingestionHost` | string | required | Unomaly host URL |

# How to use

    docker run \
        -e "ROUTE_URIS=unomaly://localhost" \
        -e "UNOMALY_INGESTION=<UNOMALY_HOST>" \
        --volume=/var/run/docker.sock:/var/run/docker.sock \
        logspout-unomaly
