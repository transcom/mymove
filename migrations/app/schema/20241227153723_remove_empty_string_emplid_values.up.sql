-- removing empty string values from the emplid column
UPDATE service_members
SET emplid = null
WHERE emplid = '';