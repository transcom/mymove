ALTER TABLE addresses ADD COLUMN IF NOT EXISTS us_post_region_cities_id uuid;

update addresses set postal_code = substr(postal_code,1,5) where length(postal_code) > 5;
update addresses set postal_code = '61866' where city = 'RANTOUL' and postal_code = '61868';
update addresses set city = 'Corona' where city = 'Coronal' and state = 'CA';
update addresses set city = 'Washington' where city = 'Washington DC' and state = 'DC';
update addresses set city = 'ANN ARBOR' where city = 'ANN HARBOR' and state = 'MI';
update addresses set city = 'Thornton' where city = 'Thorton' and state = 'CO';
update addresses set city = 'Ridgecrest' where upper(city) = 'CHINA LAKE' and state = 'CA';
update addresses set city = 'Jber' where city = 'JB Elmendorf-Richardson' and state = 'AK';
update addresses set city = 'TUSKEGEE INSTITUTE' where city = 'TUSKEGGE INSTITUTE' and state = 'AL';
update addresses set city = 'DAVIS MONTHAN AFB' where city = 'Davis-Monthan AFB' and state = 'AZ';
update addresses set city = 'MARCH AIR RESERVE BASE' where city = 'MARCH AIR FORCE BASE' and state = 'CA';
update addresses set city = 'LEMOORE NAVAL AIR STATION' where city = 'NAS Lemoore' and state = 'CA';
update addresses set city = 'POINT MUGU NAWC' where city = 'POINT MUGU' and state = 'CA';
update addresses set city = 'VANDENBERG AFB' where city = 'Vandenberg SFB' and state = 'CA';
update addresses set city = 'BUCKLEY AIR FORCE BASE' where city = 'Buckley SFB' and state = 'CO';
update addresses set city = 'CHEYENNE MOUNTAIN AFB' where city = 'CHEYENNE MTN AFB' and state = 'CO';
update addresses set city = 'PETERSON AFB' where city = 'Peterson SFB' and state = 'CO';
update addresses set city = 'USAF ACADEMY' where city in ('USAFA','U.S. Air Force Academy') and state = 'CO';
update addresses set city = 'WASHINGTON NAVY YARD' where city = 'WA NAVY YARD' and state = 'DC';
update addresses set city = 'HOMESTEAD AIR FORCE BASE' where city = 'HOMSTEAD AFB' and state = 'FL';
update addresses set city = 'LAKE WORTH' where city = 'Lake Worth Beach' and state = 'FL';
update addresses set city = 'PATRICK AIR FORCE BASE' where city = 'Patrick SFB' and state = 'FL';
update addresses set city = 'KEKAHA' where state = 'HI' and postal_code = '96752';
update addresses set city = 'FORT JOHNSON' where city = 'Fort Johnson South' and state = 'LA';
update addresses set city = 'SAINT LOUIS' where upper(city) = 'ST LOUIS' and state = 'MO';
update addresses set city = 'KEESLER AFB' where upper(city) = 'KESSLER AFB' and state = 'MS';
update addresses set city = 'POPE ARMY AIRFIELD' where upper(city) = 'POPE FIELD' and state = 'NC';
update addresses set city = 'O''NEILL' where upper(city) = 'ONEILL' and state = 'NE';
update addresses set city = 'MC GUIRE AFB' where upper(city) = 'MCGUIRE AFB' and state = 'NJ';
update addresses set city = 'WRIGHT PATTERSON AFB' where upper(city) = 'WRIGHT-PATTERSON AFB' and state = 'OH';
update addresses set city = 'SHAW AFB' where upper(city) = 'SHAW AIR FORCE BASE' and state = 'SC';
update addresses set city = 'JBSA FT SAM HOUSTON' where city = 'JBSA Fort Sam Houston' and state = 'TX';
update addresses set city = 'JBSA RANDOLPH' where city = 'JBSA Randolph AFB' and state = 'TX';
update addresses set city = 'NAVAL AIR STATION JRB' where city in ('NAS JRB Fort Worth','Naval AS JRB, Naval Air Station JRB') and state = 'TX';
update addresses set city = 'CHARLOTTESVILLE' where city = 'CHARLOTTSVILLE' and state = 'VA';
update addresses set city = 'FORT GREGG ADAMS' where city = 'Fort Gregg-Adams' and state = 'VA';
update addresses set city = 'RICHMOND' where city = 'Richmond,' and state = 'VA';
update addresses set city = 'SUFFOLK' where city = 'SUFFORLK' and state = 'VA';
update addresses set city = 'NORWICH' where city = 'NORFIELD' and state = 'VT';
update addresses set city = 'JOINT BASE LEWIS MCCHORD' where city in ('JB Lewis-McChord','Joint Base Lewis-McChord') and state = 'WA';
update addresses set city = 'ROLLING BAY' where city = 'Rollingbay' and state = 'WA';
update addresses set city = 'FE WARREN AFB' where city = 'F.E. Warren AFB' and state = 'WY';

update addresses a
   set us_post_region_cities_id = u.uprc_id
from (
	select c.city_name uprc_city,
		   s.state uprc_state,
		   upr.uspr_zip_id uprc_zip,
		   uprc.usprc_county_nm uprc_county,
		   uprc.id uprc_id
	from us_post_region_cities uprc
	join re_us_post_regions upr
	  on uprc.us_post_regions_id = upr.id
	join re_cities c
	  on uprc.cities_id = c.id
	join re_states s
	  on upr.state_id = s.id
 ) u
where upper(a.county) = u.uprc_county
and upper(a.city) = u.uprc_city
and a.postal_code = u.uprc_zip
and a.state = u.uprc_state
and a.us_post_region_cities_id is null;