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
[https://iase.disa.mil/pki-pke/pages/tools.aspx](https://iase.disa.mil/pki-pke/pages/tools.aspx). Click the "Trust Store" tab,
and download the "For DoD PKI Only" ZIP file under the "`PKI CA Certificate
Bundles: PKCS#7`" heading. If their naming convention has not changed, you want the `Certificates_PKCS7_vX.X_DoD.der.p7b` file in that archive, where `X.X` is
the current version.

Then, update the path accordingly in the `DOD_CA_PACKAGE` variable in `.envrc`.

## `devlocal` Certificate Authority

The `devlocal` CA is trusted by the system in development and test environments.

### Creating new certificates and signing them using the `devlocal` CA

To get a new certificate signed by the `devlocal` CA, the easiest way is to use
the `bin/generate-devlocal-cert.sh` script. Here's an example of running that script
to get a key pair named `partner.cer` and `partner.key` with the subject `/C=US/ST=DC/L=Washington/O=Partner/OU=Application/CN=partner.mil`.

```text
$ bin/generate-devlocal-cert.sh -o Partner -u Application -n partner.mil -f partner
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
`bin/generate-devlocal-cert.sh` provides that (and the certificate's subject)
in its output, but if you need to do it to an existing certificate, run:

```text
$ openssl x509 -outform der -in example.cer | openssl dgst -sha256
```

For human readability, the table also stores the certificate's subject. To get that, run:

```text
$ openssl x509 -noout -subject -in example.cer
```

The certificates live in the `client_certs` table.
