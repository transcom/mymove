-- Change name to be plural like other tables.
ALTER TABLE edi_processing
    RENAME TO edi_processings;

-- Also change constraint name that's based on table name.
ALTER TABLE edi_processings
    RENAME CONSTRAINT edi_processing_pkey TO edi_processings_key;

-- Add missing comments.
COMMENT ON TABLE edi_processings IS 'Stores metrics for the processing of EDIs.';
COMMENT ON COLUMN edi_processings.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN edi_processings.edi_type IS 'The type of EDI being processed (e.g., "858", "998", etc.).';
COMMENT ON COLUMN edi_processings.num_edis_processed IS 'The number of successfully processed EDIs of the given type.';
COMMENT ON COLUMN edi_processings.process_started_at IS 'Timestamp when this processing started.';
COMMENT ON COLUMN edi_processings.process_ended_at IS 'Timestamp when this processing ended.';
COMMENT ON COLUMN edi_processings.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN edi_processings.updated_at IS 'Timestamp when the record was last updated.';
