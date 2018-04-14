-- Need access to a UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Import those SCACs for all those TSPs
SELECT
  uuid_generate_v4() as id,
  scac text

