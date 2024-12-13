import React from 'react';
import PropTypes from 'prop-types';

import { ResidentialAddressShape } from 'types/address';
import { formatWeight, formatDate, dropdownInputOptions } from 'utils/formatters';
import { formatAddress } from 'utils/shipmentDisplay';
import { CheckboxField, DatePickerInput, DropdownInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { shipmentDestinationTypes } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const emptyAddressShape = {
  streetAddress1: '',
  streetAddress2: '',
  city: '',
  state: '',
  postalCode: '',
};

const shipmentDestinationTypeOptions = dropdownInputOptions(shipmentDestinationTypes);

const PrimeUIShipmentUpdateForm = ({
  editableWeightEstimateField,
  editableWeightActualField,
  editableProGearWeightActualField,
  editableSpouseProGearWeightActualField,
  editablePickupAddress,
  editableSecondaryPickupAddress,
  editableTertiaryPickupAddress,
  editableDestinationAddress,
  editableSecondaryDeliveryAddress,
  editableTertiaryDeliveryAddress,
  requestedPickupDate,
  estimatedWeight,
  actualWeight,
  actualProGearWeight,
  actualSpouseProGearWeight,
  pickupAddress,
  secondaryPickupAddress,
  tertiaryPickupAddress,
  destinationAddress,
  secondaryDeliveryAddress,
  tertiaryDeliveryAddress,
  shipmentType,
  formikProps,
}) => {
  const isNTS = shipmentType === SHIPMENT_OPTIONS.NTS;
  const isNTSR = shipmentType === SHIPMENT_OPTIONS.NTSR;
  return (
    <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
      <h2 className={styles.sectionHeader}>Shipment Dates</h2>
      <h5 className={styles.sectionHeader}>Requested Pickup</h5>
      {formatDate(requestedPickupDate)}
      <DatePickerInput name="scheduledPickupDate" label="Scheduled pickup" />
      <DatePickerInput name="actualPickupDate" label="Actual pickup" />
      <DatePickerInput name="scheduledDeliveryDate" label="Scheduled delivery" />
      <DatePickerInput name="actualDeliveryDate" label="Actual delivery" />
      <h2 className={styles.sectionHeader}>Diversion</h2>
      <CheckboxField id="diversion" name="diversion" label="Diversion" />
      <h2 className={styles.sectionHeader}>Shipment Weights</h2>
      {editableWeightEstimateField && (
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
      )}
      {!editableWeightEstimateField && (
        <>
          <dt>
            <h5 className={styles.sectionHeader}>Estimated Weight</h5>
          </dt>
          <dd data-testid="authorizedWeight">{formatWeight(estimatedWeight)}</dd>
        </>
      )}
      {editableWeightActualField && (
        <MaskedTextField
          data-testid="actualWeightInput"
          defaultValue="0"
          name="actualWeight"
          label="Actual weight (lbs)"
          id="actualWeightInput"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          signed={false} // disallow negative
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
      )}
      {!editableWeightActualField && (
        <>
          <dt>
            <h5 className={styles.sectionHeader}>Actual Weight</h5>
          </dt>
          <dd data-testid="authorizedWeight">{formatWeight(actualWeight)}</dd>
        </>
      )}
      <h2 className={styles.sectionHeader}>Pro Gear Weights</h2>
      {editableProGearWeightActualField && (
        <MaskedTextField
          data-testid="actualProGearWeightInput"
          defaultValue="0"
          name="actualProGearWeight"
          label="Actual pro gear weight (lbs)"
          id="actualProGearWeightInput"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          signed={false} // disallow negative
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
      )}
      {!editableProGearWeightActualField && (
        <>
          <dt>
            <h5 className={styles.sectionHeader}>Actual Pro Gear Weight</h5>
          </dt>
          <dd data-testid="authorizedWeight">{formatWeight(actualProGearWeight)}</dd>
        </>
      )}
      {editableSpouseProGearWeightActualField && (
        <MaskedTextField
          data-testid="actualSpouseProGearWeightInput"
          defaultValue="0"
          name="actualSpouseProGearWeight"
          label="Actual spouse pro gear weight (lbs)"
          id="actualSpouseProGearWeightInput"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          signed={false} // disallow negative
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
      )}
      {!editableSpouseProGearWeightActualField && (
        <>
          <dt>
            <h5 className={styles.sectionHeader}>Actual Pro Gear Weight</h5>
          </dt>
          <dd data-testid="authorizedWeight">{formatWeight(actualSpouseProGearWeight)}</dd>
        </>
      )}
      <h2 className={styles.sectionHeader}>Shipment Addresses</h2>
      <h5 className={styles.sectionHeader}>Pickup Address</h5>
      {editablePickupAddress && <AddressFields name="pickupAddress" locationLookup formikProps={formikProps} />}
      {!editablePickupAddress && formatAddress(pickupAddress)}
      {!isNTSR && (
        <>
          <h5 className={styles.sectionHeader}>Second Pickup Address</h5>
          {editableSecondaryPickupAddress && (
            <AddressFields name="secondaryPickupAddress" locationLookup formikProps={formikProps} />
          )}
          {!editableSecondaryPickupAddress && formatAddress(secondaryPickupAddress)}
          <h5 className={styles.sectionHeader}>Third Pickup Address</h5>
          {editableTertiaryPickupAddress && (
            <AddressFields name="tertiaryPickupAddress" locationLookup formikProps={formikProps} />
          )}
          {!editableTertiaryPickupAddress && formatAddress(tertiaryPickupAddress)}
        </>
      )}
      <h5 className={styles.sectionHeader}>Delivery Address</h5>
      {editableDestinationAddress && (
        <AddressFields name="destinationAddress" locationLookup formikProps={formikProps} />
      )}
      {!editableDestinationAddress && formatAddress(destinationAddress)}
      {!isNTS && (
        <>
          <h5 className={styles.sectionHeader}>Second Delivery Address</h5>
          {editableSecondaryDeliveryAddress && (
            <AddressFields name="secondaryDeliveryAddress" locationLookup formikProps={formikProps} />
          )}
          {!editableSecondaryDeliveryAddress && formatAddress(secondaryDeliveryAddress)}
          <h5 className={styles.sectionHeader}>Third Delivery Address</h5>
          {editableTertiaryDeliveryAddress && (
            <AddressFields name="tertiaryDeliveryAddress" locationLookup formikProps={formikProps} />
          )}
          {!editableTertiaryDeliveryAddress && formatAddress(tertiaryDeliveryAddress)}
        </>
      )}
      <DropdownInput
        label="Destination type"
        name="destinationType"
        options={shipmentDestinationTypeOptions}
        id="destinationType"
      />
    </SectionWrapper>
  );
};

PrimeUIShipmentUpdateForm.propTypes = {
  editableWeightEstimateField: PropTypes.bool,
  editableWeightActualField: PropTypes.bool,
  editableProGearWeightActualField: PropTypes.bool,
  editableSpouseProGearWeightActualField: PropTypes.bool,
  editablePickupAddress: PropTypes.bool,
  editableDestinationAddress: PropTypes.bool,
  requestedPickupDate: PropTypes.string,
  estimatedWeight: PropTypes.string,
  actualWeight: PropTypes.string,
  actualProGearWeight: PropTypes.string,
  actualSpouseProGearWeight: PropTypes.string,
  pickupAddress: ResidentialAddressShape,
  destinationAddress: ResidentialAddressShape,
  formikProps: PropTypes.object.isRequired,
};

PrimeUIShipmentUpdateForm.defaultProps = {
  editableWeightEstimateField: 0,
  editableWeightActualField: 0,
  editableProGearWeightActualField: true,
  editableSpouseProGearWeightActualField: true,
  editablePickupAddress: true,
  editableDestinationAddress: true,
  requestedPickupDate: '',
  estimatedWeight: '',
  actualWeight: '',
  actualProGearWeight: '',
  actualSpouseProGearWeight: '',
  pickupAddress: emptyAddressShape,
  destinationAddress: emptyAddressShape,
};

export default PrimeUIShipmentUpdateForm;
