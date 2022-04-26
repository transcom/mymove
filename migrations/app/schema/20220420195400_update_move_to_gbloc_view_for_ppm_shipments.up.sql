-- This view still finds the GBLOC for the first shipment of each move. PPM shipments
-- don't have full address information, though, so this looks to the postal code field
-- on PPM shipments to match it to a GBLOC. NTS-Release shipments won't have a match,
-- but they are handled in a different spot, so the null value is ok here.
CREATE OR REPLACE VIEW move_to_gbloc AS
SELECT DISTINCT ON (sh.move_id) sh.move_id AS move_id, COALESCE(pctg.gbloc, pctg_ppm.gbloc) AS gbloc
FROM mto_shipments sh
     -- try the pickup_address path
     LEFT JOIN
     (
        SELECT a.id address_id, pctg.gbloc
        FROM addresses a
        JOIN postal_code_to_gblocs pctg ON a.postal_code = pctg.postal_code
     ) pctg ON pctg.address_id = sh.pickup_address_id
     -- try the ppm_shipments path
     LEFT JOIN
     (
        SELECT ppm.shipment_id, pctg.gbloc
        FROM ppm_shipments ppm
        JOIN postal_code_to_gblocs pctg ON ppm.pickup_postal_code = pctg.postal_code
     ) pctg_ppm ON pctg_ppm.shipment_id = sh.id
WHERE sh.deleted_at IS NULL
ORDER BY sh.move_id, sh.created_at;
