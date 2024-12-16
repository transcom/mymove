-- Add enum and column to track the anti-virus processing and availability of an upload

CREATE TYPE av_status_type AS ENUM (
    'PROCESSING',
    'CLEAN',
    'INFECTED'
);

ALTER TABLE uploads ADD COLUMN IF NOT EXISTS av_status av_status_type; -- default null, will update to match s3 on first access for column

COMMENT ON TYPE av_status_type IS 'The matching type for the anti-virus status.';
COMMENT ON COLUMN uploads.av_status IS 'Column to track the anti-virus status for s3.';