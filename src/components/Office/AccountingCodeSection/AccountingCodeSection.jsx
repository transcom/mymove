import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Label, FormGroup, Radio, Grid } from '@trussworks/react-uswds';
import { useField } from 'formik';

import styles from './AccountingCodeSection.module.scss';

import { AccountingCodesShape } from 'types/accountingCodes';

const AccountingCodeSection = ({ label, fieldName, shipmentTypes, emptyMessage }) => {
  const [inputProps, , helperProps] = useField(fieldName);
  const [hasSetDefaultValue, setHasSetDefaultValue] = useState(false);

  const shipmentTypePairs = Object.entries(shipmentTypes).filter(([, value]) => !!value);

  useEffect(() => {
    if (hasSetDefaultValue === false && shipmentTypePairs.length === 1) {
      helperProps.setValue(shipmentTypePairs[0][0]);
      setHasSetDefaultValue(true);
    }
  }, [hasSetDefaultValue, shipmentTypePairs, helperProps]);

  const handleClear = () => helperProps.setValue(undefined);

  if (shipmentTypePairs.length === 0) {
    return (
      <Grid row>
        <Grid col={12}>
          <FormGroup>
            <Label>{label}</Label>
            <p className={styles.SectionDefaultText}>{emptyMessage}</p>
          </FormGroup>
        </Grid>
      </Grid>
    );
  }

  const fields = shipmentTypePairs.map(([key, value]) => {
    const isChecked = inputProps.value === key;
    const handleChange = () => helperProps.setValue(key);

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
    <Grid row>
      <Grid col={12}>
        <FormGroup className={styles.Section}>
          <Label>{label}</Label>
          {fields}
        </FormGroup>

        <button type="button" onClick={handleClear} className={styles.SectionClear}>
          Clear selection
        </button>
      </Grid>
    </Grid>
  );
};

AccountingCodeSection.propTypes = {
  label: PropTypes.string.isRequired,
  fieldName: PropTypes.string.isRequired,
  emptyMessage: PropTypes.string.isRequired,
  shipmentTypes: AccountingCodesShape.isRequired,
};

export default AccountingCodeSection;
