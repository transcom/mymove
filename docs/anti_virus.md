# Anti-Virus

This document provides an overview of the anti-virus solution employed by the MilMove project.

Anti-virus scanning is concerned with two parts of the system. The first part is viruses being injected into the source code of the project by folks with permissions to commit and merge code. The second part is viruses being uploaded to the AWS S3 bucket by users and downloaded by other users of the system.

## Source Code Scanning

Anti-virus scanning of the source code is done via a step in the CI/CD pipeline which in turn calls the `make anti_virus` target in the `Makefile`. The step uses a ClamAV Docker container to scan a copy of the checked out source code prior to building any binaries or docker images for deployment. If a virus is detected the build cannot continue and the anti-virus finding must be dealt with.

Files that are false positives (like PDF test fixtures) can be white-listed in the `scripts/anti-virus` file by updating the variable `IGNORE_FILES`.

## Upload Object Scanning in AWS S3

Anti-virus scanning of uploads is done via an AWS Lambda that responds to AWS S3 Creation events to respond in real time to uploads and scan them immediately. The application and AWS users are forbidden from downloading unscanned files or files that have been marked as infected. Infected files also cannot be re-tagged as clean, preventing circumvention of the AV solution.

### Object Tagging

The solution here relies on a terraform module named [trussworks/terraform-aws-s3-anti-virus](https://github.com/trussworks/terraform-aws-s3-anti-virus) which deploys lambda code named [upsidetravel/bucket-antivirus-function](https://github.com/upsidetravel/bucket-antivirus-function).  The lambda adds tags `av-status`, `av-timestamp` and `av-signature` to the object with the following key/values:

| Key | Values | Notes |
| --- | --- | --- |
| av-status | CLEAN/INFECTED | No other values are allowed |
| av-timestamp | RFC 3339 timestamp | This is a string |
| av-signature | ClamAV Signature String | This is the finding result for the file |
| av-notes | GENERATED/MANUAL | Notes not added by the lambda but via the application or other scripts |

The application or manual scripts can also be invoked to scan files or mark them as known to be CLEAN. These tools are requested to add a tag `av-notes` for auditing purposes. However, no other values for `av-status` other than `CLEAN` or `INFECTED` are allowed.

### Scanned buckets

The scanned buckets include:

- transcom-ppp-app-devlocal-us-west-2
- transcom-ppp-app-experimental-us-west-2
- transcom-ppp-app-staging-us-west-2
- transcom-ppp-app-prod-us-west-2

All key-prefixes are scanned in these buckets including `app/` and `secure-migrations/`. In all cases the threat is that a user uploads a file that is downloaded by another user and infects their machine. There is a threat model around an ECS container handling viruses (like migrations) but this is not the primary reason all these objects are scanned.
