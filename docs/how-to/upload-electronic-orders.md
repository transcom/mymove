# How to Upload Electronic Orders Using your CAC

## Requirements

The requirements are that you have a CAC and a CAC Reader. The recommended reader is the [Type C Smart Card Reader
Saicoo DOD Military USB-C Common Access Card Reader](https://www.amazon.com/Reader-Saicoo-Military-Compatible-Windows/dp/B071NT53M7/ref=sr_1_4).

To get going you will also need to install some software on your machine:

1. `brew install opensc` gives you tools like `pkcs11-tool` and `pkcs15-tool`
1. Install the official  Military CAC package from [CAC Key Packages](http://militarycac.org/MacVideos.htm#CACKey_packages) to get the driver `/usr/local/lib/pkcs11/cackey.dylib`

It's important that you disable default smart card access on your OSX machine. [Read the published instructions](http://militarycac.org/macuninstall.htm#Mojave_(10.14),_High_Sierra_(10.13.x),_and_Sierra_(10.12.x)_Built_in_Smart_Card_Ability). It boils down to removing your CAC from the reader, running this command and restarting your laptop:

```sh
sudo defaults write /Library/Preferences/com.apple.security.smartcard DisabledTokens -array com.apple.CryptoTokenKit.pivtoken
```

This should get you in the state where you can use your CAC to extract certs. You may also want to run these commands
to check:

```sh
cac-prereqs
prereqs
```

## Generating a Secure Migration

Orders does mutual TLS authentication and then authorizes you by comparing a SHA 256 hash of your public certificate
stored on your CAC. To get that information into the system you have to create a secure migration using your
SHA 256 fingerprint and the Subject on your CAC certificate.

To get the Fingerprint and Subject you can run these commands:

```sh
cac-extract-fingerprint
cac-extract-subject
```

Now you need to generate the secure migration with these scripts:

```sh
FINGERPRINT=`cac-extract-fingerprint`
SUBJECT=`cac-extract-subject`
milmove gen orders-migration --name "${USER}_cac" -f "${FINGERPRINT}" -s "${SUBJECT}"
```

You will see output like:

```text
2019/08/22 17:32:38 new migration file created at: "tmp/20190822173238_cgilmer_cac.up.sql"
2019/08/22 17:32:38 new migration file created at:  "local_migrations/20190822173238_cgilmer_cac.up.sql"
2019/08/22 17:32:38 new migration appended to manifest at: "/dir/transcom/mymove/migrations_manifest.txt"
```

The generation script will provide three files and update the `migrations_manifest.txt` for you:

* A stub local migration in the `local_migrations/` folder
* A migration to upload to AWS S3 in the `tmp/` folder

It is important only to upload the migration from the `tmp/` directory to the Staging and Experimental environments.
**DO NOT UPLOAD THIS MIGRATION TO PRODUCTION AS THERE SHOULD BE NO USE CASE FOR USING A PERSONAL CAC TO UPLOAD ORDERS
IN PROD**.

## Uploading Electronic Orders Locally

Use the transcom/nom repo with [sample navy orders data](https://drive.google.com/drive/folders/1dxOO9uXSOWfjQiKMzwX3bmRqBJfBLldi). It's important that the SSNs match the ones in the DMDC Contractor Test database. You can see the [set of contractor test SSN's](https://drive.google.com/file/d/1vfxEaC6cadFtMlTGFZsy95P52poKLaXA/view).

For testing locally you can **temporarily** do this (your file names will be different):

```sh
rm local_migrations/20190822181328_cgilmer_cac.up.sql
mv tmp/20190822181328_cgilmer_cac.up.sql migrations/
update-migrations-manifest
make db_dev_migrate
```

Validate your client certificate was updated (use your name):

```sh
psql-dev
> select * from client_certs where subject ILIKE 'chris';
```

Now to test this with transcom/nom you need to enable the Mutual TLS listener and then run the server. To do so modify your `.envrc.local` with this content:

```sh
export MUTUAL_TLS_ENABLED=1
```

Also the DMDC Host needs to be set correctly in the `.envrc` file:

```sh
export IWS_RBS_HOST="pkict.dmdc.osd.mil"
```

Then run:

```sh
direnv allow
make server_run
```

To continue you need to get the Token from the CAC with a script in transcom/mymove:

```sh
cac-extract-token-label
```

Grab the `token label` from `Slot 0`. It should have your name and a number in it.

Now over in your git checkout of the transcom/nom repo. Then download the [sample csv](https://drive.google.com/open?id=1-zxetfRhLEpnx1SBTAveoTLpwEzp3fK-) into the repo. And run these commands (**NOTE:** you will need your CAC personal PIN to do this operation):

```sh
make bin/nom
TOKEN=`ENTERYOURTOKEN`
MODULE=`/usr/local/lib/pkcs11/cackey.dylib`
bin/nom -host orderslocal -port 9443 -insecure -pkcs11module "${MODULE}" --tokenlabel "${TOKEN}" nom_demo_20190404.csv
PIN: ********
```

## Updating the Sample CSV for transcom/nom

The data in transcom/nom `sample.csv` is generated from data in the fake records hosted by the DMDC. Copies of
the fake data exist in CSV/Excel files in the [USTC MilMove -> Integrations -> Identity Web Services -> Developer Samples](https://drive.google.com/drive/folders/16k7eG4j5vSBQIX_eTWnoXqiae1T0ysiq) folder. The latest set of data is [Cust2675_TRANSCOM_20190823_Demo2](https://drive.google.com/drive/folders/16k7eG4j5vSBQIX_eTWnoXqiae1T0ysiq). If you need to update
this data you will need to contact DMDC as they refresh the data from time to time.
