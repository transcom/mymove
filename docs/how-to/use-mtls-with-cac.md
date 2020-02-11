# How to Upload Electronic Orders Using your CAC

## Requirements

The requirements are that you have a CAC and a CAC Reader. The recommended reader is the [Type C Smart Card Reader
Saicoo DOD Military USB-C Common Access Card Reader](https://www.amazon.com/Reader-Saicoo-Military-Compatible-Windows/dp/B071NT53M7/ref=sr_1_4).

To get going you will also need to install some software on your machine:

1. `brew install opensc` gives you tools like `pkcs11-tool` and `pkcs15-tool`

## MacOS 10.14 and earlier ONLY

**NOTE:** Skip this section if you are on MacOS `10.15.X` or later!!

Install the official Military CAC package from [CAC Key Packages](http://militarycac.org/MacVideos.htm#CACKey_packages) to get the driver `/usr/local/lib/pkcs11/cackey.dylib`

It's important that you disable default smart card access on your OSX machine. [Read the published instructions](http://militarycac.org/macuninstall.htm#Mojave_(10.14),_High_Sierra_(10.13.x),_and_Sierra_(10.12.x)_Built_in_Smart_Card_Ability). It boils down to removing your CAC from the reader, running this command and restarting your laptop:

```sh
sudo defaults write /Library/Preferences/com.apple.security.smartcard DisabledTokens -array com.apple.CryptoTokenKit.pivtoken
```

## Prerequisites

To see if you are in a place to use your CAC to extract certs you need to run these commands:

```sh
cac-prereqs
prereqs
```

## How it works

Mutual TLS is used in the app to authentication traffic. Users are authorized by comparing a SHA 256 hash of the
public certificate in the presented certificate, which is stored on the user's CAC.
To get that information into the system we must create a secure migration that contains both the
SHA 256 fingerprint and the Subject on the CAC certificate.

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

The generation script will provide three files and update the `migrations_manifest.txt` for you:

* A stub local secure migration in the `migrations/app/secure/` folder
* A migration to upload to AWS S3 in the `tmp/` folder

It is important only to upload the migration from the `tmp/` directory to the Staging and Experimental environments.
**DO NOT UPLOAD THIS MIGRATION TO PRODUCTION AS THERE SHOULD BE NO USE CASE FOR USING A PERSONAL CAC TO UPLOAD ORDERS
IN PROD**.

## For local testing only

For testing locally you can **temporarily** do this (your file names will be different):

```sh
cp tmp/20190822181328_cgilmer_cac.up.sql migrations/app/secure/
update-migrations-manifest
make db_dev_migrate
```

Validate your client certificate was updated (use your name):

```sh
psql-dev
> select * from client_certs where subject ILIKE '%chris%';
```

Now to test this with transcom/nom you need to enable the Mutual TLS listener and then run the server. To do so modify your `.envrc.local` with this content:

```sh
export MUTUAL_TLS_ENABLED=1
```

## Manually Generating a Secure Migration

**NOTE:**  Only follow these steps if you need to manually extract values from a CAC that you don't have physical access to.

To get the Fingerprint and Subject the user can run these commands:

```sh
cac-extract-fingerprint
cac-extract-subject
```

Now you need to generate the secure migration with these scripts:

```sh
FINGERPRINT=`cac-extract-fingerprint`
SUBJECT=`cac-extract-subject`
milmove gen certs-migration --name "${USER}_cac" -f "${FINGERPRINT}" -s "${SUBJECT}"
```

The output is the same as in the above steps.
