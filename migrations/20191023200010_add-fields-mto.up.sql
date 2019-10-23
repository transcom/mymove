CREATE TYPE affiliation AS ENUM (
    'ARMY',
    'COAST_GUARD',
    'MARINES',
    'NAVY'
);
-- TODO maybe this should be a table with branch and rank together since some ranks aren't
-- TODO possible for some branches?
CREATE TYPE rank AS ENUM (
    'E_1',
    'E_2',
    'E_3',
    'E_4',
    'E_5',
    'E_6',
    'E_7',
    'E_8',
    'E_9',
    'O_1_ACADEMY_GRADUATE',
    'O_2',
    'O_3',
    'O_4',
    'O_5',
    'O_6',
    'O_7',
    'O_8',
    'O_9',
    'O_10',
    'W_1',
    'W_2',
    'W_3',
    'W_4',
    'W_5',
    'ACADEMY_CADET',
    'AVIATION_CADET',
    'CIVILIAN_EMPLOYEE',
    'MIDSHIPMAN'
);

ALTER TABLE move_task_orders
    -- TODO will there still be a concept of a move
    -- TODO and do some of these belong there?
    ADD COLUMN customer uuid REFERENCES service_members,
    ADD COLUMN origin_duty_station uuid REFERENCES duty_stations,
    ADD COLUMN destination_duty_station uuid REFERENCES duty_stations,
    ADD COLUMN pickup_address uuid REFERENCES addresses,
    ADD COLUMN destination_address uuid REFERENCES addresses,
    ADD COLUMN requested_pickup_dates date,
    ADD COLUMN customer_remarks text,
    ADD COLUMN weight_entitlement int,
    ADD COLUMN sit_entitlement int,
    ADD COLUMN pov_entitlement bool,
    ADD COLUMN nts_entitlement bool;
