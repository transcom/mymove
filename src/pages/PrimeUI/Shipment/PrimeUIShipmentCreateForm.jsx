import React, { useState, useEffect } from 'react';
import { Radio, FormGroup, Label, Textarea } from '@trussworks/react-uswds';
import { Field, useField, useFormikContext } from 'formik';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';

import { SHIPMENT_OPTIONS, SHIPMENT_TYPES, FEATURE_FLAG_KEYS } from 'shared/constants';
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
  const { values, setFieldTouched, setFieldValue } = useFormikContext();
  const { shipmentType } = values;
  const { sitExpected, hasProGear } = values.ppmShipment;
  const { hasTrailer } = values.boatShipment;
  const [, , checkBoxHelperProps] = useField('diversion');
  const [, , divertedFromIdHelperProps] = useField('divertedFromShipmentId');
  const [isChecked, setIsChecked] = useState(false);
  const [enableBoat, setEnableBoat] = useState(false);
  const [enableMobileHome, setEnableMobileHome] = useState(false);

  const hasShipmentType = !!shipmentType;
  const isPPM = shipmentType === SHIPMENT_OPTIONS.PPM;
  const isBoat = shipmentType === SHIPMENT_TYPES.BOAT_HAUL_AWAY || shipmentType === SHIPMENT_TYPES.BOAT_TOW_AWAY;
  const isMobileHome = shipmentType === SHIPMENT_TYPES.MOBILE_HOME;

  let {
    hasSecondaryPickupAddress,
    hasSecondaryDestinationAddress,
    hasTertiaryPickupAddress,
    hasTertiaryDestinationAddress,
  } = '';

  if (isPPM) {
    ({
      hasSecondaryPickupAddress,
      hasSecondaryDestinationAddress,
      hasTertiaryPickupAddress,
      hasTertiaryDestinationAddress,
    } = values.ppmShipment);
  } else {
    ({
      hasSecondaryPickupAddress,
      hasSecondaryDestinationAddress,
      hasTertiaryPickupAddress,
      hasTertiaryDestinationAddress,
    } = values);
  }

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

  useEffect(() => {
    const fetchData = async () => {
      setEnableBoat(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.BOAT));
      setEnableMobileHome(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.MOBILE_HOME));
    };
    fetchData();
  }, []);

  let shipmentTypeOptions = Object.values(SHIPMENT_TYPES).map((value) => ({ key: value, value }));
  if (!enableBoat) {
    // Disallow the Prime from choosing Boat shipments if the feature flag is not enabled
    shipmentTypeOptions = shipmentTypeOptions.filter(
      (e) => e.key !== SHIPMENT_TYPES.BOAT_HAUL_AWAY && e.key !== SHIPMENT_TYPES.BOAT_TOW_AWAY,
    );
  }
  if (!enableMobileHome) {
    // Disallow the Prime from choosing Mobile Home shipments if the feature flag is not enabled
    shipmentTypeOptions = shipmentTypeOptions.filter((e) => e.key !== SHIPMENT_TYPES.MOBILE_HOME);
  }

  return (
    <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
      <h2 className={styles.sectionHeader}>Shipment Type</h2>
      <DropdownInput label="Shipment type" name="shipmentType" options={shipmentTypeOptions} id="shipmentType" />

      {isPPM && (
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
            locationLookup
            formikProps={{
              setFieldTouched,
              setFieldValue,
            }}
            render={(fields) => (
              <>
                <p>What address are the movers picking up from?</p>
                {fields}
                <h4>Second Pickup Address</h4>
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
                      title="Yes, there is a second pickup address"
                      checked={hasSecondaryPickupAddress === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="no-secondary-pickup"
                      data-testid="no-secondary-pickup"
                      label="No"
                      name="ppmShipment.hasSecondaryPickupAddress"
                      value="false"
                      title="No, there is not a second pickup address"
                      checked={hasSecondaryPickupAddress !== 'true' && hasTertiaryPickupAddress !== 'true'}
                    />
                  </div>
                </FormGroup>
                {hasSecondaryPickupAddress === 'true' && (
                  <>
                    <h5 className={styles.sectionHeader}>Second Pickup Address</h5>
                    <AddressFields
                      name="ppmShipment.secondaryPickupAddress"
                      locationLookup
                      formikProps={{
                        setFieldTouched,
                        setFieldValue,
                      }}
                    />

                    <h4>Third Pickup Address</h4>
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
                          title="Yes, there is a tertiary pickup address"
                          checked={hasTertiaryPickupAddress === 'true'}
                        />
                        <Field
                          as={Radio}
                          id="no-tertiary-pickup"
                          data-testid="no-tertiary-pickup"
                          label="No"
                          name="ppmShipment.hasTertiaryPickupAddress"
                          value="false"
                          title="No, there is not a tertiary pickup address"
                          checked={hasTertiaryPickupAddress !== 'true'}
                        />
                      </div>
                    </FormGroup>
                  </>
                )}
                {hasTertiaryPickupAddress === 'true' && hasSecondaryPickupAddress === 'true' && (
                  <>
                    <h5 className={styles.sectionHeader}>Third Pickup Address</h5>
                    <AddressFields
                      name="ppmShipment.tertiaryPickupAddress"
                      locationLookup
                      formikProps={{
                        setFieldTouched,
                        setFieldValue,
                      }}
                    />
                  </>
                )}
              </>
            )}
          />
          <h2 className={styles.sectionHeader}>Destination Info</h2>
          <AddressFields
            name="ppmShipment.destinationAddress"
            legend="Delivery Address"
            address1LabelHint="Optional"
            locationLookup
            formikProps={{
              setFieldTouched,
              setFieldValue,
            }}
            render={(fields) => (
              <>
                {fields}
                <h4>Second Delivery Address</h4>
                <FormGroup>
                  <p>
                    Will the movers deliver any belongings to a second address? (Must be near the delivery address.
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
                      checked={hasSecondaryDestinationAddress !== 'true' && hasTertiaryDestinationAddress !== 'true'}
                    />
                  </div>
                </FormGroup>
                {hasSecondaryDestinationAddress === 'true' && (
                  <>
                    <h5 className={styles.sectionHeader}>Second Delivery Address</h5>
                    <AddressFields
                      name="ppmShipment.secondaryDestinationAddress"
                      locationLookup
                      formikProps={{
                        setFieldTouched,
                        setFieldValue,
                      }}
                    />

                    <h4>Third Delivery Address</h4>
                    <FormGroup>
                      <p>
                        Will the movers pick up any belongings from a third address? (Must be near the pickup address.
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
                          title="Yes, there is a third delivery address"
                          checked={hasTertiaryDestinationAddress === 'true'}
                        />
                        <Field
                          as={Radio}
                          id="no-tertiary-destination"
                          data-testid="no-tertiary-destination"
                          label="No"
                          name="ppmShipment.hasTertiaryDestinationAddress"
                          value="false"
                          title="No, there is not a third delivery address"
                          checked={hasTertiaryDestinationAddress !== 'true'}
                        />
                      </div>
                    </FormGroup>
                  </>
                )}
                {hasTertiaryDestinationAddress === 'true' && hasSecondaryDestinationAddress === 'true' && (
                  <>
                    <h5 className={styles.sectionHeader}>Third Delivery Address</h5>
                    <AddressFields
                      name="ppmShipment.tertiaryDestinationAddress"
                      locationLookup
                      formikProps={{
                        setFieldTouched,
                        setFieldValue,
                      }}
                    />
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
        </>
      )}
      {hasShipmentType && !isPPM && (
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
          <AddressFields
            name="pickupAddress"
            locationLookup
            formikProps={{
              setFieldTouched,
              setFieldValue,
            }}
            render={(fields) => (
              <>
                {fields}
                <h4>Second Pickup Address</h4>
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
                      name="hasSecondaryPickupAddress"
                      value="true"
                      title="Yes, there is a second pickup address"
                      checked={hasSecondaryPickupAddress === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="no-secondary-pickup"
                      data-testid="no-secondary-pickup"
                      label="No"
                      name="hasSecondaryPickupAddress"
                      value="false"
                      title="No, there is not a second pickup address"
                      checked={hasSecondaryPickupAddress !== 'true'}
                    />
                  </div>
                </FormGroup>
                {hasSecondaryPickupAddress === 'true' && (
                  <>
                    <h5 className={styles.sectionHeader}>Second Pickup Address</h5>
                    <AddressFields
                      name="secondaryPickupAddress"
                      locationLookup
                      formikProps={{
                        setFieldTouched,
                        setFieldValue,
                      }}
                    />

                    <h4>Third Pickup Address</h4>
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
                          name="hasTertiaryPickupAddress"
                          value="true"
                          title="Yes, there is a tertiary pickup address"
                          checked={hasTertiaryPickupAddress === 'true'}
                        />
                        <Field
                          as={Radio}
                          id="no-tertiary-pickup"
                          data-testid="no-tertiary-pickup"
                          label="No"
                          name="hasTertiaryPickupAddress"
                          value="false"
                          title="No, there is not a tertiary pickup address"
                          checked={hasTertiaryPickupAddress !== 'true'}
                        />
                      </div>
                    </FormGroup>
                  </>
                )}
                {hasTertiaryPickupAddress === 'true' && hasSecondaryPickupAddress === 'true' && (
                  <>
                    <h5 className={styles.sectionHeader}>Third Pickup Address</h5>
                    <AddressFields
                      name="tertiaryPickupAddress"
                      locationLookup
                      formikProps={{
                        setFieldTouched,
                        setFieldValue,
                      }}
                    />
                  </>
                )}
              </>
            )}
          />

          <h3 className={styles.sectionHeader}>Destination Info</h3>
          <AddressFields
            name="destinationAddress"
            legend="Delivery Address"
            locationLookup
            formikProps={{
              setFieldTouched,
              setFieldValue,
            }}
            render={(fields) => (
              <>
                {fields}

                <h4>Second Delivery Address</h4>
                <FormGroup>
                  <p>
                    Will the movers pick up any belongings from a second address? (Must be near the pickup address.
                    Subject to approval.)
                  </p>
                  <div className={formStyles.radioGroup}>
                    <Field
                      as={Radio}
                      id="has-secondary-destination"
                      data-testid="has-secondary-destination"
                      label="Yes"
                      name="hasSecondaryDestinationAddress"
                      value="true"
                      title="Yes, there is a second delivery address"
                      checked={hasSecondaryDestinationAddress === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="no-secondary-destination"
                      data-testid="no-secondary-destination"
                      label="No"
                      name="hasSecondaryDestinationAddress"
                      value="false"
                      title="No, there is not a second delivery address"
                      checked={hasSecondaryDestinationAddress !== 'true'}
                    />
                  </div>
                </FormGroup>
                {hasSecondaryDestinationAddress === 'true' && (
                  <>
                    <h5 className={styles.sectionHeader}>Second Delivery Address</h5>
                    <AddressFields
                      name="secondaryDestinationAddress"
                      locationLookup
                      formikProps={{
                        setFieldTouched,
                        setFieldValue,
                      }}
                    />

                    <h4>Third Delivery Address</h4>
                    <FormGroup>
                      <p>
                        Will the movers pick up any belongings from a third address? (Must be near the pickup address.
                        Subject to approval.)
                      </p>
                      <div className={formStyles.radioGroup}>
                        <Field
                          as={Radio}
                          id="has-tertiary-destination"
                          data-testid="has-tertiary-destination"
                          label="Yes"
                          name="hasTertiaryDestinationAddress"
                          value="true"
                          title="Yes, there is a third delivery address"
                          checked={hasTertiaryDestinationAddress === 'true'}
                        />
                        <Field
                          as={Radio}
                          id="no-tertiary-destination"
                          data-testid="no-tertiary-destination"
                          label="No"
                          name="hasTertiaryDestinationAddress"
                          value="false"
                          title="No, there is not a third delivery address"
                          checked={hasTertiaryDestinationAddress !== 'true'}
                        />
                      </div>
                    </FormGroup>
                  </>
                )}
                {hasTertiaryDestinationAddress === 'true' && hasSecondaryDestinationAddress === 'true' && (
                  <>
                    <h5 className={styles.sectionHeader}>Third Delivery Address</h5>
                    <AddressFields
                      name="tertiaryDestinationAddress"
                      locationLookup
                      formikProps={{
                        setFieldTouched,
                        setFieldValue,
                      }}
                    />
                  </>
                )}
              </>
            )}
          />
        </>
      )}
      {isBoat && (
        <>
          <h2 className={styles.sectionHeader}>Boat Model Info</h2>
          <MaskedTextField
            label="Year"
            id="boatShipment.yearInput"
            name="boatShipment.year"
            mask={Number}
            maxLength={4}
          />
          <TextField label="Make" id="boatShipment.makeInput" name="boatShipment.make" />
          <TextField label="Model" id="boatShipment.modelInput" name="boatShipment.model" />
          <h2 className={styles.sectionHeader}>Boat Dimensions</h2>
          <figure>
            <figcaption>
              Dimensions must meet at least one of the following criteria to qualify as a separate boat shipment:
            </figcaption>
            <ul>
              <li>Over 14 feet in length</li>
              <li>Over 6 feet 10 inches in width</li>
              <li>Over 6 feet 5 inches in height</li>
            </ul>
          </figure>
          <MaskedTextField
            label="Length (Feet)"
            id="boatShipment.lengthInFeetInput"
            name="boatShipment.lengthInFeet"
            mask={Number}
            min={Number.MIN_SAFE_INTEGER}
            max={Number.MAX_SAFE_INTEGER}
          />
          <MaskedTextField
            label="Length (Inches)"
            id="boatShipment.lengthInInchesInput"
            name="boatShipment.lengthInInches"
            mask={Number}
          />
          <MaskedTextField
            label="Width (Feet)"
            id="boatShipment.widthInFeetInput"
            name="boatShipment.widthInFeet"
            mask={Number}
            min={Number.MIN_SAFE_INTEGER}
            max={Number.MAX_SAFE_INTEGER}
          />
          <MaskedTextField
            label="Width (Inches)"
            id="boatShipment.widthInInchesInput"
            name="boatShipment.widthInInches"
            mask={Number}
          />
          <MaskedTextField
            label="Height (Feet)"
            id="boatShipment.heightInFeetInput"
            name="boatShipment.heightInFeet"
            mask={Number}
            min={Number.MIN_SAFE_INTEGER}
            max={Number.MAX_SAFE_INTEGER}
          />
          <MaskedTextField
            label="Height (Inches)"
            id="boatShipment.heightInInchesInput"
            name="boatShipment.heightInInches"
            mask={Number}
          />
          <h2 className={styles.sectionHeader}>Trailer</h2>
          <CheckboxField label="Has Trailer" id="boatShipment.hasTrailerInput" name="boatShipment.hasTrailer" />
          {hasTrailer && (
            <CheckboxField
              label="Trailer is Roadworthy"
              id="boatShipment.isRoadworthyInput"
              name="boatShipment.isRoadworthy"
            />
          )}
          <h2 className={styles.sectionHeader}>Remarks</h2>
          <Label htmlFor="counselorRemarksInput">Counselor Remarks</Label>
          <Field id="counselorRemarksInput" name="counselorRemarks" as={Textarea} className={`${formStyles.remarks}`} />
        </>
      )}
      {isMobileHome && (
        <>
          <h2 className={styles.sectionHeader}>Mobile Home Model Info</h2>
          <MaskedTextField
            label="Year"
            id="mobileHomeShipment.yearInput"
            name="mobileHomeShipment.year"
            mask={Number}
            maxLength={4}
          />
          <TextField label="Make" id="mobileHomeShipment.makeInput" name="mobileHomeShipment.make" />
          <TextField label="Model" id="mobileHomeShipment.modelInput" name="mobileHomeShipment.model" />
          <h2 className={styles.sectionHeader}>Mobile Home Dimensions</h2>
          <MaskedTextField
            label="Length (Feet)"
            id="mobileHomeShipment.lengthInFeetInput"
            name="mobileHomeShipment.lengthInFeet"
            mask={Number}
            min={Number.MIN_SAFE_INTEGER}
            max={Number.MAX_SAFE_INTEGER}
          />
          <MaskedTextField
            label="Length (Inches)"
            id="mobileHomeShipment.lengthInInchesInput"
            name="mobileHomeShipment.lengthInInches"
            mask={Number}
          />
          <MaskedTextField
            label="Width (Feet)"
            id="mobileHomeShipment.widthInFeetInput"
            name="mobileHomeShipment.widthInFeet"
            mask={Number}
            min={Number.MIN_SAFE_INTEGER}
            max={Number.MAX_SAFE_INTEGER}
          />
          <MaskedTextField
            label="Width (Inches)"
            id="mobileHomeShipment.widthInInchesInput"
            name="mobileHomeShipment.widthInInches"
            mask={Number}
          />
          <MaskedTextField
            label="Height (Feet)"
            id="mobileHomeShipment.heightInFeetInput"
            name="mobileHomeShipment.heightInFeet"
            mask={Number}
            min={Number.MIN_SAFE_INTEGER}
            max={Number.MAX_SAFE_INTEGER}
          />
          <MaskedTextField
            label="Height (Inches)"
            id="heightInches"
            name="mobileHomeShipment.heightInInches"
            mask={Number}
            max={11}
          />
          <h2 className={styles.sectionHeader}>Remarks</h2>
          <Label htmlFor="counselorRemarksInput">Counselor Remarks</Label>
          <Field id="counselorRemarksInput" name="counselorRemarks" as={Textarea} className={`${formStyles.remarks}`} />
        </>
      )}
    </SectionWrapper>
  );
};

PrimeUIShipmentCreateForm.propTypes = {};

PrimeUIShipmentCreateForm.defaultProps = {};

export default PrimeUIShipmentCreateForm;
