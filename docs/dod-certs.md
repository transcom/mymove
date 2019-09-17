# DoD Certificates

MilMove has two kinds of server certificates. The first is a normal, commercially-signed cert stored in AWS Certificate Manager (ACM). We use these certificates for the user-facing properties like <https://my.move.mil/>. The other kind of certificate is signed by DISA, and we use those for communications (inbound and outbound) with other DoD entities.

As of this writing, our DoD certificates expire in **September of 2021.** These certificates will need to be renewed before then.

## Current Certificates

Here are the current DoD-signed certificates and how they are used:

| Subject | X509v3 Alternative Names | Environment | CA | Expires | Inbound | Outbound |
| ------- | ------------------------ | ----------- | -- | ------- | ------- | -------- |
| C=US, O=U.S. Government, OU=DoD, OU=PKI, OU=USTRANSCOM, CN=my.move.mil | `my.move.mil, api.move.mil, dps.move.mil, gex.move.mil, office.move.mil, orders.move.mil, tsp.move.mil` | Production | DOD SW CA-54 | Sep 13 14:14:27 2021 GMT | <https://orders.move.mil> | DMDC Identity Web Services (prod) |
| C=US, O=U.S. Government, OU=DoD, OU=PKI, OU=USTRANSCOM, CN=my.staging.move.mil | `my.staging.move.mil, api.staging.move.mil, dps.staging.move.mil, gex.staging.move.mil, office.staging.move.mil, orders.staging.move.mil, tsp.staging.move.mil` | Staging | DOD SW CA-54 | Sep 13 14:20:04 2021 GMT | <https://orders.staging.move.mil> | DMDC Identity Web Services (Contractor Test) |
| C=US, O=U.S. Government, OU=DoD, OU=PKI, OU=USTRANSCOM, CN=my.experimental.move.mil | `my.experimental.move.mil, api.experimental.move.mil, dps.experimental.move.mil, gex.experimental.move.mil, office.experimental.move.mil, orders.experimental.move.mil, tsp.experimental.move.mil` | Experimental | DOD SW CA-54 | Sep 13 14:22:14 2021 GMT | <https://orders.experimental.move.mil> | DMDC Identity Web Services (Contractor Test) |
| C=US, O=U.S. Government, OU=DoD, OU=PKI, OU=USAF, CN=mymove.sddc.army.mil | `mymove.sddc.army.mil` | Production | DOD SW CA-54 | Sep 11 16:28:48 2021 GMT | <https://mymove.sddc.army.mil> | |
| C=US, O=U.S. Government, OU=DoD, OU=PKI, OU=USAF, CN=mymove-staging.sddc.army.mil | `mymove-staging.sddc.army.mil` | Staging | DOD SW CA-54 | Sep 11 16:30:10 2021 GMT | <https://mymove-staging.sddc.army.mil> | |
| C=US, O=U.S. Government, OU=DoD, OU=PKI, OU=USAF, CN=mymove-experimental.sddc.army.mil | `mymove-experimental.sddc.army.mil` | Experimental | DOD SW CA-54 | Sep 11 16:31:02 2021 GMT | <https://mymove-experimental.sddc.army.mil> | |

## Getting a certificate signed by DISA: create a Certificate Signing Request

Generate a Certificate Signing Request (CSR) for each certificate you wish to register or renew.

### CSR Config File

The easiest way to do this with OpenSSL is to create a configuration file with the certificate details and feed that to the openssl command. For example, here is the config file I made for the production my.move.mil CSR:

```ini
[req]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C=US
ST=IL
L=Belleville
O=USTRANSCOM
OU=USTRANSCOM
emailAddress=dp3-integrations@truss.works
CN = my.move.mil

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = my.move.mil
DNS.2 = api.move.mil
DNS.3 = dps.move.mil
DNS.4 = gex.move.mil
DNS.5 = office.move.mil
DNS.6 = orders.move.mil
DNS.7 = tsp.move.mil
```

Absent other guidance, I used the same C, ST, L, O, and OU values as the cert presented by <https://www.ustranscom.mil/>.

Note the alternate names - this file covers a lot of FQDNs, even ones where we won't present the DoD-signed cert as a server. I did this both to preserve flexibility, and to ensure that when MilMove connects to other DoD entities using this cert, validation has the greatest chance of succeeding. If you don't need alternate names or other extensions, delete the `[ req_ext ]` and `[ alt_names ]` sections, along with the line `req_extensions = req_ext` in the `[req]` section.

### New Certificate

The command to generate the CSR with a new private key (no passphrase) is

`openssl req -nodes -new -config <config file path> -keyout <CommonNameHere>.key -out <CommonNameHere>.csr`

### Renew Certificate

To renew a certificate, you probably want to reuse the existing private key. In that case, the command is

`openssl req -new -config <config file path> -key <CommonNameHere>.key -out <CommonNameHere>.csr`

NOTE: do not trust openssl’s `-x509toreq` option. It strips the alternative names and the e-mail address from the certificate. It’s just a convenience option anyway; use the config file input to be certain you’re getting exactly the certificate you want.

### Checking the CSR

To double-check that the CSR contains the right information, the command is

`openssl req -text -noout -verify -in <csr filename>`

## Submitting the CSR

### Choosing a Certificate Authority

For our purposes, you want the most recent (i.e., highest numbered) DOD SW-CA.

### Get Request Numbers

Using either Google Chrome or Microsoft Edge (NOT Internet Explorer) on a NIPR machine, fill in the webform at the desired CA’s website; for DOD SW CA-54, that’s <https://ee-sw-ca-54.csd.disa.mil/ca/ee/ca>. You will need a NIPR account associated with your CAC. On the "2048-bit SSL Server Enrollment Form," you can paste the contents of the CSR file you generated. You will also need to enter the same certificate details, like the CN and alternate names, into the other fields.

For each CSR you submit, you will get a Request Number. Make a note of that, because the PKI RA will need it to approve your request.

### Getting the PKI RA to approve your CSR(s)

Now that DISA has your CSR, you need to contact the appropriate Registration Authority for your system to approve your request. As of this writing, the appropriate RA for USTRANSCOM is the Air Force PKI RA.

Air Force PKI Help Desk: <https://intelshare.intelink.gov/sites/usaf-pki/> afpki.ra@us.af.mil

### Fill out AF RA 2842-2

Fill out the form AF RA 2842-2 using Adobe Acrobat Reader DC. (Do not use non-Adobe PDF products for this, you will have a bad time.) The form is available from the [AF PKI site on Intelink](https://intelshare.intelink.gov/sites/usaf-pki/).

Some guidance on filling out this form:

1. CERTIFICATE INFORMATION
   1. The CERTIFICATE TYPE is “Application Server”.
   1. If you are doing multiple CSRs, the “CERTIFICATE COMMON NAME (CN) / Fully Qualified Domain Name” is “Multiple Device Request”.
   1. For the REQUEST INFORMATION OR ALT LOGIN TOKEN DETAILS, select the same CA that you submitted the CSR to.
   1. Enter the DEVICE INFORMATION Type / OS / Application, e.g., “Application Server / Alpine Linux 3.7 / MilMove”.
1. CERTIFICATE ACCEPTED BY - should be the same person who sent the CSR to DISA. Because you will be digitally signing the form, you can leave 2.e and 2.f blank.
1. Leave section 3 blank.
1. On the second page, fill in the CN, CA, and Request number for each CSR. You don’t need to worry about the alternate names, just the CNs.

Once you have filled out all of the other fields, digitally sign the PDF in section 2.h with your CAC.

### Submit the form to the RA

Send a digitally signed e-mail (using Outlook on a NIPR machine, or with Mail.app) to afpki.ra@us.af.mil. Attach the completed and digitally signed 2842-2 form. The AF RA expects applicants to be on AFNET, so you will need to justify why you are using the PDF process instead of the automated Windows-based process. Here’s what worked for me:

>I am reaching out on behalf of MilMove to complete the CSR process. I have attached the signed AF RA 2842-2 form detailing the certificates that we need and acknowledging my responsibilities.
>
>MilMove is not able to use the LTMA MMP template for obtaining certificates. MilMove is not internal to AFNET or to USTRANSCOM's network and does not run on Windows.

### Download the signed certificates

Assuming the PKI RA accepts your submission, you will get a notification email including instructions on how to download your signed certificates. In short, in Chrome or Edge go back to the CA’s website on NIPR (e.g., <https://ee-id-sw-ca-54.csd.disa.mil/ca/ee/ca/>), click “Retrieval”, and enter the Request numbers from before.

## Reviewing and Verifying certificates

### Reviewing x.509 certificates

To check a Base64-encoded x.509 certificate:

`openssl x509 -text -noout -in <x509 certificate>.cer`

### Reviewing CA certificate chain in PKCS#7 format

To check a Base64-encoded PKCS#7 format certificate chain:

`openssl pkcs7 -text -print_certs -noout -in <cert chain>.p7b`
