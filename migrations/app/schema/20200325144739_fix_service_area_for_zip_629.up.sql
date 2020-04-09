-- Fix service area for zip 629 (Vienna, IL); previously was 252.
UPDATE tariff400ng_zip3s
SET service_area = '256'
WHERE zip3 = '629';
