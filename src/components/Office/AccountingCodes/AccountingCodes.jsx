import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label, FormGroup, Radio, Button } from '@trussworks/react-uswds';
import { useField } from 'formik';

import styles from './AccountingCodes.module.scss';

import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';

const ShipmentTypeShape = PropTypes.shape({
  hhg: PropTypes.string,
  nts: PropTypes.string,
});

const AccountingCodeSection = ({ label, fieldName, shipmentTypes }) => {
  const [inputProps, , helperProps] = useField(fieldName);
  const handleClear = () => helperProps.setValue(undefined);

  const shipmentTypePairs = Object.entries(shipmentTypes);

  if (shipmentTypePairs.length === 0) {
    return (
      <FormGroup>
        <Label>{label}</Label>
        <p className={styles.SectionDefaultText}>No {fieldName.toUpperCase()} code entered.</p>
      </FormGroup>
    );
  }

  const fields = shipmentTypePairs.map(([key, value]) => {
    const isChecked = inputProps.value === value || shipmentTypePairs.length === 1;
    const handleChange = () => helperProps.setValue(value);

    return (
      <Radio
        key={key}
        id={`${fieldName}-${key}`}
        label={`${value} (${key.toUpperCase()})`}
        name={fieldName}
        value={value}
        title={`${value} (${key.toUpperCase()})`}
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

      <button type="button" onClick={handleClear} className={styles.SectionClear}>
        Clear selection
      </button>
    </>
  );
};

AccountingCodeSection.propTypes = {
  label: PropTypes.string.isRequired,
  fieldName: PropTypes.string.isRequired,
  shipmentTypes: ShipmentTypeShape.isRequired,
};

const AccountingCodes = ({ optional, TACs, SACs, onEditCodesClick }) => {
  const hasCodes = Object.keys(TACs).length + Object.keys(SACs).length > 0;

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset>
        <h2>
          Accounting codes
          {optional && (
            <span className="float-right">
              <span className={formStyles.optional}>Optional</span>
            </span>
          )}
        </h2>

        <AccountingCodeSection label="TAC" fieldName="tac" shipmentTypes={TACs} />
        <AccountingCodeSection label="SAC (optional)" fieldName="sac" shipmentTypes={SACs} />

        <Button onClick={onEditCodesClick} secondary={hasCodes}>
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
  TACs: ShipmentTypeShape,
  SACs: ShipmentTypeShape,
  onEditCodesClick: PropTypes.func,
};

export default AccountingCodes;
