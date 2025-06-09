--B-23635   Jonathan Spight  ub allowances table

ALTER TABLE ub_allowances
DROP CONSTRAINT IF EXISTS ub_allowances_grade_check;

-- these must be executed in this order or the updated constraint will fail because the current values no longer align
UPDATE ub_allowances
SET grade = 'O-1'
WHERE grade = 'O_1_ACADEMY_GRADUATE';

UPDATE ub_allowances
SET grade = REPLACE(grade, '_', '-')
WHERE grade LIKE '%\_%';


ALTER TABLE ub_allowances
ADD CONSTRAINT grade CHECK (
  grade IN (
    'E-1', 'E-2', 'E-3', 'E-4', 'E-5', 'E-6', 'E-7', 'E-8', 'E-9',
    'O-1', 'O-2', 'O-3', 'O-4', 'O-5', 'O-6', 'O-7', 'O-8', 'O-9', 'O-10',
    'W-1', 'W-2', 'W-3', 'W-4', 'W-5'
  )
);
