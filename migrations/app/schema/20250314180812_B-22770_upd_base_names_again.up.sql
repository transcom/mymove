--update duty loc by name in case uuids don't match in all envs
update duty_locations set name = 'Fort Bragg, NC 28307', address_id = 'c13715ec-68d9-4c77-ae9a-5a652ddd3787', updated_at = now() where name = 'Fort Liberty, NC 28307';

update duty_locations set name = 'Fort Bragg, NC 28310', address_id = 'd8769fb0-e130-46b2-9191-509663bab4b4', updated_at = now() where name = 'Fort Liberty, NC 28310';