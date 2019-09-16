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

When looking up objects that have a one-to-many relationship, ORMs such as Pop can fire off n+1 queries to the database to do the look up for n number of child objects + 1 for the original parent object. Depending on the size of n this will cause performance issues loading such lists of objects that have many children. This is a problem that eager loading seeks to solve by reducing the number of queries by looking child relationships up immediately. For more through description of the issue please read the following references.

* [What is the "N+1 selects problem" in Object-Relational Mapping?](https://stackoverflow.com/questions/97197/what-is-the-n1-selects-problem-in-orm-object-relational-mapping)
* [N+1 Queries and How to Avoid Them!](https://medium.com/@bretdoucette/n-1-queries-and-how-to-avoid-them-a12f02345be5) -- This uses examples from Ruby on Rails but the concept is the same

#### Excessive Joins (e.g. open-ended *Eager* call)

* [The Dangerous Subtleties of LEFT JOIN and COUNT() in SQL](https://www.xaprb.com/blog/2009/04/08/the-dangerous-subtleties-of-left-join-and-count-in-sql/)
* [More Dangerous Subtleties of JOINs in SQL](https://alexpetralia.com/posts/2017/7/19/more-dangerous-subtleties-of-joins-in-sql)
