CREATE INDEX ON notifications (notification_type);
CREATE INDEX ON notifications (service_member_id);

CREATE INDEX ON personally_procured_moves (original_move_date);
CREATE INDEX ON personally_procured_moves (reviewed_date); -- done for post_move_survey script

CREATE INDEX ON office_phone_lines (is_dsn_number);
