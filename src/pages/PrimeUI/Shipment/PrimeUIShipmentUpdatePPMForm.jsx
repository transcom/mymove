import React from 'react';
import { Radio, FormGroup, Label, Textarea } from '@trussworks/react-uswds';
import { Field, useFormikContext } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { CheckboxField, DatePickerInput, DropdownInput } from 'components/form/fields';
import { dropdownInputOptions } from 'utils/formatters';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const sitLocationOptions = dropdownInputOptions(LOCATION_TYPES);

const PrimeUIShipmentUpdatePPMForm = () => {
  const { values } = useFormikContext();
  const {
    sitExpected,
    hasProGear,
    hasSecondaryPickupAddress,
    hasTertiaryPickupAddress,
    hasSecondaryDestinationAddress,
    hasTertiaryDestinationAddress,
  } = values.ppmShipment;

  return (
    <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
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
                Will the movers pick up any belongings from a second address? (Must be near the pickup address. Subject
                to approval.)
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
            {hasSecondaryPickupAddress === 'true' && (
              <>
                <AddressFields name="ppmShipment.secondaryPickupAddress" />
                <h4>Third pickup location</h4>
                <FormGroup>
                  <p>
                    Will the movers pick up any belongings from a third address? (Must be near the pickup address.
                    Subject to approval.)
                  </p>
                  <div className={formStyles.radioGroup}>
                    <Field
                      as={Radio}
                      id="has-tertiary-pickup"
                      data-testid="has-tertiary-pickup"
                      label="Yes"
                      name="ppmShipment.hasTertiaryPickupAddress"
                      value="true"
                      title="Yes, there is a third pickup location"
                      checked={hasTertiaryPickupAddress === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="no-tertiary-pickup"
                      data-testid="no-tertiary-pickup"
                      label="No"
                      name="ppmShipment.hasTertiaryPickupAddress"
                      value="false"
                      title="No, there is not a third pickup location"
                      checked={hasTertiaryPickupAddress !== 'true'}
                    />
                  </div>
                </FormGroup>
                {hasTertiaryPickupAddress === 'true' && <AddressFields name="ppmShipment.tertiaryPickupAddress" />}
              </>
            )}
          </>
        )}
      />
      <h2 className={styles.sectionHeader}>Destination Info</h2>
      <AddressFields
        name="ppmShipment.destinationAddress"
        legend="Delivery Address"
        render={(fields) => (
          <>
            {fields}
            <h4>Second delivery address</h4>
            <FormGroup>
              <p>
                Will the movers deliver any belongings to a second address? (Must be near the delivery address. Subject
                to approval.)
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
              <>
                <AddressFields name="ppmShipment.secondaryDestinationAddress" />
                <h4>Third destination location</h4>
                <FormGroup>
                  <p>
                    Will the movers pick up any belongings from a third address? (Must be near the Delivery Address.
                    Subject to approval.)
                  </p>
                  <div className={formStyles.radioGroup}>
                    <Field
                      as={Radio}
                      id="has-tertiary-destination"
                      data-testid="has-tertiary-destination"
                      label="Yes"
                      name="ppmShipment.hasTertiaryDestinationAddress"
                      value="true"
                      title="Yes, there is a third destination location"
                      checked={hasTertiaryDestinationAddress === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="no-tertiary-destination"
                      data-testid="no-tertiary-destination"
                      label="No"
                      name="ppmShipment.hasTertiaryDestinationAddress"
                      value="false"
                      title="No, there is not a third destination location"
                      checked={hasTertiaryDestinationAddress !== 'true'}
                    />
                  </div>
                </FormGroup>
                {hasTertiaryDestinationAddress === 'true' && (
                  <AddressFields name="ppmShipment.tertiaryDestinationAddress" />
                )}
              </>
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
    </SectionWrapper>
  );
};

export default PrimeUIShipmentUpdatePPMForm;
