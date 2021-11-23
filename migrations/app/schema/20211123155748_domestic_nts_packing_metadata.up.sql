-- Update the code and name for domestic NTS packing (don't need the "factor" part)
UPDATE re_services
SET code = 'DNPK',
	name = 'Domestic NTS packing'
WHERE code = 'DNPKF';

-- For consistency, go ahead and do the same for international NTS packing even though
-- we're not doing OCONUS moves yet
UPDATE re_services
SET code = 'INPK',
	name = 'International NTS packing'
WHERE code = 'INPKF';
