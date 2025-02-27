-- TODO: Add this to the DML manifest
INSERT INTO re_shipment_type_prices (
    id,
    contract_id,
    service_id,
    market,
    factor,
    created_at,
    updated_at
)
VALUES
(
    'FDBF472F-4C6C-4D38-8482-A86D3F797473',
    -- TODO: Update this to Beth's branch where contract ID is set to 070f7c82-fad0-4ae8-9a83-5de87a56472e
    -- we use a static contract ID as no data present allows for lookups
    '51393fa4-b31c-40fe-bedf-b692703c46eb',
    (
      SELECT id
      FROM re_services
      WHERE code = 'INPK'
      LIMIT 1
    ),
    'O',
    -- fetched from https://dp3.atlassian.net/wiki/spaces/MT/pages/2720890895/International+Pricing#INPK%3A
    -- TODO: When Beth's upstream is in int replace this with a secure migration and real data
    -- TODO: This is a fake value
    1.15,
    NOW(),
    NOW()
),
(
    '59CC8348-1A22-4DC5-9F60-29C0A599060C',
    -- TODO: Update this to Beth's branch where contract ID is set to 070f7c82-fad0-4ae8-9a83-5de87a56472e
    -- we use a static contract ID as no data present allows for lookups
    '51393fa4-b31c-40fe-bedf-b692703c46eb',
    (
      SELECT id
      FROM re_services
      WHERE code = 'INPK'
      LIMIT 1
    ),
    'C',
    -- fetched from https://dp3.atlassian.net/wiki/spaces/MT/pages/2720890895/International+Pricing#INPK%3A
    -- TODO: When Beth's upstream is in int replace this with a secure migration and real data
    -- TODO: This is a fake value
    0.30,
    NOW(),
    NOW()
)
-- In case it's pre-existing
ON CONFLICT (contract_id, service_id, market)
DO
UPDATE
  SET factor = EXCLUDED.factor,
      updated_at = EXCLUDED.updated_at;
