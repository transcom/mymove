# Migration Files

Migration files are split by application:

- app: The app service which manages my, office, admin, etc
- orders: The orders service which manages only electronic orders (tbd)

Inside each application folder are two folders:

- schema: The schema migrations that do not contain data that is considered "secure"
- secure: A local copy of "secure" migrations that are a placeholder for the same migrations hosted in AWS S3

Each app also has a file named `migrations_manifest.txt`. This file lists the complete set of migrations required for
the application to run.
