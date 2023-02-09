# Identity Web Services

We are using the "Real-Time Broker Service" (RBS) part of the Identity Web Services (IWS) which is a service
of the Defense Manpower Data Center (DMDC).

The utility in this folder is helpful in understanding the code used in translating between SSN and EDIPI for
service members.

## Example usage

The way to use this code is:

```sh
$ bin/iws -dod_ca_package config/tls/Certificates_PKCS7_v5.9_DoD.der.p7b -ssn 666585049 -last "Escudero OCampo"
Identity Web Services: Real-Time Broker Service (REST)
Host: pkict.dmdc.osd.mil
Operation: pids-P
SSN: 666585049
Last Name: Escudero OCampo
First Name:
No match
```
