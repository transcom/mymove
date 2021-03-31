-- Change name to be plural like other tables.
ALTER TABLE edi_processing
    RENAME TO edi_processings;

-- Also change constraint name that's based on table name.
ALTER TABLE edi_processings
    RENAME CONSTRAINT edi_processing_pkey TO edi_processings_key;
