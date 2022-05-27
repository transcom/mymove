import React from 'react';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { CheckboxField, DatePickerInput, DropdownInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { dropdownInputOptions } from 'utils/formatters';

const PrimeUIShipmentCreateForm = () => {
  return (
    <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
      <h2 className={styles.sectionHeader}>Shipment Type</h2>
      <DropdownInput
        label="Shipment type"
        name="shipmentType"
        options={dropdownInputOptions(SHIPMENT_OPTIONS)}
        id="shipmentType"
      />

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
    </SectionWrapper>
  );
};

PrimeUIShipmentCreateForm.propTypes = {};

PrimeUIShipmentCreateForm.defaultProps = {};

export default PrimeUIShipmentCreateForm;
