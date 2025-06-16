-- B-23565 Ricky Mettler initial view creation

CREATE OR REPLACE VIEW v_intl_locations AS
select icc.id icc_id,
	   c.city_name,
	   cpd.country_prn_dv_id,
	   cpd.country_prn_dv_nm,
	   cpd.country_prn_dv_cd,
	   r.country,
	   c.id intl_cities_id,
	   cpd.id re_country_prn_division_id,
	   c.country_id
  from intl_city_countries icc
  join re_intl_cities c on icc.intl_cities_id = c.id
  join re_country_prn_divisions cpd on icc.country_prn_division_id = cpd.id
  join re_countries r on icc.country_id = r.id;