UPDATE Orders
SET Grade = REPLACE(Grade, '_', '-')
WHERE Grade LIKE '%_%';

UPDATE Orders
SET grade = 'O-1'
WHERE grade = 'O_1_ACADEMY_GRADUATE';