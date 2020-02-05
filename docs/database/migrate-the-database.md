# How To Migrate the Database

If you need to change the database schema, you'll need to write a migration.

<!-- markdownlint-disable MD029 MD038 -->

## Running Migrations

To run a migration you should use the `milmove migrate` command. This allows us to leverage different authentication
methods for migrations in development and in production using the same code. To migrate you should use a command
based on your DB:

- `make db_dev_migrate`
- `make db_test_migrate`
- `make db_deployed_migrations_migrate`

The reason to use a `make` target is because it will correctly set the migration flag variables and target the correct
database with environment variables.

## Creating a migration

Use the `milmove gen <subcommand>` commands to generate models and migrations. To see a list of available subcommands,
use `milmove gen`. Those subcommands include:

- `migration`: creates a generic migration for you to populate
- `disable-user-migration`: creates a migration for disabling a user given their e-mail address
- `duty-stations-migration`: creates a migration to update duty stations given a CSV of duty station data
- `office-user-migration`: creates a migration to add office users given a CSV of new office user data
- `certs-migration`: creates a migration to add a certificate for access to electronic orders and the prime api

> **We don't use down-migrations to revert changes to the schema; any problems are to be fixed by a follow-up migration.**

### Generating a New Model

If you are generating a new model first add a new migration file with the suffix `up.sql` or `up.fizz`.

**NOTE:** You must define the `PRIMARY KEY` and any indexes you want by yourself. There is nothing in Pop, Fizz, or SQL
that will automatically make these for you. Failure to do so may break tools or lead to inefficiencies in the API
implementations.

When creating the SQL you may write the migration like this:

```sql
create table new_table
(
    id uuid
        constraint new_table_pkey primary key,
    column1 text not null,
    column2 text,
    created_at timestamp not null,
    updated_at timestamp not null
);
```

**NOTE**: PLEASE USE SQL INSTEAD OF FIZZ

If instead you must use `fizz`, the equivalent code is this (but **please use sql**):

```text
create_table("new_table") {
  t.Column("id", "uuid", {primary: true})
  t.Column("column1", "string", {"null": true})
  t.Column("column2", "string", {})
  t.Timestamps()
}
```

Next run `update-migrations-manifest` to update the `migrations_manifest.txt` file.

Then create a new models file in `pkg/models/` named after your new model like `new_table.go`. The contents will look
like:

```go
package models

import (
  "time"

  "github.com/gobuffalo/pop"
  "github.com/gobuffalo/validate"
  "github.com/gobuffalo/validate/validators"
  "github.com/gofrs/uuid"
)

// NewTable represents a new table
type NewTable struct {
  ID        uuid.UUID `json:"id" db:"id"`
  Column1   string    `json:"column1" db:"column1"`
  Column2   *string   `json:"column2" db:"colunn2"`
  CreatedAt time.Time `json:"created_at" db:"created_at"`
  UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewTables is not required by pop and may be deleted
type NewTables []NewTable

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *NewTable) Validate(tx *pop.Connection) (*validate.Errors, error) {
  return validate.Validate(
    &validators.StringIsPresent{Field: r.Column1, Name: "Column1"},
  ), nil
}
```

Now you will want to run the migration to test it out with `make db_dev_migrate` before making your PR.

### Generating a New Migration

If you are generating a new migration, use: `milmove gen migration -n <migration_name>`, which will create a placeholder migration and add it to the manifest.

### A note about `uuid_generate_v4()`

Please **do not** use `uuid_generate_v4()` in your SQL. Instead please generate a valid UUID4 value. You can
get a valid UUID4 value from [the Online UUID Generator](https://www.uuidgenerator.net/). You can also use
`python -c 'import uuid; print(str(uuid.uuid4()))'` or `brew install uuidgen; uuidgen`.

In this document anywhere you see `GENERATED_UUID4_VAL` you will need to give a unique UUID4 value (i.e. don't reuse
the same value across different tables.

#### Reasons why we avoid `uuid_generate_v4()`

We avoid the use of `uuid_generate_v4()` for scripts that add data to the database (esp generating primary keys) because

- It make running migrations multiple times end up with different results
- It makes it hard to use primary keys generated this way as foreign keys in other migrations.
- It raises the remote possibility that a migration works in one system and fails in another
- With specific UUIDs we were able to track down users in each system. When using `uuid_generate_v4()` we have no way of telling what UUID people were assigned on remote machines so we lose the ability to identify them locally.

For more details see this [slack thread](https://ustcdp3.slack.com/archives/CP6PTUPQF/p1559840327095700)

## Zero-Downtime Migrations

As a good practice, all of our migrations should create a database state that works both with the current version of the application code _and_ the new version of the application code. This allows us to run migrations before the new app code is live without creating downtime for our users. More in-depth list of migrations that might cause issues are outlined in our [google drive](https://docs.google.com/document/d/1q-Ho5NINRPpsHQI-DjmLrDlzHsBh-hUc).

Eg: If we need to rename a column, doing a traditional rename would cause the app to fail if the database changes went live before the new application code (pointing to the new column name) went live. Instead, this should be done in a two-stage process.

1. Write a migration adding a new column with the preferred name and copy the data from the old column into it. The old column will effectively be deprecated at this point.
2. After the migration and new app code have been deployed to production, write a second migration to remove the old/deprecated column.

Similarly, if a column needs to be dropped, we should deprecate the column in one pull request and then actually remove it in a follow-up pull request. Deprecation can be done by renaming the column to `deprecated_column_name`. This process has an added side affect of helping us keep our migrations reversible, since columns can always be re-added, but getting old data back into those columns is a more difficult process.

## Secure Migrations

> **Before adding SSNs or other PII, please consult with Infra.**

We are piggy-backing on the migration system for importing static datasets. This approach causes problems if the data isn't public, as all of the migrations are in this open source repository. To address this, we have what are called "secure migrations."

### Creating Secure Migrations

1. Generate new migration files: `generate-secure-migration <migration_name>`. This creates two migration files:
   - a local test file with no secret data
   - a production file to be uploaded to S3 and contain sensitive data
2. Edit the production migration first, and put whatever sensitive data in it that you need to.
3. Copy the production migration into the local test migration.
4. Scrub the test migration of sensitive data, but use it to test the gist of the production migration operation.
5. Test the local migration by running `make db_dev_migrate`. You should see it run your local migration.
6. Test the secure migration by running `make run_prod_migrations` to setup a local `deployed_migrations` database. Then run `psql-deployed-migrations< tmp/$NAME_OF_YOUR_SECURE_MIGRATION`. Verify that the updated values are in the database.
7. If you are wanting to run a secure migration for a specific non-production environment, then **skip to the next section**.
8. Upload the migration to S3 with: `upload-secure-migration <production_migration_file>`. **NOTE:** For a single environment see the next section.
9. Run `make run_prod_migrations` to verify that the upload worked and that the migration can be applied successfully. If not, you can make changes and run `upload-secure-migration` again and it will overwrite the old files.
10. Once the migration is working properly, **delete the secure migration from your `tmp` directory** if you didn't delete it in step 8.
11. Open a pull request!
12. When the pull request lands, the production migrations will be run on Staging and Prod.

### Secure Migrations for One Environment

To run a secure migration on ONLY staging (or other chosen environment), upload the migration only to the S3 environment and blank files to the others:

1. Similar to the "Upload the migration" step above, run `ENVIRONMENTS="staging" upload-secure-migration <production_migration_file>` where `ENVIRONMENTS` is a quoted list of all the environments you wish to upload to. The default is `"experimental staging prod"` but you can just do staging and production with `ENVIRONMENTS="staging prod"`
2. Check that it is listed in the S3 staging secure-migrations folder: `aws s3 ls s3://transcom-ppp-app-staging-us-west-2/secure-migrations/`
3. Check that it is NOT listed in the S3 production folder: `aws s3 ls s3://transcom-ppp-app-prod-us-west-2/secure-migrations/`
4. Now upload empty files of the same name to the prod and experimental environments: `ENVIRONMENTS="experimental prod" upload-secure-migration <empty_migration_file_with_same_name>`
5. To verify upload and that the migration can be applied use the make target corresponding to your environment:

- `make run_prod_migrations`
- `make run_staging_migrations`
- `make run_experimental_migrations`

### How Secure Migrations Work

When migrations are run the `$MIGRATION_MANIFEST` will be checked against files inside the paths listed in
`$MIGRATION_PATH` (a semicolon separated list of local `file://` or AWS S3 `s3://` paths). The migration code
will then run each migration listed in the manifest in order of the Version (which is typically a time stamp
at the front of a file).

- Look at `$MIGRATION_MANIFEST` to determine list of migrations to run (anything not listed will not be run, anything listed but missing will throw an error)
- Look at `$MIGRATION_PATH` to find files locally or in AWS S3. See the `Makefile` for examples.
- If the file is to be found on S3, it is streamed directly into memory instead of downloading.
- If it is to be found locally, the script looks for it in the listed path.

There is an example of local secure migrations [in the repo](https://github.com/transcom/mymove/blob/master/local_migrations/).

### Downloading Secure Migrations

**NOTE:** Be careful with downloading secure migrations. They often contain sensitive input and should be treated with care. When
you are done with these secure migrations please delete them from your computer.

You may need to download and inspect secure migrations. Or perhaps you need to correct a file you uploaded by mistake. Here is how you download the secure migrations:

- Download the migration to S3 with: `download-secure-migration <production_migration_file>`. You can also use the `ENVIRONMENTS` environment variable to specify one or more than one environment.
- This will put files in `./tmp/secure_migrations/${environment}`.

You can now inspect or modify and re-upload those files as needed.
