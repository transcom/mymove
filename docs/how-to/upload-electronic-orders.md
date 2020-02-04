# How to Upload Electronic Orders Using your CAC

## Requirements

You must first read the article on [How to Use mTLS with CAC](./use-mtls-with-cac.md).

## Uploading Electronic Orders Locally

Use the transcom/nom repo with [sample navy orders data](https://drive.google.com/drive/folders/1dxOO9uXSOWfjQiKMzwX3bmRqBJfBLldi). It's important that the SSNs match the ones in the DMDC Contractor Test database. You can see the [set of contractor test SSN's](https://drive.google.com/file/d/1vfxEaC6cadFtMlTGFZsy95P52poKLaXA/view).

Also the DMDC Host needs to be set correctly in the `.envrc` file:

```sh
export IWS_RBS_HOST="pkict.dmdc.osd.mil"
```

Then run:

```sh
direnv allow
make server_run
```

To continue you need to get the Token from the CAC with a script in transcom/mymove (Becomes `ENTERYOURTOKEN` in following step):

```sh
cac-extract-token-label
```

Now over in your git checkout of the transcom/nom repo. Then download the [sample csv](https://drive.google.com/open?id=1-zxetfRhLEpnx1SBTAveoTLpwEzp3fK-) into the repo. And run these commands (**NOTE:** you will need your CAC personal PIN to do this operation):

For MacOS 10.14 and earlier:

```sh
make bin/nom
TOKEN="ENTERYOURTOKEN"
MODULE="/usr/local/lib/pkcs11/cackey.dylib"
bin/nom -host orderslocal -port 9443 -insecure -pkcs11module "${MODULE}" -certlabel "Identity #0" -keylabel "Identity #0" --tokenlabel "${TOKEN}" nom_demo_20190404.csv
PIN: ********
```

For MacOS 10.15 and later

```sh
make bin/nom
TOKEN="ENTERYOURTOKEN"
MODULE="/usr/local/lib/pkcs11/opensc-pkcs11.so"
bin/nom -host orderslocal -port 9443 -insecure -pkcs11module "${MODULE}" nom_demo_20190404.csv
PIN: ********
```

## Updating the Sample CSV for transcom/nom

The data in transcom/nom `sample.csv` is generated from data in the fake records hosted by the DMDC. Copies of
the fake data exist in CSV/Excel files in the [USTC MilMove -> Integrations -> Identity Web Services -> Developer Samples](https://drive.google.com/drive/folders/16k7eG4j5vSBQIX_eTWnoXqiae1T0ysiq) folder. The latest set of data is [Cust2675_TRANSCOM_20190823_Demo2](https://drive.google.com/drive/folders/16k7eG4j5vSBQIX_eTWnoXqiae1T0ysiq). If you need to update
this data you will need to contact DMDC as they refresh the data from time to time.

Tip: If you want to skip the DMDC lookup, put 10 digit numbers (starting with any digit other than 0) instead of 9 digit numbers in the SSN column in the CSV. It's still worthwhile to test the whole pipeline, obviously, but sometimes the DMDC part is more trouble than it's worth.
