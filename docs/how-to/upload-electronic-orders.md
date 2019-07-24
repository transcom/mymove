# How to Upload Electronic Orders Using your CAC

## Requirements

The requirements are that you have a CAC and a CAC Reader. The recommended reader is the [Type C Smart Card Reader
Saicoo DOD Military USB-C Common Access Card Reader](https://www.amazon.com/Reader-Saicoo-Military-Compatible-Windows/dp/B071NT53M7/ref=sr_1_4).

To get going you will also need to install some software on your machine:

1. `brew install opensc` gives you tools like `pkcs11-tool` and `pkcs15-tool`
1. Install the official  Military CAC package from [CAC Key Packages](http://militarycac.org/MacVideos.htm#CACKey_packages) to get the driver `/usr/local/lib/pkcs11/cackey.dylib`

It's important that you disable default smart card access on your OSX machine. [Read the published instructions](http://militarycac.org/macuninstall.htm#Mojave_(10.14),_High_Sierra_(10.13.x),_and_Sierra_(10.12.x)_Built_in_Smart_Card_Ability). It boils down to removing your CAC from the reader, running this command and restarting your laptop:

```sh
sudo defaults write /Library/Preferences/com.apple.security.smartcard DisabledTokens -array com.apple.CryptoTokenKit.pivtoken
```

This should get you in the state where you can use your CAC to generate certs.

## Generating a Secure Migration

Orders does mutual TLS authentication and then authorizes you by comparing a SHA 256 hash of your public certificate
stored on your CAC. To get that information into the system you have to create a secure migration using your
SHA 256 fingerprint and the Subject on your CAC certificate.

To get the Fingerprint and Subject you can run this command:

```sh
cac-generate-fingerprint
```

It will print out both on separate lines. Copy these to a clipboard for use later.

Now you need to generate the secure migration:

```sh
export FINGERPRINT="STRING_OF_FINGERPRINT_FROM_CLIPBOARD"
export SUBJECT="STRING_OF_SUBJECT_FROM_CLIPBOARD"
milmove gen orders-migration --fingerprint "${FINGERPRINT}" --subject "${SUBJECT}" --migration-filename "${USERNAME}_cac"
```

The generation script will provide three files and update the `migrations_manifest.txt` for you:

* A secure migration in the `migrations/` folder
* A stub local migration in the `local_migrations/` folder
* A migration to upload to AWS S3 in the `tmp/` folder
