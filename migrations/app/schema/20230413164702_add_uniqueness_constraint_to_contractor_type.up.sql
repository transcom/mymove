ALTER TABLE contractors
  ADD CONSTRAINT unique_contractors_type UNIQUE (type);
