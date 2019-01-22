# IAM Authentication for Database

* Status: accepted
* Deciders: @dynamike, @pjdufour-truss
* Date: 2018-10-15

## Context and Problem Statement

Rotating passwords to the MyMove database is a multi-step process that requires synchronizing changes among application servers and migrations.  AWS now provides an easy solution to outsource authentication for an RDS database to [AWS Identity and Access Management](https://aws.amazon.com/iam/) (IAM).

[https://aws.amazon.com/about-aws/whats-new/2018/09/amazon-rds-postgresql-now-supports-iam-authentication/](https://aws.amazon.com/about-aws/whats-new/2018/09/amazon-rds-postgresql-now-supports-iam-authentication/)

The default PostgreSQL authentication method is to use [password authentication](https://www.postgresql.org/docs/10/static/auth-methods.html).  A hash of the password is stored internally within the database.  There are other external authentication methods that support connecting from the same host machine or require running additional infrastructure, e.g., an LDAP service.  However, none of those methods are worth the operational burden for managing authentication for a single application or help with the credential rotation problem.

IAM is the access manager for many AWS services and has a standardized API for the AWS CLI ([https://aws.amazon.com/cli/](https://aws.amazon.com/cli/)), Terraform ([https://www.terraform.io/](https://www.terraform.io/)), and Go SDK ([https://docs.aws.amazon.com/sdk-for-go/api/service/iam/](https://docs.aws.amazon.com/sdk-for-go/api/service/iam/)).  Should we switch to using IAM?

## Decision Drivers

* Easy to rotate passwords or access keys
* Maximize Dev-prod parity
* Strong passwords / access keys
* Works with [sqlx](https://github.com/jmoiron/sqlx), our database driver for PostgreSQL
* Works with [Pop](https://github.com/gobuffalo/pop), our ORM-like framework for interacting with our PostgreSQL database
* Works within our migrations workflow.

## Considered Options

* Internal PostgreSQL Authentication Provider
* External IAM Authentication Provider

## Decision Outcome

Chosen option: "External IAM Authentication Provider".

## Pros and Cons of the Options

### Internal PostgreSQL Authentication Provider

* `-` Credential rotation is complex.
* `+` Same authentication method for development and production.
* `+` Strong passwords are supported and the responsibility of the infrastructure team.  PostgreSQL supports password lengths up to [at most 100 bytes](https://stackoverflow.com/questions/19499058/pgsql-character-limitations).  Only a hash of the database password is stored.  PostgreSQL stores password hashes using MD5 by default, but also now supports the more secure [scram-sha-256](https://www.postgresql.org/docs/11/static/auth-password.html).
* `+` Our database driver [sqlx](https://github.com/jmoiron/sqlx) supports internal passwords out of the box.
* `+/-` Our database framework Pop ([https://github.com/gobuffalo/pop](https://github.com/gobuffalo/pop)) loads passwords from configuration files ([https://github.com/gobuffalo/pop#example-configuration-file](https://github.com/gobuffalo/pop#example-configuration-file)).  This is not a great pattern and can be refractored using [Viper](https://godoc.org/github.com/spf13/viper) in the future, but it does currently work.
* `+` As part of a migration, we can change passwords using the `ALTER ROLE` SQL command ([https://www.postgresql.org/docs/10/static/sql-alterrole.html](https://www.postgresql.org/docs/10/static/sql-alterrole.html)).

### External IAM Authentication Provider

* `+` Credential rotation is handled by AWS.  Using the [Task Role](https://forums.aws.amazon.com/thread.jspa?threadID=284417) without any use of static passwords, we can automatically retrieve a valid set of credentials, which are then used to generate an auth token to connect to the database, e.g., `authToken, err := rdsutils.BuildAuthToken(host, region, cfg.User, stscreds.NewCredentials(sess, "arn:aws:iam::[AWS ID]:role/SomeRole"))`
* `-` Different authentication method for development and production.
* `+` Following of the [Shared Responsibility Model](https://aws.amazon.com/compliance/shared-responsibility-model/), the AWS IAM service is responsible for generating strong access keys, with the specific hashing algorithm internal to AWS.
* `+` rdsutils (`github.com/aws/aws-sdk-go/service/rds/rdsutils`) provides the ability to generate ephemeral authentication tokens to connect to AWS RDS using IAM ([https://docs.aws.amazon.com/sdk-for-go/api/service/rds/rdsutils/](https://docs.aws.amazon.com/sdk-for-go/api/service/rds/rdsutils/))
* `+/-` Our database framework Pop ([https://github.com/gobuffalo/pop](https://github.com/gobuffalo/pop)) may support IAM authentication by providing the connection string including the ephemeral access token as a URL string in [ConnectionDetails](https://godoc.org/github.com/gobuffalo/pop#ConnectionDetails) to [NewConnection](https://godoc.org/github.com/gobuffalo/pop#NewConnection).  Hard to know if this will work out of the box without a proof of concept.
* `+` As part of a migration, we can enable/disable RDS authentication for a given web application role using `GRANT rds_iam TO user_xyz;`.
