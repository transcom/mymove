# Run server_test job in CircleCI container locally

The server_test job seems to fail in some cases only when running on the build server. This document will walk you through running the job inside a CircleCI container on your local machine.

## Running the tests

To run the job execute the following command:

```sh
docker-compose -f docker-compose.circle.yml --compatibility up server_test
```

The `--compatibility` flag must be used to tell docker-compose to apply the limits to the local containers even though we are not running in a docker swarm. As of version 3 of docker-compose these settings are only applied to swarm containers. See [Upgrading Version 2 to 3](https://docs.docker.com/compose/compose-file/compose-versioning/#upgrading)

### Modify the memory and CPU constraints

To adjust the memory and/or CPU constraints modify the deploy section in the `docker-compose.circle.yml` file as desired. See the [documentation](https://docs.docker.com/compose/compose-file/#resources) for more details.

```yaml
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1024M
```

## Caveat

These instructions make use of your local files and you will need to run `make clean` afterwards to have your local setup work as expected.
