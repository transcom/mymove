-- B-23301  Jim Hawks  convert existing data to first and last name columns

DO $$
DECLARE
    rec RECORD;
    parts TEXT[];
    idx INT;
BEGIN
    FOR rec IN SELECT id, name FROM backup_contacts LOOP
        -- Trim and split the name into parts (handles multiple spaces)
        parts := regexp_split_to_array(trim(rec.name), '\s+');

        -- Remove title if present at the beginning
        IF array_length(parts, 1) >= 1 AND lower(parts[1]) IN (
            'mr', 'mr.', 'mrs', 'mrs.', 'miss', 'ms', 'ms.', 'dr', 'dr.'
        ) THEN
            parts := parts[2:array_length(parts, 1)];
        END IF;

        -- Recalculate length after title removal
        idx := array_length(parts, 1);

        -- Update based on the number of parts
        IF idx = 1 THEN
            -- Only one word: use as first name
            UPDATE backup_contacts
            SET first_name = parts[1],
                last_name = ''
            WHERE id = rec.id;

        ELSIF idx = 2 THEN
            -- Two words: first and last
            UPDATE backup_contacts
            SET first_name = parts[1],
                last_name = parts[2]
            WHERE id = rec.id;

        ELSIF idx >= 3 THEN
            -- Three or more words: discard middle (2nd), use 1st and 3rd+
            UPDATE backup_contacts
            SET first_name = parts[1],
                last_name = array_to_string(parts[3:idx], ' ')
            WHERE id = rec.id;

        ELSE
            -- Empty or just a title
            UPDATE backup_contacts
            SET first_name = '',
                last_name = ''
            WHERE id = rec.id;
        END IF;
    END LOOP;
END $$;

