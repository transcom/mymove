-- Prep app param table for json storage
ALTER TABLE application_parameters
ADD COLUMN IF NOT EXISTS parameter_json JSONB;

-- Insert one-off student travel app param value for weight limits
INSERT INTO application_parameters (id, parameter_name, parameter_json)
VALUES (
        '4BEEAE29-C074-4CB6-B4AE-F222F755733C',
        'studentTravelHhgAllowance',
        '{
        "TotalWeightSelf": 350,
        "TotalWeightSelfPlusDependents": 350,
        "ProGearWeight": 0,
        "ProGearWeightSpouse": 0
        }'::jsonb
    );