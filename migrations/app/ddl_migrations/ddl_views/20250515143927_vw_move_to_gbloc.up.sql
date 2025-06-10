-- B-23540  Daniel Jordan  initial view addition, updating view to consider USMC GBLOC

DROP VIEW move_to_gbloc;
CREATE OR REPLACE VIEW move_to_gbloc
AS SELECT move_id,
    gbloc
   FROM ( SELECT DISTINCT ON (sh.move_id) sh.move_id,
            s.affiliation,
            COALESCE(pctg_m.gbloc, COALESCE(pctg.gbloc, COALESCE(pctg_oconus_bos.gbloc, COALESCE(pctg_oconus.gbloc, pctg_ppm.gbloc)))) AS gbloc
           FROM mto_shipments sh
             JOIN moves m_1 ON sh.move_id = m_1.id
             JOIN orders o ON m_1.orders_id = o.id
             JOIN service_members s ON o.service_member_id = s.id
             LEFT JOIN (SELECT a.id AS address_id,
                    'USMC'::character varying AS gbloc,
                    pctg_1.postal_code
                   FROM addresses a
                     JOIN postal_code_to_gblocs pctg_1 ON a.postal_code::text = pctg_1.postal_code::text) pctg_m ON pctg_m.address_id = sh.pickup_address_id
                     	AND s.affiliation = 'MARINES'
             LEFT JOIN ( SELECT a.id AS address_id,
                    pctg_1.gbloc AS gbloc,
                    pctg_1.postal_code
                   FROM addresses a
                     JOIN postal_code_to_gblocs pctg_1 ON a.postal_code::text = pctg_1.postal_code::text) pctg ON pctg.address_id = sh.pickup_address_id
             LEFT JOIN ( SELECT ppm.shipment_id,
                    pctg_1.gbloc
                   FROM ppm_shipments ppm
                     JOIN addresses ppm_address ON ppm.pickup_postal_address_id = ppm_address.id
                     JOIN postal_code_to_gblocs pctg_1 ON ppm_address.postal_code::text = pctg_1.postal_code::text) pctg_ppm ON pctg_ppm.shipment_id = sh.id
             LEFT JOIN ( SELECT a.id AS address_id,
                    jr.code::character varying AS gbloc,
                    ga.department_indicator
                   FROM addresses a
                     JOIN re_oconus_rate_areas ora ON a.us_post_region_cities_id = ora.us_post_region_cities_id
                     JOIN gbloc_aors ga ON ora.id = ga.oconus_rate_area_id
                     JOIN jppso_regions jr ON ga.jppso_regions_id = jr.id) pctg_oconus_bos ON pctg_oconus_bos.address_id = sh.pickup_address_id AND
                CASE
                    WHEN s.affiliation = 'AIR_FORCE'::text THEN 'AIR_AND_SPACE_FORCE'::text
                    WHEN s.affiliation = 'SPACE_FORCE'::text THEN 'AIR_AND_SPACE_FORCE'::text
                    ELSE s.affiliation
                END = pctg_oconus_bos.department_indicator::text
             LEFT JOIN ( SELECT a.id AS address_id,
                    pctg_1.code AS gbloc,
                    ga.department_indicator
                   FROM addresses a
                     JOIN re_oconus_rate_areas ora ON a.us_post_region_cities_id = ora.us_post_region_cities_id
                     JOIN gbloc_aors ga ON ora.id = ga.oconus_rate_area_id
                     JOIN jppso_regions pctg_1 ON ga.jppso_regions_id = pctg_1.id) pctg_oconus ON pctg_oconus.address_id = sh.pickup_address_id AND pctg_oconus.department_indicator IS NULL
          WHERE sh.deleted_at IS NULL
          ORDER BY sh.move_id, sh.created_at) m;