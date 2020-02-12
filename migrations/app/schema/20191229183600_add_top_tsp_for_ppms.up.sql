-- Create separate placeholder TSP to use as the top TSP when we don't have specific TSP data
INSERT INTO transportation_service_providers
(id, standard_carrier_alpha_code, created_at, updated_at, enrolled, name)
VALUES ('55f9eb49-e1bc-4b69-9a98-73642d2b504c', 'PPMX', now(), now(), false, 'Top TSP for PPMs');
