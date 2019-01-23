# How To Run Against S3 Locally

The app defaults to using the local filesystem to store uploaded files in development mode. This is to simplify development setup and
make it easy to inspect uploaded files while developing features. In some situations, though, it is desirable to instead have
uploaded files stored on the real S3 service.

## Prerequisites

You need to have followed [the instructions to setup your AWS developer credentials](https://github.com/transcom/ppp-infra/tree/master/transcom-ppp) in order for the following commands to work.

## Upload Files to S3

The environment variable `STORAGE_BACKEND` specifies if files should be stored on `s3`, on the `local` filesystem, or in `memory`. The default value is `local` in development and `s3` when the app is running in any deployed environment.

### Using the Devlocal Bucket

Assuming your AWS credentials are setup properly, this command will configure the app to upload to the `transcom-ppp-app-devlocal-us-west-2` S3 bucket:

```console
$ env STORAGE_BACKEND=s3 make server_run_standalone
```

_Please note that this does not use our usual setup to automatically reload changes to files in the `swagger` directory. Other code changes should still be detected, however._

### Listing Files on S3

If you want to see what files have been uploaded to an S3 bucket through the console, use `aws s3 ls`:

```bash
# list all files uploaded to the devlocal bucket by you
$ aws s3 ls --recursive s3://transcom-ppp-app-devlocal-us-west-2/$AWS_S3_KEY_NAMESPACE

# You can set $AWS_S3_KEY_NAMESPACE in .envrc.local, but it will default
# to your local username
$ aws s3 ls --recursive s3://transcom-ppp-app-devlocal-us-west-2/bilbo
```

### Cleaning up S3

When you are done testing, remove any uploaded objects from the S3 bucket using the following command:

```bash
$ aws s3 rm --recursive s3://transcom-ppp-app-devlocal-us-west-2/$AWS_S3_KEY_NAMESPACE
```

This will delete all files from the `dev-local` bucket that were uploaded by you. Use the above listing command to verify that this worked.
