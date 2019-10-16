# Database Guides

* [How To Backup and Restore the Development Database](how-to/backup-and-restore-dev-database.md#how-to-backup-and-restore-the-development-database)
* [How To Migrate the Database](how-to/migrate-the-database.md#how-to-migrate-the-database)

## Pop SQL logging on by default in development

Pop is an ORM which helps ease communication with the database by providing database API abstraction code. However, this can obscure the actual SQL that is being executed without a in-depth knowledge of the ORM. By enabling SQL logging in development a developer can see the queries being executed by Pop as they happen to hopefully help developers to catch issues in the setup of database calls with Pop.

If you want to turn this off _temporarily_, just prefix your command with `DB_DEBUG=0` for example:

> DB_DEBUG=0 make server_run

If you need to turn this off _permanently_ on your local instance add the following to the `.envrc.local` file

> export DB_DEBUG=0

### Some problems to look out for with SQL logging on

#### Excessive Queries (e.g. n+1 Problem)

When looking up objects that have a one-to-many relationship, ORMs such as Pop can fire off n+1 queries to the database to do the look up for n number of child objects + 1 for the original parent object. Depending on the size of n this will cause performance issues loading such lists of objects that have many children. To help prevent this in Pop, developers should avoid using unbounded `Eager()` calls and try and only fetch children when that data is needed. For more through description of the issue please read the following references.

* [What is the "N+1 selects problem" in Object-Relational Mapping?](https://stackoverflow.com/questions/97197/what-is-the-n1-selects-problem-in-orm-object-relational-mapping)
* [N+1 Queries and How to Avoid Them!](https://medium.com/@bretdoucette/n-1-queries-and-how-to-avoid-them-a12f02345be5) -- This uses examples from Ruby on Rails but the concept is the same

#### Excessive Joins (e.g. open-ended *Eager* call)

* [The Dangerous Subtleties of LEFT JOIN and COUNT() in SQL](https://www.xaprb.com/blog/2009/04/08/the-dangerous-subtleties-of-left-join-and-count-in-sql/)
* [More Dangerous Subtleties of JOINs in SQL](https://alexpetralia.com/posts/2017/7/19/more-dangerous-subtleties-of-joins-in-sql)

## Using RDS IAM for database authentication

RDS IAM authentication is the method of connecting to the database using IAM as the authentication mechanism as opposed to a conventional username and password. More information can be found [here](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html).

### Locally connecting to RDS using IAM authentication

An RDS instance must be configured with IAM authentication before connecting. All RDS in the MilMove environment has this enabled. If wishing to connect to a RDS instance first prepare the infrastructure to allow this with this [guide](https://github.com/transcom/ppp-infra/blob/master/docs/runbook/0009-accessing-the-db-with-iam.md). Once complete locally milmove server can be ran with the following

```bash
/path/to/milmove serve --db-iam --db-iam-role arn:aws:iam::AWSACCOUNT:role/CONNECTROLE  --db-region us-east-2 --db-host RDSURL  --db-ssl-mode verify-full --db-ssl-root-cert bin/rds-combined-ca-bundle.pem  --db-user db_user
```

### ECS Task connecting to RDS using IAM authentication

ECS tasks such use RDS IAM authentication to securely connect without a username or passwords to rotate. This is accomplished by ECS assigning a role to the container that is allowed to connect to a specific database via IAM.

The MilMove server through the use of environment variables will use reach out to IAM to generate a temporary connection token, almost similar to a password. This token/password is valid for only 15 minutes. To enable IAM authentication ensure these environment variables are present for `app`, `app-client-tls`, and `migration` containers. Here is a snippet of the required environment [variables](https://github.com/transcom/mymove/blob/6426a37eaf0219323aef997deed5a43e0e1a824b/config/app.container-definition.json#L32-L39) for the [app.container-definition.json](https://github.com/transcom/mymove/blob/master/config/app.container-definition.json) that is deployed.

```json
{
  "name": "DB_IAM",
  "value": "{{ .DB_IAM }}"
},
{
  "name": "DB_IAM_ROLE",
  "value": "{{ .DB_IAM_ROLE }}"
},
{
  "name": "DB_REGION",
  "value": "us-west-2"
},
{
  "name": "DB_USER",
  "value": "{{ .DB_USER }}"
},
```

Update the related environment configuration to match. Note that the database user is normally different than `master` as additional configuration is needed to allow a database user to login via IAM. MilMove convention for IAM enabled user is `ecs_user`. Below is a snippet of the [experimental environment config](https://github.com/transcom/mymove/blob/master/config/env/experimental.env):

```ini
DB_USER=ecs_user
DB_IAM=true
DB_IAM_ROLE=YOUR_CONTAINER_ROLE_ARN_HERE
```

### Reverting Task to use password authentication

In the event of a IAM failure it may be desired to revert back to conventional username and password authentication.

1. Get password from Infra from the admin vault in DP3 1Password.

1. Update the Parameter store with the new password

1. ```bash
   chamber write app-YOURENV db_password NEW_PASSWORD
   ```

1. Update the environment configuration files to disable IAM authentication. Keep in mind the database `user` will need to be set to `master`.

1. ```ini
   DB_USER=master
   DB_IAM=false
   ```
