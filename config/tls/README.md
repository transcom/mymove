# TLS certificates

## DoD certificate authority package

DISA publishes a package of the public certificates of all DoD certificate
authorities. We use this package to easily validate all DoD-signed certificates
in dev, test, and production, without having to add them to the underlying OS
cert store.

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

## `devlocal` Certificate Authority

The `devlocal` CA is trusted by the system in development and test environments.

### Creating new certificates and signing them using the `devlocal` CA

To get a new certificate signed by the `devlocal` CA, the easiest way is to use
the `scritps/generate-devlocal-cert` script. Here's an example of running that script
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

## Adding certificates to the database

Client Certificates are known to the system by their `SHA-256` digest hashes.
`scripts/generate-devlocal-cert` provides that (and the certificate's subject)
in its output, but if you need to do it to an existing certificate, run:

```text
$ openssl x509 -outform der -in example.cer | openssl dgst -sha256
```

For human readability, the table also stores the certificate's subject. To get that, run:

```text
$ openssl x509 -noout -subject -in example.cer
```

The certificates live in the `client_certs` table.
