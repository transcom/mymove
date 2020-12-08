CREATE TABLE zip3_distances
(
    id UUID NOT NULL
        CONSTRAINT zip3_distances_pkey PRIMARY KEY,
    from_zip3 CHAR(3) NOT NULL,
    to_zip3 CHAR(3) NOT NULL,
    distance_miles INTEGER NOT NULL,
    -- Defaulting these to NOW() to reduce size of input file since there are so many rows.
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    CONSTRAINT zip3_distances_ordering CHECK (from_zip3 < to_zip3)
);

CREATE UNIQUE INDEX zip3_distances_unique_zip3s ON zip3_distances (from_zip3, to_zip3);

COMMENT ON TABLE zip3_distances IS 'Stores the distances between zip3 pairs; there should only be one record for any zip3 pair, with from_zip3 always alphabetically before to_zip3.';
COMMENT ON COLUMN zip3_distances.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN zip3_distances.from_zip3 IS 'The starting zip3; this should always be alphabetically before from_zip3.';
COMMENT ON COLUMN zip3_distances.to_zip3 IS 'The ending zip3; this should always be alphabetically after to_zip3.';
COMMENT ON COLUMN zip3_distances.distance_miles IS 'The distance in miles between the from_zip3 and to_zip3.';
COMMENT ON COLUMN zip3_distances.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN zip3_distances.updated_at IS 'Timestamp when the record was first updated.';
COMMENT ON CONSTRAINT zip3_distances_ordering ON zip3_distances IS 'Ensures that from_zip3 is always alphabetically before to_zip3.';
COMMENT ON INDEX zip3_distances_unique_zip3s IS 'Ensures that we have only one entry for a from/to zip3 pair.';
