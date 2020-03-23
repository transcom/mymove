# How to Run Prime Docker

This document aims to detail how to run the Prime Docker. The Prime Docker will be utilized by contractors to access and test the Prime API.

## Assumptions

- Installed everything the repo expects you to. (Covered in the repository's README.)
- Have the required env variables in `.envrc.local` (Will be covered in either this document in the future or a greater document that will point to this.)
- Modify `/etc/hosts` to include the prime host.

```sh
  echo "127.0.0.1 primelocal" | sudo tee -a /etc/hosts
```

## Running Prime Docker

Please make sure you're in the repository.

In the terminal, run the following:

```sh
make run_prime_docker
```

Please wait until the script is done running.

You should now be able to test the Prime API. You can do so with Postman or using the Prime API client within the terminal.

The latter can be quickly up and running with the following:

```sh
make bin/prime-api-client
prime-api-client --insecure fetch-mtos
```

There will be more documentation on how to use the Prime API client soon.

When you're finished, remember to shut down the server:

```sh
docker-compose -f docker-compose.prime.yml down --remove-orphans
```