-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

-- * Leaving empty- move created inmodifying production moves show status not needed on development *

--- uncomment the following lines for testing
-- UPDATE moves SET show = false where id in (
-- 	'4f3f4bee-3719-4c17-8cf4-7e445a38d90e' -- "Advance, PPM" - locator: NOADVC (new)
-- 	, '27266e89-df79-4469-8843-05b45741a818' -- "ApproveShipment, HHGPPM" - locator: COMBO2 (ppms/accepted)
-- 	, 'fb4105cf-f5a5-43be-845e-d59fdb34f31c' -- "ReadyToInvoice, HHG" locator: DOOB (delivered)
-- 	, '9992270d-4a6f-44ea-82f6-ae3cf3277c5d' -- "ReadyForApprove, HHG" locator: NOCHKA (completed)
-- 	, 'b2ecbbe5-36ad-49fc-86c8-66e55e0697a7' -- "UserPerson2, HHGDude2" locator: ZPGVED (all)
-- );