import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label, FormGroup, Radio, Button, Grid } from '@trussworks/react-uswds';
import { useField } from 'formik';

import styles from './AccountingCodes.module.scss';

import formStyles from 'styles/form.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
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

const AccountingCodes = ({ optional, TACs, SACs, onEditCodesClick }) => {
  const hasCodes = Object.keys(TACs).length + Object.keys(SACs).length > 0;

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset>
        <h2 className={shipmentFormStyles.SectionHeader}>
          Accounting codes
          {optional && (
            <span className="float-right">
              <span className={formStyles.optional}>Optional</span>
            </span>
          )}
        </h2>

        <AccountingCodeSection
          label="TAC"
          emptyMessage="No TAC code entered."
          fieldName="tacType"
          shipmentTypes={TACs}
        />
        <AccountingCodeSection
          label="SAC (optional)"
          emptyMessage="No SAC code entered."
          fieldName="sacType"
          shipmentTypes={SACs}
        />

        <Button type="button" onClick={onEditCodesClick} secondary={hasCodes}>
          {hasCodes ? 'Add or edit codes' : 'Add code'}
        </Button>
      </Fieldset>
    </SectionWrapper>
  );
};

AccountingCodes.defaultProps = {
  optional: true,
  TACs: {},
  SACs: {},
  onEditCodesClick: () => {},
};

AccountingCodes.propTypes = {
  optional: PropTypes.bool,
  TACs: AccountingCodesShape,
  SACs: AccountingCodesShape,
  onEditCodesClick: PropTypes.func,
};

export default AccountingCodes;
