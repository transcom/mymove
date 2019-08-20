# How To Troubleshoot GEX Connection

## 1. Retrieve certificates from chamber

[Chamber](https://github.com/segmentio/chamber) is a tool for managing secrets. We use it to access secret values stored in SSM Parameter Store in AWS.

Run the following commands from within the project directory:

```console
$ chamber read app-experimental move_mil_dod_tls_key > tmp/secret.key
$ chamber read app-experimental move_mil_dod_tls_cert > tmp/secret.cert
```

You will now have the client certificate and key saved into `tmp`.**Delete these files when you are done using them.**

## 2. Verify the Connection

```console
$ curl --cert tmp/secret.cert --key tmp/secret.key --insecure -v https://gexweba.daas.dla.mil/msg_data/submit/Test
```

You should see the following if the connection was **successful**:

```console
$ curl --cert tmp/secret.cert --key tmp/secret.key --insecure -v https://gexweba.daas.dla.mil/msg_data/submit/Test
*   Trying 131.78.200.191...
* TCP_NODELAY set
* Connected to gexweba.daas.dla.mil (131.78.200.191) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* Cipher selection: ALL:!EXPORT:!EXPORT40:!EXPORT56:!aNULL:!LOW:!RC4:@STRENGTH
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/cert.pem
  CApath: none
* TLSv1.2 (OUT), TLS handshake, Client hello (1):
* TLSv1.2 (IN), TLS handshake, Server hello (2):
* TLSv1.2 (IN), TLS handshake, Certificate (11):
* TLSv1.2 (IN), TLS handshake, Server key exchange (12):
* TLSv1.2 (IN), TLS handshake, Request CERT (13):
* TLSv1.2 (IN), TLS handshake, Server finished (14):
* TLSv1.2 (OUT), TLS handshake, Certificate (11):
* TLSv1.2 (OUT), TLS handshake, Client key exchange (16):
* TLSv1.2 (OUT), TLS handshake, CERT verify (15):
* TLSv1.2 (OUT), TLS change cipher, Client hello (1):
* TLSv1.2 (OUT), TLS handshake, Finished (20):
* TLSv1.2 (IN), TLS change cipher, Client hello (1):
* TLSv1.2 (IN), TLS handshake, Finished (20):
* SSL connection using TLSv1.2 / ECDHE-RSA-AES256-GCM-SHA384
* ALPN, server did not agree to a protocol
* Server certificate:
*  subject: C=US; O=U.S. GOVERNMENT; OU=DOD; OU=PKI; OU=DLA; CN=gexweba.daas.dla.mil
*  start date: Feb 20 16:13:52 2018 GMT
*  expire date: Feb 20 16:13:52 2021 GMT
*  issuer: C=US; O=U.S. Government; OU=DoD; OU=PKI; CN=DOD ID SW CA-38
*  SSL certificate verify ok.
> GET /msg_data/submit/Test HTTP/1.1
> Host: gexweba.daas.dla.mil
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 401 Unauthorized
< Date: Wed, 10 Oct 2018 16:33:47 GMT
< WWW-Authenticate: Basic realm="msg_data"
< Content-Length: 452
< Content-Type: text/html; charset=iso-8859-1
< Set-Cookie: TS0=01a9f; Path=/; Domain=.gexweba.daas.dla.mil; Secure; HTTPOnly
<
<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<html><head>
<title>401 Unauthorized</title>
</head><body>
<h1>Unauthorized</h1>
<p>This server could not verify that you
are authorized to access the document
requested.  Either you supplied the wrong
credentials (e.g., bad password), or your
browser doesn't understand how to supply
the credentials required.</p>
<hr>
<address>Apache Server at gexweba.daas.dla.mil Port 443</address>
</body></html>
* Connection #0 to host gexweba.daas.dla.mil left intact
```

If you see other output, then there is most likely an issue with the server verifying the client certificates.
