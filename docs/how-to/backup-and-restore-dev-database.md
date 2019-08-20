# How To Backup and Restore the Development Database

## Backup

Run `scripts/db-backup` to backup the `dev_db` to a file on the filesystem. The backup files stored in `tmp/db/`.

```console
$ scripts/db-backup clean-state
```

If you'd like to backup a database other than `dev_db`, specify it by setting the value of the `DB_NAME` environment variable. In bash, the following command will backup `test_db`:

```console
$ DB_NAME=test_db scripts/db-backup clean-slate
```

## Restore a Backup

Run `scripts/db-restore` to overwrite the `dev_db` database with the contents of the named backup:

```console
$ scripts/db-restore clean-state
```

**This is a destructive command!** All data currently in `dev_db` will be removed when this command is run.

If you'd like to restore to a database other than `dev_db`, specify it by setting the value of the `DB_NAME` environment variable. In bash, that looks like this:

```console
$ DB_NAME=test_db scripts/db-restore clean-slate
```

## List Existing Backups

When called without a backup name `scripts/db-restore` will list available backups:

```console
$ scripts/db-restore
Available backups are:

         clean-slate   6.1M   Oct 30 16:31:12
                boom   5.9M   Oct 19 10:45:03
```
