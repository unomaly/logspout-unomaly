# logspout-unomaly
Unomaly adapter for Logspout. More documentation can be found [in Unomaly docs](https://unomaly.io/docs/connect/logspout/).

Expects to ingest JSON log lines, and will send JSON blobs up to Unomaly, annotated with the current logspout stream, container, container ID, hostname, and docker image name.

If the log lines being streamed through Logspout aren't JSON, the contents of the message will be tucked under a `"message"` key in the Unomaly payload, alongside the metadata mentioned above.

# Building

To build the Unomaly Logspout Docker image, run:
* `make docker`

# Configuration and invocation

This module can be configured either by setting environment variables in
Docker, or by using the [Logspout routesapi](https://github.com/gliderlabs/logspout/tree/master/routesapi). The following variables are available:

Env. Variable | routesapi key | Type | Required? | Description |
| --- | --- | --- | --- | -----|
| `UNOMALY_WRITE_KEY` | `writeKey` | string | required | Your Unomaly team's write key. |
| `UNOMALY_DATASET` | `dataset` | string | required | The name of the destination dataset in your Unomaly account. It will be created if it does not already exist. |
| `UNOMALY_SAMPLE_RATE` | `sampleRate` | integer | optional | Sample your event stream: send 1 out of every N events |

### Environment variables

Configure the logspout-unomaly image via environment variables and run the container:

    docker run \
        -e "ROUTE_URIS=unomaly://localhost" \
        -e "UNOMALY_WRITE_KEY=<YOUR_WRITE_KEY>" \
        -e "UNOMALY_DATASET=<YOUR_DATASET>" \
        --volume=/var/run/docker.sock:/var/run/docker.sock \
        --publish=127.0.0.1:8000:80 \
        unomalyio/logspout-unomaly:1.13

### routesapi

Configuration can be set after the logspout-unomaly image is already running via routesapi:

    docker run \
        --volume=/var/run/docker.sock:/var/run/docker.sock \
        --publish=127.0.0.1:8000:80 \
        unomalyio/logspout-unomaly:1.13

    curl $(docker port `docker ps -lq` 80)/routes \
        -X POST \
            -d '{"adapter": "unomaly",
                 "address": "unomaly://localhost",
                 "options": {"writeKey":"<YOUR_WRITE_KEY>",
                             "dataset":"<YOUR_DATASET>"}}'
