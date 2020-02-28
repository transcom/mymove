# How to Test Prime API on staging and experimental

## Overview

In order to properly test the Prime API, you will need the ability to access it on staging and experimental.

## Requirements

You will first have to complete steps to create CAC [access](https://github.com/transcom/mymove/blob/master/docs/how-to/use-mtls-with-cac.md).

## Prerequisites

Download DoD certificates:

1. Go to this military CAC [website](https://militarycac.com/macnotes.htm#which_exact_CAC).
2. If you are a Safari or Chrome user scroll down to step 5.
3. If you are a Firefox user scroll down to step 5a.
4. Download the required files.
5. Confirm in Keychain on your Mac that you have all the certificates.
6. There will likely be a few certs that your Mac won't trust. You will need to manually enable `Always Trust` for these certificates.

## Testing Prime API

1. If testing on staging, run this command:

    ```sh
    go run ./cmd/prime-api-client --cac --hostname api.staging.move.mil --port 443
    ```

2. If testing on experimental, run this command:

    ```sh
    go run ./cmd/prime-api-client --cac --hostname api.experimental.move.mil --port 443
    ```

3. You will be prompted to enter your CAC pin. This will be the same six digit pin you created when picking up your CAC.

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
      }
    ]
    ```
