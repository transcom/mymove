-- There are no electronic orders in production yet, but the uniqueness rules
-- have changed since testing began. This truncation only affects local dev
-- and experimental, and allows the appropriate indexes to be created.
DELETE FROM electronic_orders;
