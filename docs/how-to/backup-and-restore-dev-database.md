# How To Backup and Restore the Development Database

## Backup

Run `bin/db-backup` to backup the `dev_db` to a file on the filesystem. The backup files stored in `tmp/db/`.

```console
$ bin/db-backup clean-state
```

## Restore a Backup

Run `bin/db-restore` to overwrite the `dev_db` database with the contents of the named backup:

```console
$ bin/db-restore clean-state
```

**This is a destructive command!** All data currently in `dev_db` will be removed when this command is run.

## List Existing Backups

When called without a backup name `bin/db-restore` will list available backups:

```console
$ bin/db-restore
Available backups are:

         clean-slate   6.1M   Oct 30 16:31:12
                boom   5.9M   Oct 19 10:45:03
```
