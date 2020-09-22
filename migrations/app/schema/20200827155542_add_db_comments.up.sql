-- access_codes
COMMENT ON TABLE access_codes IS 'Holds access code needed to log into the customer app. Access codes prevent customers who were not part of the "trial" from submitting moves. The access codes are given to customers by the office';
COMMENT ON COLUMN access_codes.service_member_id IS 'A foreign key that points to the service_members table';
COMMENT ON COLUMN access_codes.code IS 'A code used to allow a customer access to submit a move';
COMMENT ON COLUMN access_codes.move_type IS 'Date & time the access code was created';
COMMENT ON COLUMN access_codes.created_at IS 'Date & time the access code was last updated';
COMMENT ON COLUMN access_codes.claimed_at IS 'Date & time the access code was used';

-- addresses
COMMENT ON TABLE addresses IS 'Holds all address information';
COMMENT ON COLUMN addresses.street_address_1 IS 'First street address value for address record';
COMMENT ON COLUMN addresses.street_address_2 IS 'Second street address value for address record';
COMMENT ON COLUMN addresses.city IS 'City value for address record';
COMMENT ON COLUMN addresses.state IS 'State value for address record';
COMMENT ON COLUMN addresses.postal_code IS 'Postal code value for address record';
COMMENT ON COLUMN addresses.created_at IS 'Date & time the address was created';
COMMENT ON COLUMN addresses.updated_at IS 'Date & time the address was last updated';
COMMENT ON COLUMN addresses.street_address_3 IS 'Third street address value for address record';
COMMENT ON COLUMN addresses.country IS 'Country address value for address record';

-- backup_contacts
COMMENT ON TABLE backup_contacts IS 'Holds all information regarding a backup contact for the customer';
COMMENT ON COLUMN backup_contacts.service_member_id IS 'A foreign key that points to the service_members table';
COMMENT ON COLUMN backup_contacts.name IS 'The name of the backup contact';
COMMENT ON COLUMN backup_contacts.email IS 'The email of the backup contact';
COMMENT ON COLUMN backup_contacts.phone IS 'The phone number of the backup contact';
COMMENT ON COLUMN backup_contacts.permission IS 'An enum with 3 possible values: None, View, Edit. Meanings: None: can contact only, View: can view all move details, Edit: can view and edit all move details';
COMMENT ON COLUMN backup_contacts.created_at IS 'Date & time the backup contacts was created';
COMMENT ON COLUMN backup_contacts.updated_at IS 'Date & time the backup contacts was last updated';

-- documents
COMMENT ON TABLE documents IS 'Holds information about uploaded documents';
COMMENT ON COLUMN documents.created_at IS 'Date & time the document was created';
COMMENT ON COLUMN documents.updated_at IS 'Date & time the document was last updated';
COMMENT ON COLUMN documents.service_member_id IS 'A foreign key that points to the service_members table';
COMMENT ON COLUMN documents.deleted_at IS 'Date & time document was deleted';

-- duty_station_names
COMMENT ON TABLE duty_station_names IS 'Holds information regarding alternate names for a duty station (used for duty station lookups)';
COMMENT ON COLUMN duty_station_names.name IS 'Any alternate name for a duty station other than the official name (common names, abbreviations, etc)';
COMMENT ON COLUMN duty_station_names.duty_station_id IS 'A foreign key that points to the duty stations table';
COMMENT ON COLUMN duty_station_names.created_at IS 'Date & time the duty station name was created';
COMMENT ON COLUMN duty_station_names.updated_at IS 'Date & time the duty station name was last updated';


-- duty_stations
COMMENT ON TABLE duty_stations IS 'Holds information about the duty stations';
COMMENT ON COLUMN duty_stations.name IS 'The name of the duty station';
COMMENT ON COLUMN duty_stations.affiliation IS 'The affiliation of the duty station (Army, Air Force, Navy, Marines, Coast Guard';
COMMENT ON COLUMN duty_stations.address_id IS 'A foreign key that points to the address table';
COMMENT ON COLUMN duty_stations.created_at IS 'Date & time the duty station was created';
COMMENT ON COLUMN duty_stations.updated_at IS 'Date & time the duty station was last updated';
COMMENT ON COLUMN duty_stations.transportation_office_id IS 'A foreign key that points to the transportation_offices table';

-- orders
COMMENT ON COLUMN orders.has_dependents IS 'Does the customer''s orders include any dependents?';
COMMENT ON COLUMN orders.created_at IS 'Date & time the orders were created';
COMMENT ON COLUMN orders.updated_at IS 'Date & time the orders were last updated';
COMMENT ON COLUMN orders.uploaded_orders_id IS 'A foreign key that points to the document table';
COMMENT ON COLUMN orders.status IS 'Date & time the address was last updated';
COMMENT ON COLUMN orders.department_indicator IS 'Name of the service branch. NAVY_AND_MARINES, ARMY, AIR_FORCE, COAST_GUARD';
COMMENT ON COLUMN orders.spouse_has_pro_gear IS 'Does the spouse have any pro-gear';
COMMENT ON COLUMN orders.sac IS 'Shipment Account Classification - used for accounting';
COMMENT ON COLUMN orders.confirmation_number IS 'This column is not used and should be deleted';
COMMENT ON COLUMN orders.entitlement_id IS 'A foreign key that points to the entitlements table';

-- notifications
COMMENT ON TABLE notifications IS 'Holds information about the notifications (emails) sent to customers';
COMMENT ON COLUMN notifications.service_member_id IS 'A foreign key that points to the service_members table';
COMMENT ON COLUMN notifications.ses_message_id IS 'Uuid returned after a successful sent email message';
COMMENT ON COLUMN notifications.notification_type IS 'The type of notification sent to the customer including: move approved, move canceled, move reviewed, move submitted, and payment reminder';
COMMENT ON COLUMN notifications.created_at IS 'Date & time the notification was created';

-- personally_procured_moves
COMMENT ON TABLE personally_procured_moves IS 'Holds information about the personally procured moves - moves when customers move themselves';
COMMENT ON COLUMN personally_procured_moves.move_id IS 'A foreign key that points to the moves table';
COMMENT ON COLUMN personally_procured_moves.size IS 'The size of a move: Large, Medium, Small';
COMMENT ON COLUMN personally_procured_moves.weight_estimate IS 'The estimated weight the customer think they will move';
COMMENT ON COLUMN personally_procured_moves.created_at IS 'Date & time the personally procured move was created';
COMMENT ON COLUMN personally_procured_moves.updated_at IS 'Date & time the personally procured move was last updated';
COMMENT ON COLUMN personally_procured_moves.pickup_postal_code IS 'The pickup (origin) zip entered during the PPM setup process. This zip is used for pricing';
COMMENT ON COLUMN personally_procured_moves.additional_pickup_postal_code IS 'An additional zipcode if the customer needs to pick up items from another location - an office perhaps';
COMMENT ON COLUMN personally_procured_moves.destination_postal_code IS 'The destination zipcode, which is currently the zip of the destination duty station. This zip is used for pricing';
COMMENT ON COLUMN personally_procured_moves.days_in_storage IS 'Number of days that a customer will put their things in temporary storage - max of 90 days';
COMMENT ON COLUMN personally_procured_moves.status IS 'The status of the personally procured move. Values can be: DRAFT, SUBMITTED, APPROVED, COMPLETED, CANCELED, PAYMENT_REQUESTED';
COMMENT ON COLUMN personally_procured_moves.has_additional_postal_code IS 'A boolean to determine if the user will have an additional postal code';
COMMENT ON COLUMN personally_procured_moves.has_sit IS 'A boolean to determine if the user wants to use storage in transit';
COMMENT ON COLUMN personally_procured_moves.has_requested_advance IS 'A Boolean to determine if the requested an advance';
COMMENT ON COLUMN personally_procured_moves.advance_id IS 'A foreign key that points to the reimbursements table';
COMMENT ON COLUMN personally_procured_moves.estimated_storage_reimbursement IS 'The estimated value of the SIT reimbursements from the rate engine';
COMMENT ON COLUMN personally_procured_moves.mileage IS 'The mileage between the pickup postal code and destination postal code';
COMMENT ON COLUMN personally_procured_moves.planned_sit_max IS 'The maximum SIT reimbursement for the planned SIT duration';
COMMENT ON COLUMN personally_procured_moves.sit_max IS 'Maximum SIT reimbursement for maximum SIT duration. Typically 90 days';
COMMENT ON COLUMN personally_procured_moves.incentive_estimate_min IS 'The minimum of the estimate range returned from  the rate engine';
COMMENT ON COLUMN personally_procured_moves.incentive_estimate_max IS 'The maximum of the estimate range returned from the rate engine';
COMMENT ON COLUMN personally_procured_moves.advance_worksheet_id IS 'A foreign key that points to the documents table';
COMMENT ON COLUMN personally_procured_moves.net_weight IS 'Total weight moved (actual). This number is the sum of (total weight - empty weight) for all weight tickets.';
COMMENT ON COLUMN personally_procured_moves.original_move_date IS 'The date the customer plans to move';
COMMENT ON COLUMN personally_procured_moves.actual_move_date IS 'The actual date the customer moved';
COMMENT ON COLUMN personally_procured_moves.total_sit_cost IS 'The total cost of SIT returned from rate engine';
COMMENT ON COLUMN personally_procured_moves.submit_date IS 'Date & time the customer submitted the PPM';
COMMENT ON COLUMN personally_procured_moves.approve_date IS 'Date & time the office user approved a customer''s PPM';
COMMENT ON COLUMN personally_procured_moves.reviewed_date IS 'Date & time the office user reviewed weight tickets and expenses entered by the customer';
COMMENT ON COLUMN personally_procured_moves.has_pro_gear IS 'A boolean to indicate if the customer says they have pro-gear';
COMMENT ON COLUMN personally_procured_moves.has_pro_gear_over_thousand IS 'Does the customer have pro-gear that weighs over 1000 lbs? If so, that is handled differently and may require a visit from the PPO office';

-- reimbursements
COMMENT ON TABLE reimbursements IS 'Holds information about reimbursements to a customer';
COMMENT ON COLUMN reimbursements.requested_amount IS 'The reimbursement amount the customer is requesting in cents';
COMMENT ON COLUMN reimbursements.method_of_receipt IS 'The way the customer wants to be reimbursed: OTHER (any other payment type other than GTCC), GTCC (Govt travel charge card)';
COMMENT ON COLUMN reimbursements.status IS 'The current status of the reimbursement: DRAFT, REQUESTED, APPROVED, REJECTED, PAID';
COMMENT ON COLUMN reimbursements.requested_date IS 'Date the reimbursement was requested';
COMMENT ON COLUMN reimbursements.created_at IS 'Date & time the reimbursement was created';
COMMENT ON COLUMN reimbursements.updated_at IS 'Date & time the reimbursement was last updated';

-- service_members NOTE: using customer and not customer because that's what some existing comments use
COMMENT ON TABLE service_members IS 'Holds information about a customer';
COMMENT ON COLUMN service_members.user_id IS 'A foreign key that points to the users table';
COMMENT ON COLUMN service_members.rank IS 'The customer''s rank';
COMMENT ON COLUMN service_members.middle_name IS 'The customer''s middle name';
COMMENT ON COLUMN service_members.suffix IS 'The customer''s suffix';
COMMENT ON COLUMN service_members.secondary_telephone IS 'The customer''s secondary phone number';
COMMENT ON COLUMN service_members.phone_is_preferred IS 'Does the customer prefer a phone call';
COMMENT ON COLUMN service_members.email_is_preferred IS 'Does the customer prefer an email';
COMMENT ON COLUMN service_members.residential_address_id IS 'A foreign key that points to the addresses table - containing the customer''s residential address';
COMMENT ON COLUMN service_members.backup_mailing_address_id IS 'A foreign key that points to the addresses table - containing the customer''s backup mailing address';
COMMENT ON COLUMN service_members.created_at IS 'Date & time the customer was created';
COMMENT ON COLUMN service_members.updated_at IS 'Date & time the customer was last updated';
COMMENT ON COLUMN service_members.social_security_number_id IS 'A foreign key that points to the social_security_numbers table';
COMMENT ON COLUMN service_members.duty_station_id IS 'A foreign key that points to the duty station table - containing the customer''s current duty station';
COMMENT ON COLUMN service_members.requires_access_code IS 'A boolean value that controls if a customer needs to enter an access code to submit a move';

-- signed_certifications
COMMENT ON TABLE signed_certifications IS 'Holds information about when the customer signed the certificate';
COMMENT ON COLUMN signed_certifications.submitting_user_id IS 'A foreign key that points to the users table';
COMMENT ON COLUMN signed_certifications.move_id IS 'A foreign key that points to the moves table';
COMMENT ON COLUMN signed_certifications.certification_text IS 'The legalese text the customer agrees to. Value is hard coded and stored in: src/scenes/Legalese/legaleseText.js -> ppmPaymentLegal';
COMMENT ON COLUMN signed_certifications.signature IS 'Currently hard coded to, CHECKBOX, coming from the frontend';
COMMENT ON COLUMN signed_certifications.date IS 'Date & time the customer signed';
COMMENT ON COLUMN signed_certifications.created_at IS 'Date & time the notification was created';
COMMENT ON COLUMN signed_certifications.updated_at IS 'Date & time the notification was last updated';
COMMENT ON COLUMN signed_certifications.personally_procured_move_id IS 'A foreign key that points to the personally_procured_moves table';
COMMENT ON COLUMN signed_certifications.certification_type IS 'A certification type: PPM, PPM_PAYMENT, HHG';

-- social_security_numbers
COMMENT ON TABLE social_security_numbers IS 'Holds information regarding a customer''s social security number';
COMMENT ON COLUMN social_security_numbers.encrypted_hash IS 'A hashed version of a customer''s social security number';
COMMENT ON COLUMN social_security_numbers.created_at IS 'Date & time the social security number was created';
COMMENT ON COLUMN social_security_numbers.updated_at IS 'Date & time the social security number was last updated';

-- uploads
COMMENT ON TABLE uploads IS 'Holds information regarding files that are uploaded';
COMMENT ON COLUMN uploads.filename IS 'The filename of the upload';
COMMENT ON COLUMN uploads.bytes IS 'The number of bytes of the upload';
COMMENT ON COLUMN uploads.content_type IS 'The mime type of the upload';
COMMENT ON COLUMN uploads.checksum IS 'A checksum value of the upload';
COMMENT ON COLUMN uploads.created_at IS 'Date & time the uploads was created';
COMMENT ON COLUMN uploads.updated_at IS 'Date & time the uploads was last updated';
COMMENT ON COLUMN uploads.storage_key IS 'The resulting path to where the upload is on S3';
COMMENT ON COLUMN uploads.deleted_at IS 'Date & time of when the uploads was deleted';
COMMENT ON COLUMN uploads.upload_type IS 'Who created the upload: USER, PRIME';

-- user_uploads
COMMENT ON TABLE user_uploads IS 'Holds information that joins the uploads to the corresponding documents and users';
COMMENT ON COLUMN user_uploads.document_id IS 'A foreign key that points to the documents table';
COMMENT ON COLUMN user_uploads.uploader_id IS 'A foreign key that points to the users table';
COMMENT ON COLUMN user_uploads.upload_id IS 'A foreign key that points to the uploaded table';
COMMENT ON COLUMN user_uploads.created_at IS 'Date & time the user uploads was created';
COMMENT ON COLUMN user_uploads.updated_at IS 'Date & time the user uploads was last updated';
COMMENT ON COLUMN user_uploads.deleted_at IS 'Date & time the user uploads was deleted';
