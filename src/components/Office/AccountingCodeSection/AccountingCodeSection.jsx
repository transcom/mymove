import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Label, FormGroup, Radio } from '@trussworks/react-uswds';
import { useField } from 'formik';

import styles from './AccountingCodeSection.module.scss';

import { AccountingCodesShape } from 'types/accountingCodes';

const AccountingCodeSection = ({ label, fieldName, shipmentTypes, emptyMessage }) => {
  const [inputProps, , helperProps] = useField(fieldName);
  const [hasSetDefaultValue, setHasSetDefaultValue] = useState(false);

  const shipmentTypePairs = Object.entries(shipmentTypes).filter(([, value]) => !!value);

  const handleFormValueChange = (value) => {
    helperProps.setValue(value);
  };

  useEffect(() => {
    if (hasSetDefaultValue === false && shipmentTypePairs.length === 1) {
      helperProps.setValue(shipmentTypePairs[0][0]);
      setHasSetDefaultValue(true);
    }
  }, [hasSetDefaultValue, shipmentTypePairs, helperProps]);

  const handleClear = () => handleFormValueChange(undefined);

  if (shipmentTypePairs.length === 0) {
    return (
      <FormGroup>
        <Label>{label}</Label>
        <p className={styles.SectionDefaultText}>{emptyMessage}</p>
      </FormGroup>
    );
  }

  const fields = shipmentTypePairs.map(([key, value]) => {
    const isChecked = inputProps.value === key;
    const handleChange = () => handleFormValueChange(key);

    return (
      <Radio
        key={key}
        id={`${fieldName}-${key}`}
        label={`${value} (${key})`}
        name={fieldName}
        value={key}
        title={`${value} (${key})`}
        checked={isChecked}
        onChange={handleChange}
      />
    );
  });

  return (
    <>
      <FormGroup className={styles.Section}>
        <Label>{label}</Label>
        {fields}
      </FormGroup>

      <button
        type="button"
        onClick={handleClear}
        className={styles.SectionClear}
        data-testid={`clearSelection-${fieldName}`}
      >
        Clear selection
      </button>
    </>
  );
};

AccountingCodeSection.propTypes = {
  label: PropTypes.string.isRequired,
  fieldName: PropTypes.string.isRequired,
  emptyMessage: PropTypes.string.isRequired,
  shipmentTypes: AccountingCodesShape.isRequired,
};

export default AccountingCodeSection;
