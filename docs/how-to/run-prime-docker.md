# How to Run Prime Docker

This document aims to detail how to run the Prime Docker. The Prime Docker will be utilized by contractors to access and test both the Prime API and the web applications.

## Assumptions

- Install the following libraries:
  - go
  - docker
  - docker-compose
  - direnv
  - jq
  - yarn
- If don't already have `.envrc.local`, copy `.envrc.local.template` to `.envrc.local`.
  - values will be added to this in a following story
- Modify `/etc/hosts` to include the prime, office, and milmove hosts.

```sh
  echo "127.0.0.1 primelocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 officelocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 milmovelocal" | sudo tee -a /etc/hosts
```

## Running Prime Docker

Please make sure you're in the `mymove` repository.

In your terminal, run the following:

```sh
make run_prime_docker
```

Please wait until the script is done running.

## Accessing Prime API

You should now be able to test the Prime API. You can do so with [Postman](make-a-sample-prime-api-call.md) or using the Prime API client within the terminal.

The latter can be quickly up and running with the following:

```sh
rm -f bin/prime-api-client
make bin/prime-api-client
prime-api-client --insecure fetch-mtos
```

You should see something like this:

```json
    [
        {
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
        },
        ...
    ]
```

There will be more documentation on how to use the Prime API client soon.

## Accessing Web Applications

While the container is still running, you should be able to access the different web applications

- [MilMove](http://milmovelocal:4000)
- [Office](http://officelocal:4000)

After the page loads in for either url, on the top right, you should be able to see a link for "Local Sign-In". This will show a list of users you can log in as. On the MilMove side, this will allow you to test as a service member at different stages within a move. Within the Office app, these different users represent different office user roles.

## Shutting Down Docker

When you're finished testing, remember to shut down the server:

```sh
make docker_compose_down
```
