-- any office users created prior to implementation of transportation_office_assignments table
-- may not have a row in the transportation_office_assignments table, so lets populate it if that is the case.
-- please note: the id col on transportation office assignments table is NOT a primary key column, its a foreign column to the office_users table;
-- updates to transportation_office_assignments table will be handled in E-06880

INSERT INTO transportation_office_assignments (id, transportation_office_id, created_at, updated_at, primary_office)
SELECT ou.id, ou.transportation_office_id, now(), now(), true
FROM office_users ou
WHERE NOT EXISTS (
    SELECT 1
    FROM transportation_office_assignments toa
    WHERE toa.id = ou.id
)
AND ou.transportation_office_id IS NOT NULL;
