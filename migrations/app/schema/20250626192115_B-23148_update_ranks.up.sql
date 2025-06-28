DELETE FROM ranks WHERE rank_abbv = 'PSG' AND affiliation = 'ARMY';

UPDATE ranks
SET rank_order = rank_order - 1
WHERE affiliation = 'ARMY' AND rank_order > 21;

UPDATE ranks
SET rank_abbv = '1SG'
WHERE id = '50f29cf0-fd75-452a-969e-a64ac19f3775';

UPDATE ranks
SET rank_order = rank_order + 1
WHERE affiliation = 'AIR_FORCE' AND rank_order > 10;

UPDATE ranks
SET rank_order = rank_order + 1
WHERE id = '43b69e7b-99a3-488f-8cea-3abe63d6f20a';

UPDATE ranks
SET rank_order = rank_order + 1
WHERE id = '2cf8e36a-20fb-41fe-9268-d3d1f0219d1a';
