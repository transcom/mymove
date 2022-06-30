import React from 'react';
import { Label, Textarea } from '@trussworks/react-uswds';
import { Field } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField, DatePickerInput, DropdownInput } from 'components/form/fields';
import { dropdownInputOptions } from 'utils/formatters';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const sitLocationOptions = dropdownInputOptions(LOCATION_TYPES);

const PrimeUIShipmentUpdatePPMForm = () => {
  return (
    <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
      <h2 className={styles.sectionHeader}>Dates</h2>
      <DatePickerInput
        label="Expected Departure Date"
        id="expectedDepartureDateInput"
        name="ppmShipment.expectedDepartureDate"
      />
      <DatePickerInput label="Actual Move Date" id="actualMoveDateInput" name="ppmShipment.actualMoveDate" />
      <h2 className={styles.sectionHeader}>Origin Info</h2>
      <TextField
        label="Pickup Postal Code"
        id="pickupPostalCodeInput"
        name="ppmShipment.pickupPostalCode"
        maxLength={10}
      />
      <TextField
        label="Secondary Pickup Postal Code"
        id="secondaryPickupPostalCodeInput"
        name="ppmShipment.secondaryPickupPostalCode"
        maxLength={10}
      />
      <h2 className={styles.sectionHeader}>Destination Info</h2>
      <TextField
        label="Destination Postal Code"
        id="destinationPostalCodeInput"
        name="ppmShipment.destinationPostalCode"
        maxLength={10}
      />
      <TextField
        label="Secondary Destination Postal Code"
        id="secondaryDestinationPostalCodeInput"
        name="ppmShipment.secondaryDestinationPostalCode"
        maxLength={10}
      />
      <h2 className={styles.sectionHeader}>Storage In Transit (SIT)</h2>
      <CheckboxField label="SIT Expected" id="sitExpectedInput" name="ppmShipment.sitExpected" />
      <DropdownInput
        label="SIT Location"
        id="sitLocationInput"
        name="ppmShipment.sitLocation"
        options={sitLocationOptions}
      />
      <MaskedTextField
        label="SIT Estimated Weight (lbs)"
        id="sitEstimatedWeightInput"
        name="ppmShipment.sitEstimatedWeight"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
      />
      <DatePickerInput
        label="SIT Estimated Entry Date"
        id="sitEstimatedEntryDateInput"
        name="ppmShipment.sitEstimatedEntryDate"
      />
      <DatePickerInput
        label="SIT Estimated Departure Date"
        id="sitEstimatedDepartureDateInput"
        name="ppmShipment.sitEstimatedDepartureDate"
      />
      <h2 className={styles.sectionHeader}>Weights</h2>
      <MaskedTextField
        label="Estimated Weight (lbs)"
        id="estimatedWeightInput"
        name="ppmShipment.estimatedWeight"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
      />
      <MaskedTextField
        label="Net Weight (lbs)"
        id="netWeightInput"
        name="ppmShipment.netWeight"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
      />
      <CheckboxField label="Has Pro Gear" id="hasProGearInput" name="ppmShipment.hasProGear" />
      <MaskedTextField
        label="Pro Gear Weight (lbs)"
        id="proGearWeightInput"
        name="ppmShipment.proGearWeight"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
      />
      <MaskedTextField
        label="Spouse Pro Gear Weight (lbs)"
        id="spouseProGearWeightInput"
        name="ppmShipment.spouseProGearWeight"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
      />
      <h2 className={styles.sectionHeader}>Remarks</h2>
      <Label htmlFor="counselorRemarks">Counselor Remarks</Label>
      <Field id="counselorRemarksInput" name="counselorRemarks" as={Textarea} className={`${formStyles.remarks}`} />
    </SectionWrapper>
  );
};

export default PrimeUIShipmentUpdatePPMForm;
