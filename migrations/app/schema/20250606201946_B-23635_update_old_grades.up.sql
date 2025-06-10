UPDATE Orders
SET Grade = REPLACE(Grade, '_', '-')
WHERE (Grade LIKE 'E_%' OR Grade LIKE 'W_%' OR Grade LIKE 'O_%');

UPDATE Orders
SET grade = 'O-1'
WHERE grade = 'O-1-ACADEMY-GRADUATE';