# How to Test Prime API on staging and experimental

## Overview

In order to properly test the Prime API, you will need the ability to access it on staging and experimental.

## Requirements

You will first have to complete steps to create CAC [access](https://github.com/transcom/mymove/blob/master/docs/how-to/use-mtls-with-cac.md).


Additionally, those changes must be deployed to each environment. Merging to master will deploy to Staging. But you need to explicitly deploy to Experimental to get access there, otherwise you won't have access.

## Prerequisites

Download DoD certificates:

1. Go to this military CAC [website](https://militarycac.com/macnotes.htm#which_exact_CAC).
2. If you are a Safari or Chrome user scroll down to step 5.
3. If you are a Firefox user scroll down to step 5a.
4. Download the required files.
5. Confirm in Keychain on your Mac that you have all the certificates.
6. There will likely be a few certs that your Mac won't trust. You will need to manually enable `Always Trust` for these certificates.

## Sub-commands

At this time, there are only two sub-commands to be used within the Prime client:

- `fetch-mtos` which will fetch all the move task orders available to the Prime
- `update-mto-shipment` which will update an mto shipment with the data passed in
- `create-mto-service-item` which will create a new mto service item for a shipment

## Testing Prime API: Fetch MTOS

1. If testing on staging, run this command:

    ```sh
    go run ./cmd/prime-api-client --cac --hostname api.staging.move.mil --port 443 fetch-mtos | jq
    ```

2. If testing on experimental, run this command:

    ```sh
    go run ./cmd/prime-api-client --cac --hostname api.experimental.move.mil --port 443 fetch-mtos | jq
    ```

3. You will be prompted to enter your CAC pin. This will be the same pin you created when picking up your CAC.

4. If successful you should receive a response similar to:

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

1. If testing on staging, run this command:

    ```sh
    go run ./cmd/prime-api-client --cac --hostname api.staging.move.mil --port 443 update-mto-shipment --etag {ETAG} --filename {PATH TO FILE} | jq
    ```

2. If testing on experimental, run this command:

    ```sh
    go run ./cmd/prime-api-client --cac --hostname api.experimental.move.mil --port 443 fetch-mtos update-mto-shipment --etag {ETAG} --filename {PATH TO FILE} | jq
    ```

3. You will be prompted to enter your CAC pin. This will be the same pin you created when picking up your CAC.

4. If successful you should receive a response similar to:

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

  1. If testing on staging, run this command:

      ```sh
      go run ./cmd/prime-api-client --cac --hostname api.staging.move.mil --port 443 create-mto-service-item --filename {PATH TO FILE} | jq
      ```

  2. If testing on experimental, run this command:

      ```sh
      go run ./cmd/prime-api-client --cac --hostname api.experimental.move.mil --port 443 create-mto-service-item --filename {PATH TO FILE} | jq
      ```

  3. You will be prompted to enter your CAC pin. This will be the same pin you created when picking up your CAC.

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