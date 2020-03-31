# Anti-Virus

This folder contains anti-virus database files that are compatible with ClamAV. The two files are:

- whitelist-files.fp: Files known to cause false positives when running ClamAV scans
- whitelist-signatures.ign2: Signatures known to cause false positives when running ClamAV scans

The files are generated using tool in `scripts/anti-virus-whitelists` and example usage is this:

```sh
export AV_DIR=$PWD
export AV_IGNORE_DIR=./anti-virus/
export AV_IGNORE_FILES=pkg/testdatagen/testdata/orders.pdf
export AV_IGNORE_SIGS="PUA.Pdf.Trojan.EmbeddedJavaScript-1 orders.pdf.UNOFFICIAL"
anti-virus-whitelists
```

The files should reside alongside the virus definition files downloaded by `freshclam` wherever that is run.
