# Schema Documentation

This documentation was created using [SQLEditor](https://www.malcolmhardie.com/sqleditor/)

It was initially seeded from the (then) current `/migrations/schema.sql` and should
 be maintained as models are added to the DB.

## Installing the tool

You can either download the installer directly from the [website](https://www.malcolmhardie.com/sqleditor/)
 or install SQLEditor via brew

```console
my_machine ~:$ brew cask install sqleditor
```

If you use SQLEditor frequently, you will need to purchase a license for the tool.

## Updating and changing the model

When you make changes to the model objects/add migrations to the pop code you should also make the corresponding changes
 to the diagram in the same PR.

Reviewers should be able to match up the changes in the migrations/schema.sql to diffs in the dp3.sqs file.

NB. _We may decide to use the diffs from this tool as a way to generate migrations in the future_
