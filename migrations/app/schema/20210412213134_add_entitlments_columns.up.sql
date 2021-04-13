-- add two new columns to entitlements table
ALTER TABLE entitlements
    ADD COLUMN required_medical_equipment_weight int4;
COMMENT ON COLUMN entitlements.required_medical_equipment_weight IS 'The RME (required medical equipment) weight in pounds. A Service member or a dependent who is entitled to, and
receiving, medical care authorized by 10 U.S.C. §1071-§1110. may ship medical equipment necessary for such care. The medical equipment may be shipped in the same way as HHG, but has no weight limit.
The weight of authorized medical equipment is not included in the maximum authorized HHG weight
allowance.
1. Required medical equipment does not include a modified personally owned vehicle.
2. For medical equipment to qualify for shipment under this paragraph, an appropriate.
Uniformed Services healthcare provider must certify that the equipment is necessary for medical
treatment of the Service member or the dependent who is authorized medical care under.';
ALTER TABLE entitlements
    ADD COLUMN organizational_clothing_and_individual_equipment bool;
COMMENT ON COLUMN entitlements.organizational_clothing_and_individual_equipment IS 'A yes/no field reflecting whether the customer has OCIE (organizational clothing and individual equipment) that will need to be shipped as part of their move. Government property issued to the Service
member or employee by an Agency or Service for official use. A term specific to the Army and not other services.';