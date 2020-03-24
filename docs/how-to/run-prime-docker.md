# How to Run Prime Docker

This document aims to detail how to run the Prime Docker. The Prime Docker will be utilized by contractors to access and test the Prime API.

## Assumptions

- Install the following libraries:
  - golang
  - docker
  - docker-compose
  - direnv
  - jq
  - yarn
- If don't already have `.envrc.local`, copy `.envrc.local.template` to `.envrc.local`.
  - values will be added to this in a following story
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
rm -f bin/prime-api-client
make bin/prime-api-client
prime-api-client --insecure fetch-mtos
```

You should see something like this:

```{
  "createdAt": "2020-03-23",
  "id": "c66e2e16-4b3c-467b-a3a8-c80e46135dd2",
  "isAvailableToPrime": true,
  "isCanceled": false,
  "moveOrder": {
    "confirmationNumber": "GBQP4Q",
    "customer": {
      "branch": "COAST_GUARD",
      "currentAddress": {
        "city": "",
        "id": "00000000-0000-0000-0000-000000000000",
        "postalCode": "",
        "state": "",
        "streetAddress1": ""
      },
      ...
```

There will be more documentation on how to use the Prime API client soon.

When you're finished, remember to shut down the server:

```sh
docker-compose -f docker-compose.prime.yml down --remove-orphans
```
