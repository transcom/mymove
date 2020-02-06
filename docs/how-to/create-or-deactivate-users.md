# How To Create or Deactivate Users

## Creating Users

We now create users for the Office and Admin applications in the Admin console and not via secure migrations so that we can capture the security object audit logs.

### Creating DPS Users

For creating users let's assume that the new user's email is username@example.com.

Here is an example migration to create a DPS user (please edit as needed):

```sql
INSERT INTO public.dps_users
    (id, login_gov_email,
     created_at, updated_at, active)
    VALUES (
        GENERATED_UUID4_VAL, 'username@example.com',
        now(), now(), true
    );
```

## Deactivating Users

MilMove doesn't delete users because of both auditing concerns and CASCADE DELETE failures. Instead each
user table has a `active` boolean column that can be used to deactivate a user. Deactivating a user means the
person with valid credentials to Login.gov may not be permitted to get a session in MilMove.

There are several places you can deactivate a user, at the global level and at each application level. It's important
to deactivate users at the application level if you are concerned that a user entry was made for the user in that
table but they have not yet claimed the user entry by logging in. This is an issue to be aware of for both Office
and Admin users.

### Deactivating Users Globally

An example of deactivating a user by email:

```sql
UPDATE users SET active = false WHERE email = 'username@example.com';
```

This is the only way to deactivate Service Members.

### Deactivating Office Users

This should now be done in the admin user application:

1. Navigate to the admin app and log in
2. Click on "Office Users" in the sidebar
3. Select the user you would like to deactivate and click "Edit" in the top right
   corner. Here, you'll be able to choose whether or not that user is deactivated.

### Deactivating DPS Users

An example of deactivating a DPS user by email:

```sql
UPDATE dps_users SET active = false WHERE email = 'username@example.com';
```

### Generating a migration to deactivate a specific user

We may need to deactivate specific milmove (not office or admin) users. Until this functionality is available in the admin console, you can use the following `milmove` sub-command:

`milmove gen disable-user-migration -e EMAIL`
