-- Jira: MB-16546
-- Re-assign all moves that were normally assigned to GBLOC BKAS (Fort Liberty) to go to the BGAC queue so that the moves can receive the appropriate counseling.

UPDATE orders SET gbloc = 'BGAC' where gbloc = 'BKAS';
