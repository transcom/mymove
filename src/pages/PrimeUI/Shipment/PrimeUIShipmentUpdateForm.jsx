import React from 'react';
import PropTypes from 'prop-types';

import { ResidentialAddressShape } from 'types/address';
import { formatWeight } from 'utils/formatters';
import { dropdownInputOptions, formatDate } from 'shared/formatters';
import { formatAddress } from 'utils/shipmentDisplay';
import { CheckboxField, DatePickerInput, DropdownInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { shipmentDestinationTypes } from 'constants/shipments';

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
  editablePickupAddress,
  editableDestinationAddress,
  requestedPickupDate,
  estimatedWeight,
  actualWeight,
  pickupAddress,
  destinationAddress,
}) => {
  return (
    <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
      <h2 className={styles.sectionHeader}>Shipment Dates</h2>
      <h5 className={styles.sectionHeader}>Requested Pickup</h5>
      {formatDate(requestedPickupDate)}
      <DatePickerInput name="scheduledPickupDate" label="Scheduled pickup" />
      <DatePickerInput name="actualPickupDate" label="Actual pickup" />
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
      <h2 className={styles.sectionHeader}>Shipment Addresses</h2>
      <h5 className={styles.sectionHeader}>Pickup Address</h5>
      {editablePickupAddress && <AddressFields name="pickupAddress" />}
      {!editablePickupAddress && formatAddress(pickupAddress)}
      <h5 className={styles.sectionHeader}>Destination Address</h5>
      {editableDestinationAddress && <AddressFields name="destinationAddress" />}
      {!editableDestinationAddress && formatAddress(destinationAddress)}
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
  editablePickupAddress: PropTypes.bool,
  editableDestinationAddress: PropTypes.bool,
  requestedPickupDate: PropTypes.string,
  estimatedWeight: PropTypes.string,
  actualWeight: PropTypes.string,
  pickupAddress: ResidentialAddressShape,
  destinationAddress: ResidentialAddressShape,
};

PrimeUIShipmentUpdateForm.defaultProps = {
  editableWeightEstimateField: 0,
  editableWeightActualField: 0,
  editablePickupAddress: true,
  editableDestinationAddress: true,
  requestedPickupDate: '',
  estimatedWeight: '',
  actualWeight: '',
  pickupAddress: emptyAddressShape,
  destinationAddress: emptyAddressShape,
};

export default PrimeUIShipmentUpdateForm;
