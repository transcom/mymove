-- B-23301  Jim Hawks  convert existing data to first and last name columns

DO $$
DECLARE
    rec RECORD;
    parts TEXT[];
    idx INT;
    suffixes TEXT[] := ARRAY['jr', 'jr.', 'sr', 'sr.', 'ii', 'iii', 'iv', 'v'];
BEGIN
    FOR rec IN SELECT id, name FROM backup_contacts LOOP
        parts := regexp_split_to_array(trim(rec.name), '\s+');

        -- Remove title at the beginning
        IF array_length(parts, 1) >= 1 AND lower(parts[1]) IN (
            'mr', 'mr.', 'mrs', 'mrs.', 'miss', 'ms', 'ms.', 'dr', 'dr.'
        ) THEN
            parts := parts[2:array_length(parts, 1)];
        END IF;

        idx := array_length(parts, 1);

        IF idx = 1 THEN
            UPDATE backup_contacts
            SET first_name = parts[1],
                last_name = ''
            WHERE id = rec.id;

        ELSIF idx = 2 THEN
            UPDATE backup_contacts
            SET first_name = parts[1],
                last_name = parts[2]
            WHERE id = rec.id;

        ELSIF idx >= 3 THEN
            IF lower(parts[idx]) = ANY(suffixes) THEN
                -- Has suffix
                IF idx >= 4 THEN
                    -- Skip middle name (2nd part), use 1st and 3rd+
                    UPDATE backup_contacts
                    SET first_name = parts[1],
                        last_name = array_to_string(parts[3:idx], ' ')
                    WHERE id = rec.id;
                ELSE
                    -- Only 3 parts: treat 2nd as last name, keep suffix
                    UPDATE backup_contacts
                    SET first_name = parts[1],
                        last_name = parts[2] || ' ' || parts[3]
                    WHERE id = rec.id;
                END IF;
            ELSE
                -- No suffix: discard middle, use parts[3:] as last name
                UPDATE backup_contacts
                SET first_name = parts[1],
                    last_name = array_to_string(parts[3:idx], ' ')
                WHERE id = rec.id;
            END IF;

        ELSE
            UPDATE backup_contacts
            SET first_name = '',
                last_name = ''
            WHERE id = rec.id;
        END IF;
    END LOOP;
END $$;

