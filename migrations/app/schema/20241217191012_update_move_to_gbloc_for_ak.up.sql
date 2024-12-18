delete from postal_code_to_gblocs where postal_code in (
select uspr_zip_id from v_locations where state = 'AK');

drop view move_to_gbloc;
CREATE OR REPLACE VIEW public.move_to_gbloc
AS
SELECT move_id, gbloc FROM (
SELECT DISTINCT ON (sh.move_id) sh.move_id, s.affiliation,
    COALESCE(pctg.gbloc, coalesce(pctg_oconus_bos.gbloc, coalesce(pctg_oconus.gbloc, pctg_ppm.gbloc))) AS gbloc
   FROM mto_shipments sh
   JOIN moves m ON sh.move_id = m.id
   JOIN orders o on m.orders_id = o.id
   JOIN service_members s on o.service_member_id = s.id
     LEFT JOIN ( SELECT a.id AS address_id,
            pctg_1.gbloc, pctg_1.postal_code
           FROM addresses a
             JOIN postal_code_to_gblocs pctg_1 ON a.postal_code::text = pctg_1.postal_code::text) pctg ON pctg.address_id = sh.pickup_address_id
     LEFT JOIN ( SELECT ppm.shipment_id,
            pctg_1.gbloc
           FROM ppm_shipments ppm
             JOIN addresses ppm_address ON ppm.pickup_postal_address_id = ppm_address.id
             JOIN postal_code_to_gblocs pctg_1 ON ppm_address.postal_code::text = pctg_1.postal_code::text) pctg_ppm ON pctg_ppm.shipment_id = sh.id
     LEFT JOIN ( SELECT a.id AS address_id,
            cast(jr.code as varchar) AS gbloc, ga.department_indicator
           FROM addresses a
             JOIN re_oconus_rate_areas ora ON a.us_post_region_cities_id = ora.us_post_region_cities_id
             JOIN gbloc_aors ga ON ora.id = ga.oconus_rate_area_id
             JOIN jppso_regions jr ON ga.jppso_regions_id = jr.id
           		) pctg_oconus_bos ON pctg_oconus_bos.address_id = sh.pickup_address_id
           				and case when s.affiliation = 'AIR_FORCE' THEN 'AIR_AND_SPACE_FORCE'
           				         when s.affiliation = 'SPACE_FORCE' THEN 'AIR_AND_SPACE_FORCE'
           				         when s.affiliation = 'NAVY' THEN 'NAVY_AND_MARINES'
           				         when s.affiliation = 'MARINES' THEN 'NAVY_AND_MARINES'
           				         else s.affiliation
           				    end = pctg_oconus_bos.department_indicator
     LEFT JOIN ( SELECT a.id AS address_id,
            cast(pctg_1.code as varchar) AS gbloc, ga.department_indicator
           FROM addresses a
             JOIN re_oconus_rate_areas ora ON a.us_post_region_cities_id = ora.us_post_region_cities_id
             JOIN gbloc_aors ga ON ora.id = ga.oconus_rate_area_id
             JOIN jppso_regions pctg_1 ON ga.jppso_regions_id = pctg_1.id
           		) pctg_oconus ON pctg_oconus.address_id = sh.pickup_address_id and pctg_oconus.department_indicator is null
     WHERE sh.deleted_at IS NULL);