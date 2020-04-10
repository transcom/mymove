# How to Test Prime API locally

## Overview

This document details how to test the Prime API locally.

For how to test on staging and experimental, follow this [link](test-prime-api-staging-experimental.md).

## Requirements

You must have data generated within your database and have the server running.

If you are using the [Prime Docker](run-prime-docker.md) via `make run_prime_docker`, this has already been done for you.

If you are not, please generate the data through `make db_dev_e2e_populate && server_run`,

## Sub-commands

These are the sub-commands that can be used within the Prime client:

- `fetch-mtos` which will fetch all the move task orders available to the Prime
- `update-mto-shipment` which will update an mto shipment with the data passed in
- `create-mto-service-item` which will create a new mto service item for a shipment
- `make-available-to-prime` which will make an MTO available to the prime. This uses the support api to get an MTO in the position to be used to test prime api more effectively

## Testing Prime API: Fetch MTOS

1. Run this command:

   ```sh
   go run ./cmd/prime-api-client --insecure fetch-mtos | jq
   ```

2. If successful you should receive a response similar to:

   ```json
   [
     {
       "createdAt": "2020-01-22",
       "id": "5d4b25bb-eb04-4c03-9a81-ee0398cb7791",
       "isAvailableToPrime": true,
       "isCanceled": false,
       "moveOrderID": "6fca843a-a87e-4752-b454-0fac67aa4981",
       "mto_service_items": [{ ... }],
       "mto_shipments": [{ ... }],
       "payment_requests": [],
       "reference_id": "1234-4321",
       "updatedAt": "2020-01-22"
     },
     {...},
   ]
   ```

## Testing Prime API: Update MTO Shipment

Before updating a shipment, you need to figure out the ID, the MTO ID, and the ETag of the shipment you'll be updating. You can accomplish this by using the `fetch-mtos` sub-command. You must also have data to pass in the form of a JSON file. Data must have at least the ID and MTO ID. Additionally, any field you are attempting to update should be part of this data. For example:

```json
{
  "id": "df2725cf-ee88-4fe5-90e7-a693b551dd3a",
  "moveTaskOrderID": "5d4b25bb-eb04-4c03-9a81-ee0398cb7791",
  "actualPickupDate": "2020-03-25"
}
```

1. Run this command:

   ```sh
   go run ./cmd/prime-api-client --insecure update-mto-shipment --etag {ETAG} --filename {PATH TO FILE} | jq
   ```

2. If successful you should receive a response similar to:

   ```json
        {
           "actualPickupDate": "2020-03-25",
           "agents": [
             {
               "agentType": "RELEASING_AGENT",
               "createdAt": "0001-01-01",
               "email": "test@test.email.com",
               "firstName": "Test",
               "id": "b870fe0c-f2a1-4372-9d44-eaaf54a1b45d",
               "lastName": "Agent",
               "mtoShipmentID": "df2725cf-ee88-4fe5-90e7-a693b551dd3a",
               "updatedAt": "0001-01-01"
             },
             {...}
           ],
           "createdAt": "2020-03-25T15:42:28.992Z",
           "customerRemarks": "please treat gently",
           "destinationAddress": {...},
           "eTag": "MjAyMC0wMy0yNVQxNTo0MjoyOC45OTI1ODFa",
           "id": "df2725cf-ee88-4fe5-90e7-a693b551dd3a",
           "moveTaskOrderID": "5d4b25bb-eb04-4c03-9a81-ee0398cb7791",
            ...
         }
   ```

## Testing Prime API: Create MTO Service Item

Before creating a new MTO service item, you need to figure out the ID and the MTO ID of the shipment you'll be updating. You can accomplish this by using the `fetch-mtos` sub-command. You must also have data to pass in the form of a JSON file. Data must have at least the ID and MTO ID plus the Model Type, Pickup Postal Code, and Reason for the new service item. For example:

```json
{
  "modelType": "MTOServiceItemDOFSIT",
  "moveTaskOrderID": "5d4b25bb-eb04-4c03-9a81-ee0398cb7791",
  "mtoShipmentID": "df2725cf-ee88-4fe5-90e7-a693b551dd3a",
  "pickupPostalCode": "90802",
  "reason": "sometimes stuff happens"
}
```

1. Run this command:

   ```sh
   go run ./cmd/prime-api-client --insecure create-mto-service-item --filename {PATH TO FILE} | jq
   ```

2. If successful you should receive a response similar to:

   ```json
   {
     "pickupPostalCode": "90802",
     "reServiceCode": "DOFSIT",
     "reason": "sometimes stuff happens",
     "eTag": "MjAyMC0wMy0yNVQxNjowMTowMS41ODg0MTha",
     "id": "10873d77-b617-40a3-a66d-3b8f4389ba5e",
     "modelType": "MTOServiceItemDOFSIT",
     "moveTaskOrderID": "5d4b25bb-eb04-4c03-9a81-ee0398cb7791",
     "reServiceID": "998beda7-e390-4a83-b15e-578a24326937"
   }
   ```

3. Run this command:

   ```sh
   go run ./cmd/prime-api-client --insecure make-available-to-prime --filename {PATH TO FILE} --etag
   ```

4. If successful you should receive a response similar to:

   ```json
   {
     "pickupPostalCode": "90802",
     "reServiceCode": "DOFSIT",
     "reason": "sometimes stuff happens",
     "eTag": "MjAyMC0wMy0yNVQxNjowMTowMS41ODg0MTha",
     "id": "10873d77-b617-40a3-a66d-3b8f4389ba5e",
     "modelType": "MTOServiceItemDOFSIT",
     "moveTaskOrderID": "5d4b25bb-eb04-4c03-9a81-ee0398cb7791",
     "reServiceID": "998beda7-e390-4a83-b15e-578a24326937"
   }
   ```

## Testing Prime API (Support API): Make MTO available to prime

To make an MTO available to the prime you must have a move order with an MTO associated with it and that MTO is not currently available to the prime. ID is the only required field in the file. The order of the flag usage does not matter.

Examples usage:
`data.json`

```json
{
  "id": "6efa0519-d4cf-4c92-9eac-a8f66479a1a6"
}
```

  1. Run this command:

   ```sh
    go run ./cmd/prime-api-client make-available-to-prime --etag {Etag of MTO} --filename {PATH TO FILE} --insecure | jq
   ```

   `Example usage:`

   ```sh
    prime-api-client make-available-to-prime --etag MjAyMC0wMy0zMVQyMTo1NjozOS45NzIxNzha --filename data.json --insecure | jq
   ```

  `Successful response:`

  ```json
  {
    "createdAt": "2020-03-31",
    "eTag": "MjAyMC0wNC0wN1QxOToxMzowNi4xMzMzNjJa",
    "id": "69a60e8f-5eef-42e7-8de6-5e1bb8a2fa53",
    "isAvailableToPrime": true,
    "isCanceled": false,
    "moveOrderID": "74fc2ab0-2f5b-466f-915a-23d33da712c4",
    "referenceId": "6287-6043",
    "requestedPickupDate": "0001-01-01",
    "updatedAt": "2020-04-07"
  }
  ```

## Quick Tricks

You've seen this at the beginning of each command:

```sh
go run ./cmd/prime-api-client
```

You can shorten it by doing the following

```sh
rm -f bin/prime-api-client
make bin/prime-api-client
```

Then every time you see `go run ./cmd/prime-api-client` you can replace it with simply `prime-api-client`. For example:

```shell script
prime-api-client --insecure fetch-mtos
```
