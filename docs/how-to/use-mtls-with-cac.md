# How to Create CAC Access (for using Prime API and uploading Electronic Orders)

## Overview

Mutual TLS is used in the app to authenticate traffic. Users are authorized by comparing a SHA 256 hash of the
public certificate in the presented certificate, which is stored on the user's CAC.
To get that information into the system we must create a [secure migration](https://github.com/transcom/mymove/blob/master/docs/database/migrate-the-database.md#secure-migrations) that contains both the
SHA 256 fingerprint and the Subject on the CAC certificate.

## Requirements

The requirements are that you have a CAC and a CAC Reader. The recommended reader is the [Type C Smart Card Reader
Saicoo DOD Military USB-C Common Access Card Reader](https://www.amazon.com/Reader-Saicoo-Military-Compatible-Windows/dp/B071NT53M7/ref=sr_1_4).

To get going, install software on your machine with this command:

 `brew install opensc`

 It gives you tools like `pkcs11-tool` and `pkcs15-tool`.  Note: with these tools, you will not need to install a driver for the card reader.

### MacOS 10.14 and earlier ONLY

**NOTE:** Skip this section if you are on MacOS `10.15.X` or later!!

Install the official Military CAC package from [CAC Key Packages](http://militarycac.org/MacVideos.htm#CACKey_packages) to get the driver `/usr/local/lib/pkcs11/cackey.dylib`

It's important that you disable default smart card access on your OSX machine. [Read the published instructions](http://militarycac.org/macuninstall.htm#Mojave_(10.14),_High_Sierra_(10.13.x),_and_Sierra_(10.12.x)_Built_in_Smart_Card_Ability). It boils down to removing your CAC from the reader, running this command and restarting your laptop:

```sh
sudo defaults write /Library/Preferences/com.apple.security.smartcard DisabledTokens -array com.apple.CryptoTokenKit.pivtoken
```

## Prerequisites

1. To see if you are in a place to use your CAC to extract certs you need to run these commands:

    ```sh
      cac-prereqs
      prereqs
    ```

2. Plug in your card reader to your computer

3. Insert your CAC into the reader

## Generating a Secure Migration

To generate the secure migration run this step:

```sh
milmove gen certs-migration --name "${USER}_cac" --cac
```

You will see output like:

```text
2019/08/22 17:32:38 new migration file created at: "tmp/20190822173238_cgilmer_cac.up.sql"
2019/08/22 17:32:38 new migration file created at:  "migrations/app/secure/20190822173238_cgilmer_cac.up.sql"
2019/08/22 17:32:38 new migration appended to manifest at: "/dir/transcom/mymove/migrations/app/migrations_manifest.txt"
```

The generation script will do the following:

* Create a stub local secure migration in the `migrations/app/secure/` folder
* Create a migration to upload to AWS S3 in the `tmp/` folder
* Update the `migrations_manifest.txt`

## Preparing for Upload

You will be uploading the migration from the `tmp/` directory, ONLY to the Staging and Experimental environments.

**DO NOT UPLOAD THIS MIGRATION TO PRODUCTION AS THERE SHOULD BE NO USE CASE FOR USING A PERSONAL CAC TO UPLOAD ORDERS IN PROD**.

Before uploading we'll remove some sensitive information and verify the migration script runs locally.

### Testing Certificate Locally

For testing locally:

1. copy the `tmp` migration file contents to the corresponding file in `migrations/app/secure/` (your file name will be different):

    ```sh
    cp tmp/20190822181328_cgilmer_cac.up.sql migrations/app/secure/
    ```

2. Change the migration sql file in `migrations/app/secure/` (which is for Devlocal only) by updating the `CN` field with your GitHub username.
This will prevent your name and EDIPI from getting checked into git.

    Example:

    ```tex
    CN=LAST.FIRST.MI.EDIPI,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US
    # should be
    CN=GITHUB_USERNAME,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US
    ```

    Note: After your secure migration of `tmp` file, staging and experimental migration files will continue to have your name and EDIPI in the format `CN=LAST.FIRST.MIDDLE.EDIPI,...`).

3. Test this by running:

    ```bash
    make db_dev_migrate
    ```

4. Make a PR with the secure migration sql file located in `migrations/app/secure` and the `migrations_manifest.txt` file.

5. Validate your client certificate was updated (use your username):

    ```sh
    psql-dev
    > select * from client_certs where subject ILIKE '%GITHUB_USERNAME%';
    ```

6. To test this with transcom/nom, run the server and test the secure migration (with your CAC inserted into your reader):

    ```sh
    make db_dev_e2e_populate server_run
    go run ./cmd/prime-api-client/main.go --insecure --cac | jq .
    ```

    If the secure migration worked you should receive a response similar to

    ```json
    [
      {
        "createdAt": "2020-01-22",
        "id": "5d4b25bb-eb04-4c03-9a81-ee0398cb7791",
        "isAvailableToPrime": true,
        "isCanceled": false,
        "moveOrderID": "6fca843a-a87e-4752-b454-0fac67aa4981",
        "mto_service_items": [],
        "payment_requests": [],
        "mto_shipments": [
          {
            ...
          }
        ],
        "updatedAt": "2020-01-22"
      }
    ]
    ```



## Uploading CAC Secure Migration

Once the files have been generated by previous section next step is to upload them to the environments.

1. Upload to `experimental` and `staging` using the file in `tmp/`:

    ```bash
    ENVIRONMENTS="experimental staging" ./scripts/upload-secure-migration tmp/20200211150405_mr337_cac.up.sql
    ```

2. For `prod` **do not** upload the `tmp` file. Instead we'll upload a stub file as a place holder since the migration process will be looking for such a file in `prod` environment.

    Remove the contents of your `migrations/app/secure` file (given you've already committed it) and then paste `-- This is a stub file for user X`

    OR

    Use this command:

    ```bash
    echo "-- This is a stub file for user X" >> migrations/app/secure/20200211150405_mr337_cac.up.sql
    echo "-- This is a stub file for user X" >> tmp/20200211150405_mr337_cac.up.sql
    ENVIRONMENTS=prod ./scripts/upload-secure-migration tmp/20200211150405_mr337_cac.up.sql
    ```

    The file should be empty except for the comment `-- This is a stub file for user X`

3. Once completed confirm the upload with the command `download-secure-migration`

    * Ensure the secure migration files for `experimental` and `staging` inserts a new record into the database with sensitive information.
    * Ensure the secure migration file for `prod` is just a stub and does not insert a record

    ```bash
    download-secure-migration 20200211150405_mr337_cac.up.sql
    Downloading from: experimental
    ...executing: aws s3 cp s3://transcom-ppp-app-experimental-us-west-2/secure-migrations/20200211150405_mr337_cac.up.sql ./tmp/secure_migrations/experimental/20200211150405_mr337_cac.up.sql
    Downloading from: staging
    ...executing: aws s3 cp s3://transcom-ppp-app-staging-us-west-2/secure-migrations/20200211150405_mr337_cac.up.sql ./tmp/secure_migrations/staging/20200211150405_mr337_cac.up.sql
    Downloading from: prod
    ...executing: aws s3 cp s3://transcom-ppp-app-prod-us-west-2/secure-migrations/20200211150405_mr337_cac.up.sql ./tmp/secure_migrations/prod/20200211150405_mr337_cac.up.sql

    Files have been downloaded to these locations:

    ./tmp/secure_migrations/prod/20200211150405_mr337_cac.up.sql
    ./tmp/secure_migrations/experimental/20200211150405_mr337_cac.up.sql
    ./tmp/secure_migrations/staging/20200211150405_mr337_cac.up.sql

    Please remember to 'rm -rf ./tmp/secure_migrations' when you are finished working
    ```



## Manually Generating a Secure Migration

**NOTE:**  Only follow these steps if you need to manually extract values from a CAC that you don't have physical access to.

1. To get the Fingerprint and Subject the user can run these commands:

    ```sh
    cac-extract-fingerprint
    cac-extract-subject
    ```

2. Now you need to generate the secure migration with these scripts:

    ```sh
    FINGERPRINT=`cac-extract-fingerprint`
    SUBJECT=`cac-extract-subject`
    milmove gen certs-migration --name "${USER}_cac" -f "${FINGERPRINT}" -s "${SUBJECT}"
    ```

    The output is the same as in the above steps.
