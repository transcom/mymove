update
	transportation_offices
set
	gbloc = 'BGAC'
where
	gbloc = 'BKAS';

insert
	into
	duty_locations(id,
	name,
	address_id,
	created_at,
	updated_at,
	transportation_office_id,
	provides_services_counseling )
values(
    '9f7a1c83-26ad-4ba3-b30f-1a2e62f250f6',
    'PPPO West Point/ USMA - USA',
'6aa77b74-41a7-4a4c-ab29-986f3263495a',
now(),
now(),
'dd043073-4f1b-460f-8f8c-74403619dbaa',
true);

insert
	into duty_location_names(id,
	name,
	duty_location_id,
	created_at,
	updated_at)
values (
    '47cfbfa1-4633-440e-928e-92a6f462826e',
    'PPPO West Point/ USMA - USA',
'9f7a1c83-26ad-4ba3-b30f-1a2e62f250f6',
now(),
now()
);

insert
	into
	duty_locations(id,
	name,
	address_id,
	created_at,
	updated_at,
	transportation_office_id,
	provides_services_counseling )
values(
    '0f420f7b-72ac-43cf-bc93-f24d44ba8f93',
    'PPPO USAG Miami - USA',
'09058d36-2966-496a-aaf5-55c024404396',
now(),
now(),
'4f10d0f5-6017-4de2-8cfb-ee9252e492d5',
true);

insert
	into duty_location_names(id,
	name,
	duty_location_id,
	created_at,
	updated_at)
values (
    'fbe7677f-e1ec-47d2-bf33-57e83eded778',
    'PPPO USAG Miami - USA',
'0f420f7b-72ac-43cf-bc93-f24d44ba8f93',
now(),
now()
);
