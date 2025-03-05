UPDATE re_intl_transit_times
   SET hhg_transit_time = 10
WHERE origin_rate_area_id IN ('b80a00d4-f829-4051-961a-b8945c62c37d','5a27e806-21d4-4672-aa5e-29518f10c0aa')
   OR destination_rate_area_id IN ('b80a00d4-f829-4051-961a-b8945c62c37d','5a27e806-21d4-4672-aa5e-29518f10c0aa');

update re_intl_transit_times
   SET hhg_transit_time = 20
WHERE origin_rate_area_id IN ('9bb87311-1b29-4f29-8561-8a4c795654d4','635e4b79-342c-4cfc-8069-39c408a2decd')
   OR destination_rate_area_id IN ('9bb87311-1b29-4f29-8561-8a4c795654d4','635e4b79-342c-4cfc-8069-39c408a2decd');