-- See https://dp3.atlassian.net/wiki/spaces/MT/pages/2738716677/HHG+and+UB+Entitlements
-- Prep entitlements table for holding weight restricted
ALTER TABLE entitlements
ADD COLUMN IF NOT EXISTS is_weight_restricted boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS weight_restriction int;
-- Create pay grades table to get our static entitlements.go file to be db based
CREATE TABLE IF NOT EXISTS pay_grades (
    id uuid PRIMARY KEY NOT NULL,
    grade text NOT NULL UNIQUE,
    grade_description text,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);
-- Create household goods allowances table
CREATE TABLE IF NOT EXISTS hhg_allowances (
    id uuid PRIMARY KEY NOT NULL,
    pay_grade_id uuid NOT NULL UNIQUE REFERENCES pay_grades(id) ON DELETE CASCADE,
    total_weight_self int NOT NULL,
    total_weight_self_plus_dependents int NOT NULL,
    pro_gear_weight int NOT NULL,
    pro_gear_weight_spouse int NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);
-- Insert Max HHG allowance app value
-- camel case to match the existing standaloneCrateCap parameter
INSERT INTO application_parameters (id, parameter_name, parameter_value)
VALUES (
        'D246186B-E93B-4716-B82C-6A38EA5EAB8C',
        'maxHhgAllowance',
        '18000'
    );
-- Insert pay_grades and hhg_allowances
-- ACADEMY_CADET
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '8D8C82EA-EA8F-4D7F-9D84-8D186AB7A7C0',
        'ACADEMY_CADET',
        'Academy Cadet'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '8A43128C-B080-4D22-9BEA-6F1CBF7F7123',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'ACADEMY_CADET'
        ),
        350,
        350,
        0,
        0
    );
-- MIDSHIPMAN
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '63998729-EF74-486E-BEEA-5B519FA3812F',
        'MIDSHIPMAN',
        'Midshipman'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'C95AA341-9261-4E14-B63B-9A7262FB8EA0',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'MIDSHIPMAN'
        ),
        350,
        350,
        0,
        0
    );
-- AVIATION_CADET
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        'DF749D7E-5007-43CD-8715-2875D281F817',
        'AVIATION_CADET',
        'Aviation Cadet'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '23D6DEF4-975E-4075-A4B2-E4DC3DF3D6FF',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'AVIATION_CADET'
        ),
        7000,
        8000,
        2000,
        500
    );
-- E_1
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '6CB785D0-CABF-479A-A36D-A6AEC294A4D0',
        'E_1',
        'Enlisted Grade E_1'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '6CB785D0-CABF-479A-A36D-A6AEC294A4DE',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_1'
        ),
        5000,
        8000,
        2000,
        500
    );
-- E_2
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '5F871C82-F259-43CC-9245-A6E18975DDE0',
        'E_2',
        'Enlisted Grade E_2'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '5F871C82-F259-43CC-9245-A6E18975DDE8',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_2'
        ),
        5000,
        8000,
        2000,
        500
    );
-- E_3
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '862EB395-86D1-44AF-AD47-DEC44FBEDA30',
        'E_3',
        'Enlisted Grade E_3'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '862EB395-86D1-44AF-AD47-DEC44FBEDA3F',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_3'
        ),
        5000,
        8000,
        2000,
        500
    );
-- E_4
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        'BB55F37C-3165-46BA-AD3F-9A477F699990',
        'E_4',
        'Enlisted Grade E_4'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'BB55F37C-3165-46BA-AD3F-9A477F699991',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_4'
        ),
        7000,
        8000,
        2000,
        500
    );
-- E_5
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '3F142461-DCA5-4A77-9295-92EE93371330',
        'E_5',
        'Enlisted Grade E_5'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '3F142461-DCA5-4A77-9295-92EE9337133A',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_5'
        ),
        7000,
        9000,
        2000,
        500
    );
-- E_6
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '541AEC36-BD9F-4AD2-ABB4-D9B63E29DC80',
        'E_6',
        'Enlisted Grade E_6'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '541AEC36-BD9F-4AD2-ABB4-D9B63E29DC8C',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_6'
        ),
        8000,
        11000,
        2000,
        500
    );
-- E_7
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '523D57A1-529C-4DFD-8C33-9CB169FD29A0',
        'E_7',
        'Enlisted Grade E_7'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '523D57A1-529C-4DFD-8C33-9CB169FD29AF',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_7'
        ),
        11000,
        13000,
        2000,
        500
    );
-- E_8
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '1D909DB0-602F-4724-BD43-8F90A6660460',
        'E_8',
        'Enlisted Grade E_8'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '1D909DB0-602F-4724-BD43-8F90A666046E',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_8'
        ),
        12000,
        14000,
        2000,
        500
    );
-- E_9
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        'A5FC8FD2-6F91-492B-ABE2-2157D03EC990',
        'E_9',
        'Enlisted Grade E_9'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'A5FC8FD2-6F91-492B-ABE2-2157D03EC99B',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_9'
        ),
        13000,
        15000,
        2000,
        500
    );
-- E_9 Special Senior Enlisted
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '911208CC-3D13-49D6-9478-B0A3943435C0',
        'E_9_SPECIAL_SENIOR_ENLISTED',
        'Enlisted Grade E_9 Special Senior Enlisted'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'D219899B-251F-49E9-94B3-C073C22D9D2F',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'E_9_SPECIAL_SENIOR_ENLISTED'
        ),
        14000,
        17000,
        2000,
        500
    );
-- O_1 (Academy Graduate) / W_1 uses same as O_1
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        'B25998F4-4715-4F41-8986-4C5C8E59FC80',
        'O_1_ACADEMY_GRADUATE',
        'Officer Grade O_1 Academy Graduate'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'B25998F4-4715-4F41-8986-4C5C8E59FC84',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_1_ACADEMY_GRADUATE'
        ),
        10000,
        12000,
        2000,
        500
    );
-- O_2
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        'D1B76A01-D8E4-4BD3-98FF-FA93FF7BC790',
        'O_2',
        'Officer Grade O_2'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'D1B76A01-D8E4-4BD3-98FF-FA93FF7BC79A',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_2'
        ),
        12500,
        13500,
        2000,
        500
    );
-- O_3
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '5658D67B-D510-4226-9E56-714403BA0F10',
        'O_3',
        'Officer Grade O_3'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '5658D67B-D510-4226-9E56-714403BA0F1D',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_3'
        ),
        13000,
        14500,
        2000,
        500
    );
-- O_4
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        'E83D8F8D-F70B-4DB1-99CC-DD983D2FD250',
        'O_4',
        'Officer Grade O_4'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '0991ABBC-5400-4E6C-8BC4-195F9A602E75',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_4'
        ),
        14000,
        17000,
        2000,
        500
    );
-- O_5
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '3BC4B197-7897-4105-80A1-39A0378D7730',
        'O_5',
        'Officer Grade O_5'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'E83D8F8D-F70B-4DB1-99CC-DD983D2FD25D',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_5'
        ),
        16000,
        17500,
        2000,
        500
    );
-- O_6
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '455A112D-D1E0-4559-81E8-6DF664638F70',
        'O_6',
        'Officer Grade O_6'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '3BC4B197-7897-4105-80A1-39A0378D773E',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_6'
        ),
        18000,
        18000,
        2000,
        500
    );
-- O_7
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        'CF664124-9BAF-4187-8F28-0908C0F0A5E0',
        'O_7',
        'Officer Grade O_7'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '455A112D-D1E0-4559-81E8-6DF664638F7C',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_7'
        ),
        18000,
        18000,
        2000,
        500
    );
-- O_8
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '6E50B04A-52DC-45C9-91D9-4A7B4FA1AB20',
        'O_8',
        'Officer Grade O_8'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'CF664124-9BAF-4187-8F28-0908C0F0A5E8',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_8'
        ),
        18000,
        18000,
        2000,
        500
    );
-- O_9
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '1D6E34C3-8C6C-4D4F-8B91-F46BED3F5E80',
        'O_9',
        'Officer Grade O_9'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '6E50B04A-52DC-45C9-91D9-4A7B4FA1AB2A',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_9'
        ),
        18000,
        18000,
        2000,
        500
    );
-- O_10
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '7FA938AB-1C34-4666-A878-9B989C916D1A',
        'O_10',
        'Officer Grade O_10'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '1D6E34C3-8C6C-4D4F-8B91-F46BED3F5E85',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_10'
        ),
        18000,
        18000,
        2000,
        500
    );
-- W_1
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '6BADF8A0-B0EF-4E42-B827-7F63A3987A4B',
        'W_1',
        'Warrant Officer W_1'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '16F0F64F-728A-42A7-98B7-EA9BF289FE1A',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'W_1'
        ),
        10000,
        12000,
        2000,
        500
    );
-- W_2
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        'A687A2E1-488C-4943-B9D9-3D645A2712F4',
        'W_2',
        'Warrant Officer W_2'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'A687A2E1-488C-4943-B9D9-3D645A2712F9',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'W_2'
        ),
        12500,
        13500,
        2000,
        500
    );
-- W_3
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '5A65FB1F-4245-4178-B6A7-CC504C9CBB37',
        'W_3',
        'Warrant Officer W_3'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '5A65FB1F-4245-4178-B6A7-CC504C9CBB38',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'W_3'
        ),
        13000,
        14500,
        2000,
        500
    );
-- W_4
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '74DB5649-CF66-4AF8-939B-D3D7F1F6B7C6',
        'W_4',
        'Warrant Officer W_4'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '74DB5649-CF66-4AF8-939B-D3D7F1F6B7C7',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'W_4'
        ),
        14000,
        17000,
        2000,
        500
    );
-- W_5
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        'EA8CB0E9-15FF-43B4-9E41-7168D01E7553',
        'W_5',
        'Warrant Officer W_5'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        'EA8CB0E9-15FF-43B4-9E41-7168D01E7554',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'W_5'
        ),
        16000,
        17500,
        2000,
        500
    );
-- CIVILIAN EMPLOYEE
INSERT INTO pay_grades (id, grade, grade_description)
VALUES (
        '9E2CB9A5-ACE3-4235-9EE7-EBE4CC2A9BC9',
        'CIVILIAN_EMPLOYEE',
        'Civilian Employee'
    );
INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '9E2CB9A5-ACE3-4235-9EE7-EBE4CC2A9BC1',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'CIVILIAN_EMPLOYEE'
        ),
        18000,
        18000,
        2000,
        500
    );