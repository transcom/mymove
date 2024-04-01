import React from 'react';
import { Radio, FormGroup, Fieldset, Label, Textarea } from '@trussworks/react-uswds';
import { Field, useFormikContext } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { CheckboxField, DatePickerInput, DropdownInput } from 'components/form/fields';
import { dropdownInputOptions } from 'utils/formatters';
import Hint from 'components/Hint/index';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const sitLocationOptions = dropdownInputOptions(LOCATION_TYPES);

const PrimeUIShipmentUpdatePPMForm = () => {
  const { values } = useFormikContext();
  const { sitExpected, hasProGear } = values.ppmShipment;

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
        name="pickupAddress.address"
        render={(fields) => (
          <>
            <p>What address are the movers picking up from?</p>
            {fields}
            <h4>Second pickup location</h4>
            <FormGroup>
              <Fieldset>
                <legend className="usa-label">Will you add items to your PPM from a different address?</legend>
                <Field
                  as={Radio}
                  data-testid="yes-secondary-pickup-address"
                  id="yes-secondary-pickup-address"
                  label="Yes"
                  name="hasSecondaryPickupAddress"
                  value="true"
                  checked={values.hasSecondaryPickupAddress === 'true'}
                />
                <Field
                  as={Radio}
                  data-testid="no-secondary-pickup-address"
                  id="no-secondary-pickup-address"
                  label="No"
                  name="hasSecondaryPickupAddress"
                  value="false"
                  checked={values.hasSecondaryPickupAddress === 'false'}
                />
              </Fieldset>
            </FormGroup>
            {values.hasSecondaryPickupAddress === 'true' && (
              <>
                <AddressFields name="secondaryPickupAddress.address" />
                <Hint className={ppmStyles.hint}>
                  <p>A second origin address could mean that your final incentive is lower than your estimate.</p>
                  <p>
                    Get separate weight tickets for each leg of the trip to show how the weight changes. Talk to your
                    move counselor for more detailed information.
                  </p>
                </Hint>
              </>
            )}
          </>
        )}
      />
      <h2 className={styles.sectionHeader}>Destination Info</h2>
      <AddressFields
        name="destinationAddress.address"
        render={(fields) => (
          <>
            <p>Please input Delivery Address</p>
            {fields}
            <FormGroup>
              <Fieldset>
                <legend className="usa-label">Will you deliver part of your PPM to a different address?</legend>
                <Field
                  as={Radio}
                  data-testid="yes-secondary-destination-address"
                  id="hasSecondaryDestinationAddressYes"
                  label="Yes"
                  name="hasSecondaryDestinationAddress"
                  value="true"
                  checked={values.hasSecondaryDestinationAddress === 'true'}
                />
                <Field
                  as={Radio}
                  data-testid="no-secondary-destination-address"
                  id="hasSecondaryDestinationAddressNo"
                  label="No"
                  name="hasSecondaryDestinationAddress"
                  value="false"
                  checked={values.hasSecondaryDestinationAddress === 'false'}
                />
              </Fieldset>
            </FormGroup>
            {values.hasSecondaryDestinationAddress === 'true' && (
              <>
                <AddressFields name="secondaryDestinationAddress.address" />
                <Hint className={ppmStyles.hint}>
                  <p>A second destination ZIP could mean that your final incentive is lower than your estimate.</p>
                  <p>
                    Get separate weight tickets for each leg of the trip to show how the weight changes. Talk to your
                    move counselor for more detailed information.
                  </p>
                </Hint>
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
