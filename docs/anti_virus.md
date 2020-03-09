# Anti-Virus

This document provides an overview of the anti-virus solutions employed by the MilMove project.

Anti-virus scanning is concerned with two parts of the system. The first part is viruses being injected into the source code of the project by folks with permissions to commit and merge code. The second part is viruses being uploaded to the AWS S3 bucket by users and downloaded by other users of the system.

## Source Code Scanning

Anti-virus scanning of the source code is done via a step in the CI/CD pipeline which in turn calls the `make anti_virus` target in the `Makefile`. The step uses a ClamAV Docker container to scan a copy of the checked out source code prior to building any binaries or docker images for deployment. If a virus is detected the build cannot continue and the anti-virus finding must be dealt with.

### Virus Findings

If a report of a virus is found to have been checked into the github repository by the CircleCI anti-virus job it must be assessed as a real threat. DO NOT ASSUME IT IS TEST FLAKINESS. Check out a fresh copy of the repository, assess with the anti-virus script using `make anti_virus` and determine if the virus is indeed detected.

**NOTE:** If scanning repeatedly find the virus then open an Incident first before doing anything else. An incident for a false positive result is better than no incident for a positive result.

During the incident get help in determining if the finding is a False Positive (read Troubleshooting below) or if it needs further action. Use the incident response procedures to move forward.

### Troubleshooting

#### False Positives

Files that are false positives (like PDF test fixtures) can be white-listed in the `scripts/anti-virus` file by updating the variable `IGNORE_FILES`.

Signatures that are false positives (like `PUA.Pdf.Trojan.EmbeddedJavaScript-1`) can be white listed in the `scripts/anti-virus` file by updating the variable `IGNORE_SIGS`.

#### Outdated ClamAV

If you see this warning you can ignore it:

```text
Fri Mar  6 19:57:07 2020 -> ^Your ClamAV installation is OUTDATED!
Fri Mar  6 19:57:07 2020 -> ^Local version: 0.102.1 Recommended version: 0.102.2
Fri Mar  6 19:57:07 2020 -> DON'T PANIC! Read https://www.clamav.net/documents/upgrading-clamav
```

The ClamAV maintainers have not yet released an official version for Debian Alpine. That doesn't stop the ClamAV process for trying to determine if there is a new version available and suggesting that you update. They helpfully added `DON'T PANIC!` to try and dissuade folks from thinking this is the reason they are having trouble with ClamAV.

See [mko-x/docker-clamav#39](https://github.com/mko-x/docker-clamav/issues/39) for more information.

#### Docker container exits unpredictably

The ClamAV process starts off by updating its various databases and definitions for virus detection. Here is some of that output:

```text
Fri Mar  6 19:57:07 2020 -> daily database available for download (remote version: 25743)
Fri Mar  6 19:57:11 2020 -> Testing database: '/store/tmp/clamav-6e94b769e8d37704da1bf78d7727231c.tmp-daily.cvd' ...
Fri Mar  6 19:57:22 2020 -> Database test passed.
Fri Mar  6 19:57:22 2020 -> daily.cvd updated (version: 25743, sigs: 2209759, f-level: 63, builder: raynman)
Fri Mar  6 19:57:22 2020 -> main database available for download (remote version: 59)
Fri Mar  6 19:57:29 2020 -> Testing database: '/store/tmp/clamav-e05bb84fe2271bf55ee55ee2506ff284.tmp-main.cvd' ...
Fri Mar  6 19:57:36 2020 -> Database test passed.
Fri Mar  6 19:57:36 2020 -> main.cvd updated (version: 59, sigs: 4564902, f-level: 60, builder: sigmgr)
Fri Mar  6 19:57:36 2020 -> bytecode database available for download (remote version: 331)
Fri Mar  6 19:57:36 2020 -> Testing database: '/store/tmp/clamav-8c552c424237331bd40cdba9db6f36ff.tmp-bytecode.cvd' ...
Fri Mar  6 19:57:36 2020 -> Database test passed.
Fri Mar  6 19:57:36 2020 -> bytecode.cvd updated (version: 331, sigs: 94, f-level: 63, builder: anvilleg)
Fri Mar  6 19:57:36 2020 -> safebrowsing database available for download (remote version: 49191)
Fri Mar  6 19:57:39 2020 -> Testing database: '/store/tmp/clamav-21790a79b7578c4fa0d73a103ab50d49.tmp-safebrowsing.cvd' ...
Fri Mar  6 19:57:43 2020 -> Database test passed.
Fri Mar  6 19:57:43 2020 -> safebrowsing.cvd updated (version: 49191, sigs: 2213119, f-level: 63, builder: google)
```

Each update requires a network connection and for the website hosting the files to be available. If either the network is unstable or the website is unavailable then the docker container cannot continue startup and will exit. This is why a retry loop is added to the script, to try to catch intermittent problems. Additionally the script waits for the `/tmp/clamd.sock` unix socket to be available as an indicator that the container was able to start up correctly.

#### My database changes are not being captured

You may need to reload the service with `docker exec -it "${CONTAINER_NAME}" clamdscan --reload`.

#### How do I get inside the container to debug

On your local computer you can just run `make anti_virus` and wait until the script finishes. Once completed you can run the following command to drop into the container to debug:

```sh
docker exec -it anti_virus /bin/bash
```

This should drop you in a bash shell. The project files can be found in the `MOUNT_DIR` from the `scripts/anti-virus` script. The ClamAV database files are found at `/store`.

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
