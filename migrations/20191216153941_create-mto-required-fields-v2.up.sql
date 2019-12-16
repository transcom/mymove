ALTER TABLE move_task_orders
    ADD COLUMN origin_duty_station_id uuid REFERENCES duty_stations,
    ADD COLUMN destination_duty_station_id uuid REFERENCES duty_stations,
    ADD COLUMN pickup_address_id uuid REFERENCES addresses,
    ADD COLUMN destination_address_id uuid REFERENCES addresses,
    ADD COLUMN requested_pickup_date date
