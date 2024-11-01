-- Add new International SIT FSC services
INSERT INTO re_services
(id, code, name, created_at, updated_at)
VALUES
	('98ebdb87-cf2e-4d77-8630-1f0dc830fe41', 'IOSFSC', 'International origin SIT fuel surcharge', now(), now()),
	('b63b598c-81c3-44ea-852f-537aeef9c018', 'IDSFSC', 'International destination SIT fuel surcharge', now(), now());
