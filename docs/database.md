# Database Guides

* [How To Backup and Restore the Development Database](how-to/backup-and-restore-dev-database.md#how-to-backup-and-restore-the-development-database)
* [How To Migrate the Database](how-to/migrate-the-database.md#how-to-migrate-the-database)

## Pop SQL logging on by default in development

Pop is an ORM which helps ease communication with the database by providing database API abstraction code. However, this can obscure the actual SQL that is being executed without a in-depth knowledge of the ORM. By enabling SQL logging in development a developer can see the queries being executed by Pop as they happen to hopefully help developers to catch issues in the setup of database calls with Pop.

If you want to turn this off _temporarily_, just prefix your command with `DB_DEBUG=0` for example:

> DB_DEBUG=0 make server_run

If you need to turn this off _permanently_ on your local instance add the following to the `.envrc.local` file

> export DB_DEBUG=0
