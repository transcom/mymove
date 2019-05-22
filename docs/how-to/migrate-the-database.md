# How To Migrate the Database

If you need to change the database schema, you'll need to write a migration.

<!-- markdownlint-disable MD029 MD038 -->

## Running Migrations

To run a migration you should use the `milmove migrate` command. This is done in a similar way to the `soda migrate`
command but uses our DB connection code instead. This avoids us having to use the `database.yaml` as a config file
and allows us to leverage different authentication methods for migrations in development and in production using
the same code.  To migrate you should use a command based on your DB:

* `make db_dev_migrate`
* `make db_test_migrate`
* `make db_prod_migrations_migrate`

The reason to use a `make` target is because it will put you into the `scripts/` directory from which it is required
you run the migration so that `scripts/apply-secure-migrations.sh` is called with the correct paths for the different
files in the `./migrations` folder.

## Creating a migration

Use soda (a part of [pop](https://github.com/gobuffalo/pop/)) to generate models and migrations.

> **We don't use down-migrations to revert changes to the schema; any problems are to be fixed by a follow-up migration.**

### Generating a New Model

If you are generating a new model, use: `gen-model model-name column-name:type column-name:type ...`. The fields `id`, `created_at`, and `updated_at` are all created automatically.

### Modifying an Existing Model

If you are modifying an existing model, use `soda generate migration migration_name` and add the [Fizz commands](https://github.com/gobuffalo/fizz) yourself to the created `{migration_name}.up.fizz` file. Delete the `down.fizz` file, as we aren't using those (see note below.)

## Zero-Downtime Migrations

As a good practice, all of our migrations should create a database state that works both with the current version of the application code _and_ the new version of the application code. This allows us to run migrations before the new app code is live without creating downtime for our users. More in-depth list of migrations that might cause issues are outlined in our [google drive](https://docs.google.com/document/d/1ht57qz1ut--fqTQdLKbCqbZO_f_S0UoVSIyO6Bg-wJw).

Eg: If we need to rename a column, doing a traditional rename would cause the app to fail if the database changes went live before the new application code (pointing to the new column name) went live. Instead, this should be done in a two-stage process.

1. Write a migration adding a new column with the preferred name and copy the data from the old column into it. The old column will effectively be deprecated at this point.
2. After the migration and new app code have been deployed to production, write a second migration to remove the old/deprecated column.

Similarly, if a column needs to be dropped, we should deprecate the column in one pull request and then actually remove it in a follow-up pull request. Deprecation can be done by renaming the column to `deprecated_column_name`. This process has an added side affect of helping us keep our migrations reversible, since columns can always be re-added, but getting old data back into those columns is a more difficult process.

## Secure Migrations

> **Before adding SSNs or other PII, please consult with Infra.**

We are piggy-backing on the migration system for importing static datasets. This approach causes problems if the data isn't public, as all of the migrations are in this open source repository. To address this, we have what are called "secure migrations."

### Creating Secure Migrations

1. Generate new migration files: `generate-secure-migration <migration_name>`. This creates two migration files:
    * a local test file with no secret data
    * a production file to be uploaded to S3 and contain sensitive data
2. Edit the production migration first, and put whatever sensitive data in it that you need to.
3. Copy the production migration into the local test migration.
4. Scrub the test migration of sensitive data, but use it to test the gist of the production migration operation.
5. Test the local migration by running `make db_dev_migrate`. You should see it run your local migration.
6. Test the secure migration by running `make run_prod_migrations` to setup a local  prod_migrations` database. Then run `psql-prod-migrations< tmp/$NAME_OF_YOUR_SECURE_MIGRATION`. Verify that the updated values are in the database.
7. If you are wanting to run a secure migration for a specific non-production environment, then **skip to the next section**.
8. Upload the migration to S3 with: `upload-secure-migration <production_migration_file>`. **NOTE:** For a single environment see the next section.
9. Run `make run_prod_migrations` to verify that the upload worked and that the migration can be applied successfully. If not, you can make changes and run `upload-secure-migration` again and it will overwrite the old files.
10. Once the migration is working properly, **delete the secure migration from your `tmp` directory** if you didn't delete it in step 8.
11. Open a pull request!
12. When the pull request lands, the production migrations will be run on Staging and Prod.

### Secure Migrations for One Environment

To run a secure migration on ONLY staging (or other chosen environment), upload the migration only to the S3 environment and blank files to the others:

1. Similar to the "Upload the migration" step above, run `ENVIRONMENTS="staging" upload-secure-migration <production_migration_file>` where `ENVIRONMENTS` is a quoted list of all the environments you wish to upload to.  The default is `"experimental staging prod"` but you can just do staging and production with `ENVIRONMENTS="staging prod"`
2. Check that it is listed in the S3 staging secure-migrations folder: `aws s3 ls s3://transcom-ppp-app-staging-us-west-2/secure-migrations/`
3. Check that it is NOT listed in the S3 production folder: `aws s3 ls s3://transcom-ppp-app-prod-us-west-2/secure-migrations/`
4. Now upload empty files of the same name to the prod and experimental environments: `ENVIRONMENTS="experimental prod" upload-secure-migration <empty_migration_file_with_same_name>`
5. To verify upload and that the migration can be applied use the make target corresponding to your environment:

* `make run_prod_migrations`
* `make run_staging_migrations`
* `make run_experimental_migrations`

### How Secure Migrations Work

When secure migrations are run, `soda` will shell out to our script, `apply-secure-migration.sh`. This script will:

* Look at `$SECURE_MIGRATION_SOURCE` to determine if the migrations should be found locally (`local`, for dev & testing,) or on S3 (`s3`).
* If the file is to be found on S3, it is downloaded from `${AWS_S3_BUCKET_NAME}/secure-migrations/${FILENAME}`.
* If it is to be found locally, the script looks for it in `$SECURE_MIGRATION_DIR`.
* Regardless of where the migration comes from, it is then applied to the database by essentially doing: `psql-prod-migrations < ${FILENAME}`.

There is an example of a secure migration [in the repo](https://github.com/transcom/mymove/blob/master/migrations/20180424010930_test_secure_migrations.up.fizz).

### Downloading Secure Migrations

**NOTE:** Be careful with downloading secure migrations. They often contain sensitive input and should be treated with care. When
you are done with these secure migrations please delete them from your computer.

You may need to download and inspect secure migrations. Or perhaps you need to correct a file you uploaded by mistake. Here is how you download the secure migrations:

* Download the migration to S3 with: `download-secure-migration <production_migration_file>`. You can also use the `ENVIRONMENTS` environment variable to specify one or more than one environment.
* This will put files in `./tmp/secure_migrations/${environment}`.

You can now inspect or modify and re-upload those files as needed.
