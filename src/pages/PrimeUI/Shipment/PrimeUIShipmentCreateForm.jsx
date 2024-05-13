import React, { useState } from 'react';
import { Radio, FormGroup, Label, Textarea } from '@trussworks/react-uswds';
import { Field, useField, useFormikContext } from 'formik';

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
  const { sitExpected, hasProGear, hasSecondaryDestinationAddress, hasSecondaryPickupAddress } = values.ppmShipment;
  const [, , checkBoxHelperProps] = useField('diversion');
  const [, , divertedFromIdHelperProps] = useField('divertedFromShipmentId');
  const [isChecked, setIsChecked] = useState(false);

  const hasShipmentType = !!shipmentType;
  const isPPM = shipmentType === SHIPMENT_OPTIONS.PPM;

  // if a shipment is a diversion, then the parent shipment id will be required for input
  const toggleParentShipmentIdTextBox = (checkboxValue) => {
    if (checkboxValue) {
      checkBoxHelperProps.setValue(true);
      setIsChecked(true);
    } else {
      // set values for checkbox & divertedFromShipmentId when box is unchecked
      checkBoxHelperProps.setValue('');
      divertedFromIdHelperProps.setValue('');
      setIsChecked(false);
    }
  };

  // validates the uuid - if it ain't no good, then an error message displays
  const validateUUID = (value) => {
    const uuidRegex = /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/;

    if (!uuidRegex.test(value)) {
      return 'Invalid UUID format';
    }

    return undefined;
  };

  return (
    <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
      <h2 className={styles.sectionHeader}>Shipment Type</h2>
      <DropdownInput
        label="Shipment type"
        name="shipmentType"
        options={Object.values(SHIPMENT_OPTIONS).map((value) => ({ key: value, value }))}
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
          <AddressFields
            name="ppmShipment.pickupAddress"
            legend="Pickup Address"
            render={(fields) => (
              <>
                <p>What address are the movers picking up from?</p>
                {fields}
                <h4>Second pickup location</h4>
                <FormGroup>
                  <p>
                    Will the movers pick up any belongings from a second address? (Must be near the pickup address.
                    Subject to approval.)
                  </p>
                  <div className={formStyles.radioGroup}>
                    <Field
                      as={Radio}
                      id="has-secondary-pickup"
                      data-testid="has-secondary-pickup"
                      label="Yes"
                      name="ppmShipment.hasSecondaryPickupAddress"
                      value="true"
                      title="Yes, there is a second pickup location"
                      checked={hasSecondaryPickupAddress === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="no-secondary-pickup"
                      data-testid="no-secondary-pickup"
                      label="No"
                      name="ppmShipment.hasSecondaryPickupAddress"
                      value="false"
                      title="No, there is not a second pickup location"
                      checked={hasSecondaryPickupAddress !== 'true'}
                    />
                  </div>
                </FormGroup>
                {hasSecondaryPickupAddress === 'true' && <AddressFields name="ppmShipment.secondaryPickupAddress" />}
              </>
            )}
          />
          <h2 className={styles.sectionHeader}>Destination Info</h2>
          <AddressFields
            name="ppmShipment.destinationAddress"
            legend="Destination Address"
            render={(fields) => (
              <>
                {fields}
                <h4>Second destination address</h4>
                <FormGroup>
                  <p>
                    Will the movers deliver any belongings to a second address? (Must be near the destination address.
                    Subject to approval.)
                  </p>
                  <div className={formStyles.radioGroup}>
                    <Field
                      as={Radio}
                      data-testid="has-secondary-destination"
                      id="has-secondary-destination"
                      label="Yes"
                      name="ppmShipment.hasSecondaryDestinationAddress"
                      value="true"
                      title="Yes, there is a second destination location"
                      checked={hasSecondaryDestinationAddress === 'true'}
                    />
                    <Field
                      as={Radio}
                      data-testid="no-secondary-destination"
                      id="no-secondary-destination"
                      label="No"
                      name="ppmShipment.hasSecondaryDestinationAddress"
                      value="false"
                      title="No, there is not a second destination location"
                      checked={hasSecondaryDestinationAddress !== 'true'}
                    />
                  </div>
                </FormGroup>
                {hasSecondaryDestinationAddress === 'true' && (
                  <AddressFields name="ppmShipment.secondaryDestinationAddress" />
                )}
              </>
            )}
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
            <CheckboxField
              id="diversion"
              name="diversion"
              label="Diversion"
              onChange={(e) => toggleParentShipmentIdTextBox(e.target.checked)}
            />
            {isChecked && (
              <TextField
                data-testid="divertedFromShipmentIdInput"
                label="Diverted from Shipment ID"
                id="divertedFromShipmentIdInput"
                name="divertedFromShipmentId"
                labelHint="Required if diversion box is checked"
                validate={(value) => validateUUID(value)}
              />
            )}

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
