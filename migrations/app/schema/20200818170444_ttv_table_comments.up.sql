COMMENT ON TABLE admin_users IS 'Holds all users who have access to the admin interface, where one can perform CRUD operations on entities such as office users and admin users. Individual authenticated sessions can also be revoked via the admin interface.';
COMMENT ON COLUMN admin_users.created_at IS 'Date & time the admin user was created';
COMMENT ON COLUMN admin_users.updated_at IS 'Date & time the admin user was updated';
COMMENT ON COLUMN admin_users.user_id IS 'The foreign key that points to the user id in the users table';
COMMENT ON COLUMN admin_users.first_name IS 'The first name of the admin user';
COMMENT ON COLUMN admin_users.last_name IS 'The last name of the admin user';
COMMENT ON COLUMN admin_users.organization_id IS 'The foreign key that points to the organization id in the organizations table. Truss admin users belong to the Truss organization.';
COMMENT ON COLUMN admin_users.role IS 'An enum with two possible values: SYSTEM_ADMIN or PROGRAM_ADMIN. Note that PROGRAM_ADMIN is no longer used and there is a JIRA story to remove it.';
COMMENT ON COLUMN admin_users.email IS 'The email of the admin user';
COMMENT ON COLUMN admin_users.active IS 'A boolean that determines whether or not an admin user is active. Users that are not active are not allowed to access the admin site. See https://github.com/transcom/mymove/wiki/create-or-deactivate-users.';

COMMENT ON TABLE contractors IS 'Holds all contractors who handle moves. There is only one active contractor per type at a time, though we do not yet have a way to identify that.';
COMMENT ON COLUMN contractors.created_at IS 'Date & time the contractor was created';
COMMENT ON COLUMN contractors.updated_at IS 'Date & time the contractor was updated';
COMMENT ON COLUMN contractors.name IS 'The name of the contractor';
COMMENT ON COLUMN contractors.contract_number IS 'The government-issued contract number for the contractor.';
COMMENT ON COLUMN contractors.type IS 'A string to represent the type of contractor. Examples are Prime and NTS.';

COMMENT ON TABLE distance_calculations IS 'Represents a distance calculation in miles between an origin and destination address.';
COMMENT ON COLUMN distance_calculations.created_at IS 'Date & time the distance_calculation was created';
COMMENT ON COLUMN distance_calculations.updated_at IS 'Date & time the distance_calculation was updated';
COMMENT ON COLUMN distance_calculations.origin_address_id IS 'Represents the origin address as a foreign key to the addresses table.';
COMMENT ON COLUMN distance_calculations.destination_address_id IS 'Represents the destination address as a foreign key to the addresses table.';
COMMENT ON COLUMN distance_calculations.distance_miles IS 'The distance in miles between the origin and destination address.';

COMMENT ON TABLE dps_users IS 'Users who have permission to access MyMove - DPS integration resources';
COMMENT ON COLUMN dps_users.created_at IS 'Date & time the dps_user was created';
COMMENT ON COLUMN dps_users.updated_at IS 'Date & time the dps_user was updated';
COMMENT ON COLUMN dps_users.login_gov_email IS 'The login.gov email of the user.';
COMMENT ON COLUMN dps_users.active IS 'A boolean that determines whether or not a DPS user is active. Users that are not active are not allowed to access the DPS resources. See https://github.com/transcom/mymove/wiki/create-or-deactivate-users.';

COMMENT ON TABLE ghc_domestic_transit_times IS 'Allows calculation of the maximum transit time based on the distance and weight ranges.';
COMMENT ON COLUMN ghc_domestic_transit_times.max_days_transit_time IS 'The max transit time for the corresponding weight and distance ranges defined via the _lower and _upper columns.';
COMMENT ON COLUMN ghc_domestic_transit_times.weight_lbs_lower IS 'The minimum weight in the range.';
COMMENT ON COLUMN ghc_domestic_transit_times.weight_lbs_upper IS 'The maximum weight in the range. If 0 (zero), there is no upper bound';
COMMENT ON COLUMN ghc_domestic_transit_times.distance_miles_lower IS 'The minimum distance in the range.';
COMMENT ON COLUMN ghc_domestic_transit_times.distance_miles_upper IS 'The maximum distance in the range.';

COMMENT ON TABLE office_emails IS 'Stores email addresses for the Transportation Offices.';
COMMENT ON COLUMN office_emails.created_at IS 'Date & time the office_email was created.';
COMMENT ON COLUMN office_emails.updated_at IS 'Date & time the office_email was updated.';
COMMENT ON COLUMN office_emails.transportation_office_id IS 'A foreign key to the transportation_offices table.';
COMMENT ON COLUMN office_emails.email IS 'The email address for the transportation office.';
COMMENT ON COLUMN office_emails.label IS 'The department the email gets sent to. For example, ''Customer Service''';

COMMENT ON TABLE office_phone_lines IS 'Stores phone numbers for the Transportation Offices.';
COMMENT ON COLUMN office_phone_lines.created_at IS 'Date & time the office_phone_line was created.';
COMMENT ON COLUMN office_phone_lines.updated_at IS 'Date & time the office_phone_line was updated.';
COMMENT ON COLUMN office_phone_lines.transportation_office_id IS 'A foreign key to the transportation_offices table.';
COMMENT ON COLUMN office_phone_lines.number IS 'The phone number for the transportation office.';
COMMENT ON COLUMN office_phone_lines.type IS 'The kind of phone line, such as ''voice'' or ''fax''';
COMMENT ON COLUMN office_phone_lines.label IS 'This field is not populated locally. It''s not clear how it differs from type';
COMMENT ON COLUMN office_phone_lines.is_dsn_number IS 'A boolean that represents whether or not this number is a Defense Switched Network number. Defaults to false.';

COMMENT ON TABLE office_users IS 'Holds all users who have access to the office site.';
COMMENT ON COLUMN office_users.created_at IS 'Date & time the office user was created.';
COMMENT ON COLUMN office_users.updated_at IS 'Date & time the office user was updated.';
COMMENT ON COLUMN office_users.user_id IS 'The foreign key that points to the user id in the users table. This gets populated when the user first signs in via login.gov, which then creates the user in the users table, and the link is then made in this table.';
COMMENT ON COLUMN office_users.first_name IS 'The first name of the office user.';
COMMENT ON COLUMN office_users.last_name IS 'The last name of the office user.';
COMMENT ON COLUMN office_users.middle_initials IS 'The middle initials of the office user.';
COMMENT ON COLUMN office_users.email IS 'The email of the office user. This will match their login_gov_email in the users table.';
COMMENT ON COLUMN office_users.telephone IS 'The phone number of the office user.';
COMMENT ON COLUMN office_users.transportation_office_id IS 'The id of the transportation office the office user is assigned to.';
COMMENT ON COLUMN office_users.active IS 'A boolean that determines whether or not an office user is active. Users that are not active are not allowed to access the office site. See https://github.com/transcom/mymove/wiki/create-or-deactivate-users.';

COMMENT ON TABLE organizations IS 'Holds all organizations that admin users belong to.';
COMMENT ON COLUMN organizations.created_at IS 'Date & time the organization was created.';
COMMENT ON COLUMN organizations.updated_at IS 'Date & time the organization was updated.';
COMMENT ON COLUMN organizations.name IS 'The organization name.';
COMMENT ON COLUMN organizations.poc_email IS 'The email of the organization''s point of contact.';
COMMENT ON COLUMN organizations.poc_phone IS 'The phone number of the organization''s point of contact.';

COMMENT ON TABLE roles IS 'Holds all roles that users can have.';
COMMENT ON COLUMN roles.created_at IS 'Date & time the role was created.';
COMMENT ON COLUMN roles.updated_at IS 'Date & time the role was updated.';
COMMENT ON COLUMN roles.role_type IS 'The name of the role in snake case. Current values are: ''transportation_ordering_officer'', ''transportation_invoicing_officer'', ''customer'', ''ppm_office_users'', ''contracting_officer''.';
COMMENT ON COLUMN roles.role_name IS 'The reader-friendly capitalized name of the role.';

COMMENT ON TABLE users_roles IS 'A join table between users and roles to identify which users have which roles.';
COMMENT ON COLUMN users_roles.created_at IS 'Date & time the users_roles was created.';
COMMENT ON COLUMN users_roles.updated_at IS 'Date & time the users_roles was updated.';
COMMENT ON COLUMN users_roles.deleted_at IS 'Date & time the users_roles was deleted.';
COMMENT ON COLUMN users_roles.user_id IS 'The id of the user being referenced.';
COMMENT ON COLUMN users_roles.role_id IS 'The id of the role being referenced.';

COMMENT ON TABLE transportation_offices IS 'Holds all known transportation offices where office users are assigned.';
COMMENT ON COLUMN transportation_offices.created_at IS 'Date & time the transportation_office was created.';
COMMENT ON COLUMN transportation_offices.updated_at IS 'Date & time the transportation_office was updated.';
COMMENT ON COLUMN transportation_offices.shipping_office_id IS 'This is a foreign key that points back to this table. This does not seem right and will be removed in a separate cleanup PR.';
COMMENT ON COLUMN transportation_offices.name IS 'The name of the transportation office.';
COMMENT ON COLUMN transportation_offices.address_id IS 'The id of the transportation office''s address from the addresses table.';
COMMENT ON COLUMN transportation_offices.latitude IS 'The latitude of the transportation office.';
COMMENT ON COLUMN transportation_offices.longitude IS 'The longitude of the transportation office.';
COMMENT ON COLUMN transportation_offices.hours IS 'The hours of operation in freeform text format.';
COMMENT ON COLUMN transportation_offices.services IS 'The various services offered in freeform text format.';
COMMENT ON COLUMN transportation_offices.note IS 'Unclear what this field is used for. It is not populated locally.';
COMMENT ON COLUMN transportation_offices.gbloc IS 'A 4-character code representing the geographical area this transportation office is part of. This maps to the code field in the jppso_regions table.';

COMMENT ON TABLE users IS 'Holds all users. Anyone who signs in to any of the mymove apps is automatically created in this table after signing in with login.gov.';
COMMENT ON COLUMN users.created_at IS 'Date & time the user was created.';
COMMENT ON COLUMN users.updated_at IS 'Date & time the user was updated.';
COMMENT ON COLUMN users.login_gov_uuid IS 'The login.gov uuid of the user.';
COMMENT ON COLUMN users.login_gov_email IS 'The login.gov email of the user.';
COMMENT ON COLUMN users.active IS 'A boolean that determines whether or not a user is active. Users that are not active are not allowed to access the mymove apps. See https://github.com/transcom/mymove/wiki/create-or-deactivate-users.';
COMMENT ON COLUMN users.current_mil_session_id IS 'This field gets populated when a user signs into the mil app. The string matches the session id stored in Redis. It is used to allow an admin user to revoke the session if necessary.';
COMMENT ON COLUMN users.current_admin_session_id IS 'This field gets populated when a user signs into the admin app. The string matches the session id stored in Redis. It is used to allow an admin user to revoke the session if necessary.';
COMMENT ON COLUMN users.current_office_session_id IS 'This field gets populated when a user signs into the office app. The string matches the session id stored in Redis. It is used to allow an admin user to revoke the session if necessary.';

COMMENT ON TABLE weight_ticket_set_documents IS 'Documents the vehicles used to transport goods: their type, make, model, empty weight, full weight, date the weight was measured, and whether or not documentation is missing.';
COMMENT ON COLUMN weight_ticket_set_documents.created_at IS 'Date & time the weight_ticket_set_document was created.';
COMMENT ON COLUMN weight_ticket_set_documents.updated_at IS 'Date & time the weight_ticket_set_document was updated.';
COMMENT ON COLUMN weight_ticket_set_documents.deleted_at IS 'Date & time the weight_ticket_set_document was deleted.';
COMMENT ON COLUMN weight_ticket_set_documents.weight_ticket_set_type IS 'An enum with 4 possible values: CAR, CAR_TRAILER, BOX_TRUCK, PRO_GEAR';
COMMENT ON COLUMN weight_ticket_set_documents.move_document_id IS 'The id of the move_document this weight_ticket_set_document is associated with.';
COMMENT ON COLUMN weight_ticket_set_documents.empty_weight IS 'The empty weight in pounds.';
COMMENT ON COLUMN weight_ticket_set_documents.empty_weight_ticket_missing IS 'A boolean representing whether or not the empty weight ticket is missing.';
COMMENT ON COLUMN weight_ticket_set_documents.full_weight IS 'The full weight in pounds.';
COMMENT ON COLUMN weight_ticket_set_documents.full_weight_ticket_missing IS 'A boolean representing whether or not the full weight ticket is missing.';
COMMENT ON COLUMN weight_ticket_set_documents.weight_ticket_date IS 'The date the weight was measured and recorded on a ticket.';
COMMENT ON COLUMN weight_ticket_set_documents.trailer_ownership_missing IS 'A boolean representing whether or not the trailer ownership is missing.';
COMMENT ON COLUMN weight_ticket_set_documents.vehicle_make IS 'The make of the vehicle used to transport the goods.';
COMMENT ON COLUMN weight_ticket_set_documents.vehicle_model IS 'The model of the vehicle used to transport the goods.';
COMMENT ON COLUMN weight_ticket_set_documents.vehicle_nickname IS 'The nickname of the vehicle used to transport the goods.';

COMMENT ON TABLE jppso_regions IS 'Holds all JPPSO region names and codes. This is used to map states to regions. This table is not currently used, but will be soon in order to associate a TOO with a specific JPPSO, which will allow the TOO to filter the list of moves by region.';
COMMENT ON COLUMN jppso_regions.created_at IS 'Date & time the jppso_region was created.';
COMMENT ON COLUMN jppso_regions.updated_at IS 'Date & time the jppso_region was updated.';
COMMENT ON COLUMN jppso_regions.code IS 'The 4-character code for the region.';
COMMENT ON COLUMN jppso_regions.name IS 'The human-readable name of the region.';

COMMENT ON TABLE jppso_region_state_assignments IS 'Maps US states to JPPSO regions. This table is not currently used, but will be soon in order to associate a TOO with a specific JPPSO, which will allow the TOO to filter the list of moves by region.';
COMMENT ON COLUMN jppso_region_state_assignments.created_at IS 'Date & time the jppso_region_state_assignment was created.';
COMMENT ON COLUMN jppso_region_state_assignments.updated_at IS 'Date & time the jppso_region_state_assignment was updated.';
COMMENT ON COLUMN jppso_region_state_assignments.jppso_region_id IS 'The JPPSO region this state is part of. A foreign key to the jppso_regions table.';
COMMENT ON COLUMN jppso_region_state_assignments.state_name IS 'The full capitalized US state name.';
COMMENT ON COLUMN jppso_region_state_assignments.state_abbreviation IS 'The two-letter state abbreviation.';
