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

To get a new certificate signed by the `devlocal` CA, the easiest technique is
to have `openssl` generate the key, cert, and certificate signing request in a
single pass. First, make a text file with the certificate details, like this `example.txt`:

```text
[req]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
C=US
ST=IL
L=Belleville
O=Not USTRANSCOM
OU=Not The Real Thing
CN=localhost

[req_ext]
```

Then, from within this directory, run

```text
$ openssl req -nodes -new -config example.txt -keyout example.key -out example.csr
```

That gives you the private key and the CSR. Finally, get the signed cert:

```text
$ openssl x509 -req -in example.csr -CA devlocal-ca.pem -CAkey devlocal-ca.key -CAcreateserial -out example.cer -days 3652 -sha256
```

The resulting `example.cer` certificate is now signed by the `devlocal` CA, valid for 10 years.

## Adding certificates to the database

Client Certificates are known to the system by their `SHA-256` digest hashes. To
get the hash for a cert, run the following (or tell the partner to run it):

```text
$ openssl x509 -outform der -in example.cer | openssl dgst -sha256
```

For human readability, the row in the database also stores the certificate's
subject. To get that, run:

```text
$ openssl x509 -noout -subject -in example.cer | sed -e 's/subject= \///g;s/\//, /g'
```

(That `sed` script removes the prefix "subject= /" from `openssl`'s output and converts the subject's "/" characters to more readable commas.)