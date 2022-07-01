# TLS certificates

## Certificates in this Directory

A description of the certificates in this directory will helpful:

| Name | Function |
| --- | --- |
| `api.demo.dp3.us.cert` | Certificate for api.demo.dp3.us |
| `api.demo.dp3.us.p7b` | Certificate chain for api.demo.dp3.us (non-ATO) |
| `api.loadtest.dp3.us.crt` | Certificate for api.loadtest.dp3.us |
| `api.loadtest.dp3.us.chain.der.p7b` | Certificate chain for api.loadtest.dp3.us (non-ATO) |
| `Certificates_PKCS7_v5.6_DoD.der.p7b` | |
| `Certificates_PKCS7_v5.11_WCF.pem.p7b` | The collection of certs from which the dod-wcf-* certificates are derived. |
| `devlocal-ca.key` | Devlocal CA Key |
| `devlocal-ca.pem` | Devlocal CA PEM |
| `devlocal-ca.srl` | Devlocal CA Serial |
| `devlocal-client_auth_secret.key` | Client auth secret JWT key. |
| `devlocal-faux-(air-force/all/army-hrc/coast-guard/marine-corps/navy)-orders.(cer/key)` | Certs signed by Devlocal CA for Orders API testing |
| `devlocal-https.(key/pem)` | a self-signed TLS cert/key pair |
| `devlocal-mtls.(cer/key)` | Certs signed by Devlocal CA for mTLS testing |
| `dod-sw-ca-54.pem` | DoD SW CA-54 package |
| `dod-wcf-intermediate-ca-1-.pem` | DoD WCF Intermediate CA 1 for allowing TLS connectivity to AWS services in the BCAP |
| `dod-wcf-root-ca-1-.pem` | DoD WCF Root CA 1 for allowing TLS connectivity to AWS services in the BCAP |

## DoD certificate authority package

DISA publishes a package of the public certificates of all DoD certificate
authorities. We use this package to validate all DoD-signed certificates in
dev, test, and production, without having to add them to the underlying OS cert
store.

### Updating the DoD certificate authority package

If one of our partners wants to use a certificate signed by a CA newer than the
package, then get the newest version of the package from DISA. As of this
writing, this package can be found at
[PKI/PKE Document Library](https://public.cyber.mil/pki-pke/pkipke-document-library/).
Search for the term `pkcs` and look for the "PKI CA Certificate Bundles: PKCS#7 For DoD PKI Only - Version X.X".
You can download the certificate from the provided link to get a zip file.

For reference only, the current [pkcs7 v5.6 archive](https://dl.dod.cyber.mil/wp-content/uploads/pki-pke/zip/certificates_pkcs7_v5-6_dod.zip)
can be downloaded and unzipped with:

```sh
curl https://dl.dod.cyber.mil/wp-content/uploads/pki-pke/zip/certificates_pkcs7_v5-6_dod.zip -o certificates_pkcs7_v5-6_dod.zip
unzip certificates_pkcs7_v5-6_dod.zip -d tmp/certs/
cp tmp/certs/Certificates_PKCS7_v5.6_DoD.der.p7b config/tls/
chmod +x config/tls/Certificates_PKCS7_v5.6_DoD.der.p7b
```

Then, update the path accordingly in the `DOD_CA_PACKAGE` variable in `.envrc` and in other places in the code base.

### Updating the DoD SW CA 54 cert

**NOTE:** The certificates for move.mil are signed by DoD SW CA 54. In order to update this cert we also need the
site's certs to be resigned by the new cert. This means unless the cert is expiring it should not be changed.

To download the DoD SW CA cert you will need to visit the [DoD PKI Management Portal](https://crl.gds.disa.mil/).
Select the CA you wish to download and "Submit Selection". Then under the "Certificate Authority Certificate" section
use the "Download" link.

```sh
openssl x509 -inform der -in ~/Downloads/DODSWCA_54.cer -subject -issuer > config/tls/dod-sw-ca-54.pem
```

Then, update the `MOVE_MIL_DOD_CA_CERT` variable in the `.envrc` file and in other places in the code base.

If you need to check the valid dates on the certificate you can use this command:

```sh
$ openssl x509 -inform der -in DODSWCA_54.cer -noout -dates
notBefore=Nov 22 13:51:28 2016 GMT
notAfter=Nov 23 13:51:28 2022 GMT
```

You can read this as the certificate is valid until November 23, 2022.

## DoD WCF certificate authority certificates

These certificates are required when running the app through the BCAP. When running traffic through
the BCAP, various AWS services present TLS certificate chains that include these certificates,
rather than the standard AWS certificate chains. Therefore, standard calls to AWS services will
fail due to certificate validation issues without these certificates.

When updating these, make sure to include .crt at the end of the filename when placing them into
the destination container at /usr/local/share/ca-certificates/ or else update-ca-certificates won't
add them to /etc/ssl/certs/ca-certificates.crt. This is important because Go looks at that file
when searching for trusted CA certificates (<https://golang.org/src/crypto/x509/root_linux.go>) in
both Debian-based (e.g. Distroless) and Alpine-based containers.

## `devlocal` Certificate Authority

The `devlocal` CA is trusted by the system in development and test environments.

### Creating new certificates and signing them using the `devlocal` CA

To get a new certificate signed by the `devlocal` CA, the easiest way is to use
the `scripts/generate-devlocal-cert` script. Here's an example of running that script
to get a key pair named `partner.cer` and `partner.key` with the subject `/C=US/ST=DC/L=Washington/O=Partner/OU=Application/CN=partner.mil`.

```text
$ scripts/generate-devlocal-cert -o Partner -u Application -n partner.mil -f partner
Generating a 2048 bit RSA private key
.............................+++
........................+++
writing new private key to '/Users/me/go/src/github.com/transcom/mymove/partner.key'
-----
Signature ok
subject=/C=US/ST=DC/L=Washington/O=Partner/OU=Application/CN=partner.mil
Getting CA Private Key
SHA256 digest: 21d45ee839ef3416b361d25acc3aa6437cde87e04bfd98619cdc3ec8d47faee7
```

An example for generating the faux orders certs is this:

```sh
scripts/generate-devlocal-cert -o "Not Coast Guard" -u "Not Coast Guard Orders" -n localhost -f devlocal-faux-coast-guard-orders
Generating a RSA private key
..+++++
.....+++++
writing new private key to '/Users/cgilmer/Projects/transcom/mymove/devlocal-faux-coast-guard-orders.key'
-----
Signature ok
subject=/C=US/ST=DC/L=Washington/O=Not Coast Guard/OU=Not Coast Guard Orders/CN=localhost
Getting CA Private Key
SHA256 digest: (stdin)= c5f3d9127e756209c6090b6ade8044e138cac82d2606cf85a7e9e381c4b7b2ac
```

## Adding certificates to the database

Client Certificates are known to the system by their `SHA-256` digest hashes.
`scripts/generate-devlocal-cert` provides that (and the certificate's subject)
in its output, but if you need to do it to an existing certificate, run:

```text
$ mutual-tls-extract-fingerprint config/tls/devlocal-faux-coast-guard-orders.cer
10ac6d7bdb7003093ad82880a2c7ea496e8dc3d50217da6170de83c0be826507
```

For human readability, the table also stores the certificate's subject. To get that, run:

```text
$ mutual-tls-extract-subject config/tls/devlocal-faux-coast-guard-orders.cer
CN=localhost,OU=Not Coast Guard Orders,O=Not Coast Guard,L=Washington,ST=DC,C=US
```

The certificates live in the `client_certs` table for both the MilMove and Orders applications.

## Certs for api.demo.dp3.us and api.loadtest.dp3.us

This certificate and chain were made with SSLMate. The key/pair and chain are held in the Engineer Vault for 1Password.
