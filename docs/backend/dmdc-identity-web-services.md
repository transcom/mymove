# DMDC Identity Web Services

DMDC offers a suite of services under the "Identity Web Services" (IWS) program. These services provide clients limited access to DEERS (Defense Enrollment Eligibility Reporting System) for looking up personnel information on anyone affiliated with the DoD, including service members, civilian employees, contractors, and dependents.

MilMove needs to look up personnel information for several reasons:

1. When given the SSN of a member, get the member's EDIPI, to properly associate electronic Orders and avoid saving the SSN in clear-text.
1. Save time during account creation by pre-populating fields like name, branch of service, and rank.
1. Using a member's EDIPI, get the member's SSN, to facilitate PPM payments or support other legacy paperwork that has yet to move to EDIPI for unique identifiers.
1. Look up a member's work e-mail address, to notify a member who has yet to create an account that they have electronic PCS Orders in MilMove and that they can start the move process.

MilMove uses IWS's "Real-time Broker Service" (RBS). This is an allegedly RESTful webservice. It utilizes XML payloads.

## Protocol

DMDC calls RBS RESTful, but it's really not. The client sends a GET request with the parameters encoded in the path of the URL (which is RESTful) but still using the `=` sign (which is not RESTful).

The HTTP status code it returns is meaningless, even in error cases.

Instead, the response body contains one of two completely unrelated types of XML documents - one for normal operation, and another for protocol errors (like malformed data). After determining which of the two document types the body contains, the client must look for the reason and return codes in the document.

The "normal" document type is specified in `2675_DDS_get_cac_data.xsd`. Find it in the Google drive.

## Authentication

We have provided DMDC with our DoD-signed certificates. Our staging and experimental certs are allowed to connect to their CT environment, while our production cert is allowed to connect to their production environment.

As of this writing, IWS presents an ordinary commercial signed cert, although they once presented a DoD-signed cert. Perhaps that means they would trust clients with commercial certs now instead of requiring DoD-signed certs. **As a practical matter, when connecting to RBS, we trust certs signed by both the DoD and the usual commercial CAs.** Who knows, DMDC might switch back.

## Enclaves

IWS has two enclaves - production, where the real data lives, and "Contractor Test", where fake data lives.

It should go without saying that any data retrieved from IWS's Production Environment is PII and needs to be appropriately protected.

As for Contractor Test, it is _shared with other customers,_ some of whom _**have write access**_. These other customers (IIRC this includes Tricare and the VA) won't delete any records, but they frequently change the populations associated with a given test record. (They want to simulate a person separating from the military, for example.)

Therefore, do not write any tests against CT that expect a test record to never change in any details. You can rely on the SSN-to-EDIPI to stay constant.

On the other hand, this resembles production, as records of real people obviously change over time.

Unfortunately, CT is infrequently "refreshed," where the previous data set is obliterated and a new data set is created with all-new fake people. When this happens, reach out to our DMDC Project Officer and request the updated data. In the past, they supplied the data as an Excel spreadsheet in an attachment to an encrypted email.

## Network

We have also provided DMDC with MilMove's outbound static IPs. At the time of this writing, they are not limiting us to our staging and experimental IPs _in CT_, so connecting from developer machines is allowed assuming you have the staging / experimental private DoD keys handy.

## Request Types

We can look up personnel information using one of the following - SSN + name, EDIPI, work e-mail, or CAC token.

### PIDS : SSN + name

Provide SSN, name, and the DOB (if available, which it usually isn't), and get 0, 1, or more person records back, along with a MatchReasonCode that explains why IWS returned what it did. In the interest of brevity, only the reason codes returned for SSN queries (as opposed to sponsor queries) are listed below. They are:

Code | Enum | Description
-|-|-
PAB | MatchReasonCodeMultiple | More than one record matched the provided SSN and last name ("more than one PN_ID matched the provided criteria")
PMB | MatchReasonCodeLimited | SSN matched a record, but the other criteria (like last name) didn't match ("the person matched on PN_ID and PN_ID_TYP_CD only")
PMC | MatchReasonCodeFull | SSN + last name (and first name, if provided) matched just one record ("the person matched on PN_ID, PN_ID_TYP_CD and at least one additional criterion")
PNB | MatchReasonCodeNone | Not found ("no person matched the provided PN_ID and PN_ID_TYP_CD combination")

### EDI : DoD ID Number / EDIPI

Provide EDIPI, and get a person record back, or an error.

### wkEma: Work E-mail

Provide work e-mail, and get a person record back, or an error. This query is implemented on our end but not currently used.

### TIDS: CAC Token

Provide certificate information from a CAC and get a person record back, or an error.

This query is unimplemented on our end and untested. This query could potentially be very useful if we start taking CAC logins directly from service members and want to accelerate profile creation and fetch electronic Orders. The URL to use will look like `https://pkict.dmdc.osd.mil/appj/rbs/rest/op=tids03/customer=2675/schemaName=get_cac_data/schemaVersion=1.0/PKIC_AUTH_NM=CN=DOD%20CA-30,%20OU=PKI,%20OU=DoD,%20O=U.S.%20Government,%20C=US/PKIC_SER_ID=000000`

## Populations

MilMove has access to many, but not all, of the "populations" of people in DEERS. It is possible for individuals to be part of multiple populations at the same time. It is also possible to leave or join a population due to events like Separation.

| Category Code | Population | Accessible |
| ------------- | ---------- | ---------- |
| A | Active Duty Member | **Yes** |
| B | Presidential Appointee | **Yes** |
| C | DoD / Uniformed Service Civil Service Employee | **Yes** |
| D | Disabled American Veteran | **Yes** |
| F | Former member | **Yes** |
| J | Service Academy Student | **Yes** |
| K | NAF DoD / Uniformed Service employee | **Yes** |
| N | National Guard Member | **Yes** |
| Q | Gray Area Retiree | **Yes** |
| R | Retiree | **Yes** |
| V | Reserve Member | **Yes** |
| Y | Civilian Retiree | **Yes** |
| E | DoD and Uniformed Service contract employee | _No_ |
| H | Medal of Honor recipient | _No_ |
| I | Non-DoD civil service employee, except Presidential appointee | _No_ |
| L | Lighthouse Service (obsolete) | _No_ |
| M | Non-federal Agency civilian associates | _No_ |
| O | Non-DoD contract employee | _No_ |
| T | Foreign Affiliate | _No_ |
| U | DoD OCONUS Hire | _No_ |
| W | DoD Beneficiary | _No_ |

## Pitfalls

* RBS is usually fast (&lt;1 seconds), but sometimes takes a long time (&gt;25 seconds) to respond to requests. Keep this in mind when setting timeouts.
* If you try to look up an individual but do not have access to any of that individual's populations, RBS will simply return no match, as if that individual does not exist.
* SSNs do not have a 1-to-1 relationship to individuals. Individuals may be issued more than one SSN in their lifetimes for completely legitimate reasons. One SSN can also be associated with multiple people, due to typos or fraud. To account for this, **RBS returns a MatchReasonCode along with possibly multiple matches.** Clients are supposed to send as much information (names and/or DOB) with the PIDS request to increase the odds of returning a single match.
* RBS will still return a match if it has a person with the provided SSN even if the provided name does not match; it gives the MatchReasonCode "PMB" in this case.
