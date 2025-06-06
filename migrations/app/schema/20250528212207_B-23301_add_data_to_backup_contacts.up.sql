-- B-23301  Jim Hawks  convert existing data to first and last name columns

DO $$
DECLARE
    rec RECORD;
    parts TEXT[];
    idx INT;
BEGIN
    FOR rec IN SELECT id, name FROM backup_contacts LOOP
        -- Trim and split the name string into an array of words
        parts := regexp_split_to_array(trim(rec.name), '\s+');

        -- Remove title if it exists
        IF array_length(parts, 1) >= 1 AND lower(parts[1]) IN ('mr', 'mrs', 'miss', 'dr') THEN
            parts := parts[2:array_length(parts, 1)];
        END IF;

        -- Get new length of array
        idx := array_length(parts, 1);

        -- Update based on how many words are left
        IF idx >= 1 THEN
            UPDATE backup_contacts
            SET first_name = parts[1],
                last_name = CASE
                              WHEN idx >= 2 THEN parts[idx]
                              ELSE ''
                            END
            WHERE id = rec.id;
        ELSE
            -- Handle case where name is empty or just a title
            UPDATE backup_contacts
            SET first_name = '',
                last_name = ''
            WHERE id = rec.id;
        END IF;
    END LOOP;
END $$;

