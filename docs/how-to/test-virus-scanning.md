# How to Test Virus Scanning

In experimental, staging, and production we run an anti-virus scan using [ClamAV](https://www.clamav.net) to ensure that files uploaded by users
aren't infected with viruses or malware.

If you want to test this functionality on any of the deployed environments, you will need access to a

The most common way to test AV software is using [EICAR test files](https://en.wikipedia.org/wiki/EICAR_test_file). Unfortunately, these files are plain text and
as a result are rejected by Milmove for not being one of the allowed content types (JPG, PNG, or PDF).

Instead, we can use one of the test files that is provided by the ClamAV project itself. These test files are generated
when the project is built and are not installed by the `clamav` formula in homebrew. To access them without needing to build the
entire ClamAV project we can extract them from a Debian package.

Please don't check this file into our repository as it is under the GNU license and this project is not.

```shell script
$ curl http://http.us.debian.org/debian/pool/main/c/clamav/clamav-testfiles_0.101.4+dfsg-1_all.deb -o clamav-testfiles_0.101.4+dfsg-1_all.deb
$ ar x clamav-testfiles_0.101.4+dfsg-1_all.deb
$ tar -xzvf data.tar.xz

# The test files are available in ./usr/share/clamav-testfiles/
```

You can use the PDF filed located within the directory used above at `usr/share/clamav-testfiles/clam.pdf`. This file will
be flagged as containing `Clamav.Test.File-6` by ClamAV.

## Running ClamAV locally

Should you need to debug ClamAV on your development machine, you can install it using the following commands:

```shell script
$ brew install clamav
$ cp /usr/local/etc/clamav/freshclam.conf.sample /usr/local/etc/clamav/freshclam.conf

# Edit /usr/local/etc/clamav/fleshclam.conf to remove the "Example" line

$ freshclam # downloads virus definitions
$ clamscan -va $PATH_TO_FILE_OR_DIR_TO_SCAN
```

It takes `clamscan` a while to run, especially the first time. Subsequent runs on very small inputs took nearly a minute on
a recent MacBookPro:

```shell script
$ clamscan -va infected/*
Scanning infected/clam.pdf
infected/clam.pdf: Clamav.Test.File-6 FOUND
infected/clam.pdf!(1): Clamav.Test.File-6 FOUND

----------- SCAN SUMMARY -----------
Known viruses: 6508139
Engine version: 0.102.0
Scanned directories: 0
Scanned files: 1
Infected files: 1
Data scanned: 0.03 MB
Data read: 0.03 MB (ratio 1.00:1)
Time: 49.708 sec (0 m 49 s)
```
