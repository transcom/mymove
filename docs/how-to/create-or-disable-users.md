# How To Create or Disable Users

For all users you will [create a secure migration](./migrate-the-database.md#secure-migrations). Please
only create a migration corresponding to the environment where you intend to add or disable users. For instance,
please only add a `staging` secure migration if you intend to add staging users and leave the `prod`
and `experimental` migrations empty.

## Creating Users

The only users for which Truss is responsible to create are the users of the Office and TSP apps. These
users are usually team members in the staging environment and and actual JPPSO and TSP personnel in the
production environment.

### A note about `uuid_generate_v4()`

Please **do not** use `uuid_generate_v4()` in your SQL. Instead please generate a valid UUID4 value. You can
get a valid UUID4 value from [the Online UUID Generator](https://www.uuidgenerator.net/). You can also use
`python -c 'import uuid; print(str(uuid.uuid4()))'` or `brew install uuidgen; uuidgen`.

In this document anywhere you see `GENERATED_UUID4_VAL` you will need to give a unique UUID4 value (i.e. don't reuse
the same value across different tables.

### Creating Office Users

For creating users let's assume that the new user's email is username@example.com.

For Truss Office users in the staging environment please use these values:

| User Email | `transportation_office_id` |
| --- | --- |
| `username@example.com` | `0931a9dc-c1fd-444a-b138-6e1986b1714c` |

Here is an example migration to create an office user (please edit as needed):

```sql
INSERT INTO public.office_users
    (id, user_id,
     last_name, first_name, middle_initials,
     email, telephone,
     transportation_office_id,
     created_at, updated_at, disabled)
    VALUES (
           GENERATED_UUID4_VAL, NULL,
           'Jones', 'Alice', NULL,
           'username@example.com', '(415) 891-0828',
           '0931a9dc-c1fd-444a-b138-6e1986b1714c',
            now(), now(), false
     );
```

Writing this migration by hand, can get tedious if there are many office users to add, so

### Creating TSP Users

For creating users let's assume that the new user's email is username@example.com.

For Truss TSP users in the staging environment please use these values:

| User Email | `transportation_service_provider_id` |
| --- | --- |
| `username+pyvl@example.com` | `c71bdb14-ed86-4c92-bf06-93c0865f5070` |
| `username+dlxm@example.com` | `b98d3deb-abe9-4609-8d6e-36b2c50873c0` |
| `username+ssow@example.com` | `b6f06674-1b6b-4b93-9ec6-293d5d846876` |

Here is an example migration to create a TSP user (please edit as needed):

```sql
INSERT INTO public.tsp_users
    (id, user_id,
     last_name, first_name, middle_initials,
     email, telephone,
     transportation_service_provider_id,
     created_at, updated_at, disabled)
    VALUES (
        GENERATED_UUID4_VAL, NULL,
        'Jones', 'Alice', NULL,
        'username@example.com', '(415) 891-0828',
        'c71bdb14-ed86-4c92-bf06-93c0865f5070',
        now(), now(), false
    );
```

However, if you are creating Truss TSP users in the staging environment then you'll want this instead:

```sql
INSERT INTO public.tsp_users
    (id, user_id,
     last_name, first_name, middle_initials,
     email, telephone,
     transportation_service_provider_id,
     created_at, updated_at, disabled)
    VALUES (
        GENERATED_UUID4_VAL, NULL,
        'Jones', 'Alice', NULL,
        'username+pyvl@example.com', '(415) 891-0828',
        'c71bdb14-ed86-4c92-bf06-93c0865f5070',
        now(), now(), false
    );
INSERT INTO public.tsp_users
    (id, user_id,
     last_name, first_name, middle_initials,
     email, telephone,
     transportation_service_provider_id,
     created_at, updated_at, disabled)
    VALUES (
        GENERATED_UUID4_VAL, NULL,
        'Jones', 'Alice', NULL,
        'username+dlxm@example.com', '(415) 891-0828',
        'b98d3deb-abe9-4609-8d6e-36b2c50873c0',
        now(), now(), false
    );
INSERT INTO public.tsp_users
    (id, user_id,
     last_name, first_name, middle_initials,
     email, telephone,
     transportation_service_provider_id,
     created_at, updated_at, disabled)
    VALUES (
        GENERATED_UUID4_VAL, NULL,
        'Jones', 'Alice', NULL,
        'username+ssow@example.com', '(415) 891-0828',
        'b6f06674-1b6b-4b93-9ec6-293d5d846876',
        now(), now(), false
    );
```

### Creating DPS Users

For creating users let's assume that the new user's email is username@example.com.

Here is an example migration to create a DPS user (please edit as needed):

```sql
INSERT INTO public.dps_users
    (id, login_gov_email,
     created_at, updated_at, disabled)
    VALUES (
        GENERATED_UUID4_VAL, 'username@example.com',
        now(), now(), false
    );
```

## Disabling Users

MilMove doesn't delete users because of both auditing concerns and CASCADE DELETE failures. Instead each
user table has a `disabled` boolean column that can be used to disable a user. Disabling a user means the
person with valid credentials to Login.gov may not be permitted to get a session in MilMove.

There are several places you can disable a user, at the global level and at each application level. It's important
to disable users at the application level if you are concerned that a user entry was made for the user in that
table but they have not yet claimed the user entry by logging in.  This is an issue to be aware of for both Office
and TSP users.

### Disabling Users Globally

An example of disabling a user by email:

```sql
UPDATE users SET disabled = true WHERE email = 'username@example.com';
```

This is the only way to disable Service Members.

### Disabling Office Users

An example of disabling an Office user by email:

```sql
UPDATE office_users SET disabled = true WHERE email = 'username@example.com';
```

### Disabling TSP Users

An example of disabling a TSP user by email:

```sql
UPDATE tsp_users SET disabled = true WHERE email = 'username@example.com';
```

### Disabling DPS Users

An example of disabling a DPS user by email:

```sql
UPDATE dps_users SET disabled = true WHERE email = 'username@example.com';
```
