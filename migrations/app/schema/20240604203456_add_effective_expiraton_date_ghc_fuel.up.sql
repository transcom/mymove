ALTER Table ghc_diesel_fuel_prices
ADD column IF NOT EXISTS effective_date date,
ADD column IF NOT EXISTS end_date date;

-- update current records with effective date and end date
-- business rule is that the diesel fuel prices are posted on Mondays and are effective Tuesday and end the following Monday
update ghc_diesel_fuel_prices set effective_date = publication_date + interval '1' day;
update ghc_diesel_fuel_prices set end_date = effective_date + interval '6' day;