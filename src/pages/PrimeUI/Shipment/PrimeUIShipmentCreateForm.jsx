import React from 'react';
import { Label, Textarea } from '@trussworks/react-uswds';
import { Field, useFormikContext } from 'formik';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { CheckboxField, DatePickerInput, DropdownInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import formStyles from 'styles/form.module.scss';
import { dropdownInputOptions } from 'utils/formatters';
import { LOCATION_TYPES } from 'types/sitStatusShape';

const sitLocationOptions = dropdownInputOptions(LOCATION_TYPES);

const PrimeUIShipmentCreateForm = () => {
  const { values } = useFormikContext();
  const { shipmentType } = values;
  const { sitExpected, hasProGear } = values.ppmShipment;

  const hasShipmentType = !!shipmentType;
  const isPPM = shipmentType === SHIPMENT_OPTIONS.PPM;

  return (
    <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
      <h2 className={styles.sectionHeader}>Shipment Type</h2>
      <DropdownInput
        label="Shipment type"
        name="shipmentType"
        options={dropdownInputOptions(SHIPMENT_OPTIONS)}
        id="shipmentType"
      />

      {isPPM ? (
        <>
          <h2 className={styles.sectionHeader}>Dates</h2>
          <DatePickerInput
            label="Expected Departure Date"
            id="ppmShipment.expectedDepartureDateInput"
            name="ppmShipment.expectedDepartureDate"
          />
          <h2 className={styles.sectionHeader}>Origin Info</h2>
          <TextField
            label="Pickup Postal Code"
            id="ppmShipment.pickupPostalCodeInput"
            name="ppmShipment.pickupPostalCode"
            maxLength={10}
          />
          <TextField
            label="Secondary Pickup Postal Code"
            id="ppmShipment.secondaryPickupPostalCodeInput"
            name="ppmShipment.secondaryPickupPostalCode"
            maxLength={10}
          />
          <h2 className={styles.sectionHeader}>Destination Info</h2>
          <TextField
            label="Destination Postal Code"
            id="ppmShipment.destinationPostalCodeInput"
            name="ppmShipment.destinationPostalCode"
            maxLength={10}
          />
          <TextField
            label="Secondary Destination Postal Code"
            id="ppmShipment.secondaryDestinationPostalCodeInput"
            name="ppmShipment.secondaryDestinationPostalCode"
            maxLength={10}
          />
          <h2 className={styles.sectionHeader}>Storage In Transit (SIT)</h2>
          <CheckboxField label="SIT Expected" id="ppmShipment.sitExpectedInput" name="ppmShipment.sitExpected" />
          {sitExpected && (
            <>
              <DropdownInput
                label="SIT Location"
                id="ppmShipment.sitLocationInput"
                name="ppmShipment.sitLocation"
                options={sitLocationOptions}
              />
              <MaskedTextField
                label="SIT Estimated Weight (lbs)"
                id="ppmShipment.sitEstimatedWeightInput"
                name="ppmShipment.sitEstimatedWeight"
                mask={Number}
                scale={0} // digits after point, 0 for integers
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
              />
              <DatePickerInput
                label="SIT Estimated Entry Date"
                id="ppmShipment.sitEstimatedEntryDateInput"
                name="ppmShipment.sitEstimatedEntryDate"
              />
              <DatePickerInput
                label="SIT Estimated Departure Date"
                id="ppmShipment.sitEstimatedDepartureDateInput"
                name="ppmShipment.sitEstimatedDepartureDate"
              />
            </>
          )}
          <h2 className={styles.sectionHeader}>Weights</h2>
          <MaskedTextField
            label="Estimated Weight (lbs)"
            id="ppmShipment.estimatedWeightInput"
            name="ppmShipment.estimatedWeight"
            mask={Number}
            scale={0} // digits after point, 0 for integers
            signed={false} // disallow negative
            thousandsSeparator=","
            lazy={false} // immediate masking evaluation
          />
          <CheckboxField label="Has Pro Gear" id="ppmShipment.hasProGearInput" name="ppmShipment.hasProGear" />
          {hasProGear && (
            <>
              <MaskedTextField
                label="Pro Gear Weight (lbs)"
                id="ppmShipment.proGearWeightInput"
                name="ppmShipment.proGearWeight"
                mask={Number}
                scale={0} // digits after point, 0 for integers
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
              />
              <MaskedTextField
                label="Spouse Pro Gear Weight (lbs)"
                id="ppmShipment.spouseProGearWeightInput"
                name="ppmShipment.spouseProGearWeight"
                mask={Number}
                scale={0} // digits after point, 0 for integers
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
              />
            </>
          )}
          <h2 className={styles.sectionHeader}>Remarks</h2>
          <Label htmlFor="counselorRemarksInput">Counselor Remarks</Label>
          <Field id="counselorRemarksInput" name="counselorRemarks" as={Textarea} className={`${formStyles.remarks}`} />
        </>
      ) : (
        hasShipmentType && (
          <>
            <h2 className={styles.sectionHeader}>Shipment Dates</h2>
            <DatePickerInput name="requestedPickupDate" label="Requested pickup" />

            <h2 className={styles.sectionHeader}>Diversion</h2>
            <CheckboxField id="diversion" name="diversion" label="Diversion" />

            <h2 className={styles.sectionHeader}>Shipment Weights</h2>

            <MaskedTextField
              data-testid="estimatedWeightInput"
              defaultValue="0"
              name="estimatedWeight"
              label="Estimated weight (lbs)"
              id="estimatedWeightInput"
              mask={Number}
              scale={0} // digits after point, 0 for integers
              signed={false} // disallow negative
              thousandsSeparator=","
              lazy={false} // immediate masking evaluation
            />

            <h2 className={styles.sectionHeader}>Shipment Addresses</h2>
            <h5 className={styles.sectionHeader}>Pickup Address</h5>
            <AddressFields name="pickupAddress" />
            <h5 className={styles.sectionHeader}>Destination Address</h5>
            <AddressFields name="destinationAddress" />
          </>
        )
      )}
    </SectionWrapper>
  );
};

PrimeUIShipmentCreateForm.propTypes = {};

PrimeUIShipmentCreateForm.defaultProps = {};

export default PrimeUIShipmentCreateForm;
