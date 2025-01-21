-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.

--add is_less_50_miles bool
alter table re_intl_other_prices 
add column if not exists is_less_50_miles  bool;

comment on column re_intl_other_prices.is_less_50_miles is 'Denotes if price for SIT pickup or delivery is for shipments where distance from pickup or delivery and SIT is less than 50 miles';

--alter unique constraint
alter table re_intl_other_prices drop constraint re_intl_other_prices_unique_key;

alter table re_intl_other_prices add constraint re_intl_other_prices_unique_key 
unique (contract_id,service_id,is_peak_period,rate_area_id,is_less_50_miles);


--update re_intl_other_prices.id to prod id and price
DO
$$
declare
	i record;
begin
	for i in 
		(select c.code, b.id prod_id, b.per_unit_cents prod_price, a.id from re_intl_other_prices a, re_intl_other_prices_prod b, re_services c
			where a.service_id = b.service_id
			  and a.is_peak_period = b.is_peak_period
			  and a.rate_area_id = b.rate_area_id
			  and a.contract_id = b.contract_id
			  and a.service_id = c.id
			and a.contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e')
	loop
		begin
			update re_intl_other_prices
			   set id = i.prod_id,
				   per_unit_cents = i.prod_price
			 where id = i.id;
	    exception when others then null;
		end;
		
	end loop;
	
end $$;

--set is_less_50_miles
--current IOPSIT - is_less_50_miles = true
--current IDDSIT - is_less_50_miles = false

update re_intl_other_prices
   set is_less_50_miles = true
 where service_id = '6f4f6e31-0675-4051-b659-89832259f390'
   and is_less_50_miles is null;

update re_intl_other_prices
   set is_less_50_miles = false
 where service_id = '28389ee1-56cf-400c-aa52-1501ecdd7c69'
   and is_less_50_miles is null;

--add missing recs - re_intl_other_prices.sql