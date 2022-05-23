import React from 'react';
import { PropTypes } from 'prop-types';

import styles from './OriginZIPInfo.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField, DatePickerInput } from 'components/form/fields';

const OriginZIPInfo = ({ setFieldValue, currentZip }) => {
  const setOriginZipToCurrentZip = (isChecked) => {
    setFieldValue('useResidentialAddressZIP', isChecked);
    if (isChecked) {
      setFieldValue('pickupPostalCode', currentZip);
    } else {
      setFieldValue('pickupPostalCode', '');
    }
  };

  return (
    <SectionWrapper className={styles.OriginZIPInfo}>
      <h2>Origin info</h2>
      <DatePickerInput label="Planned departure date" name="expectedDepartureDate" required />
      <div className="display-inline-block">
        <TextField label="Origin ZIP" id="pickupPostalCode" name="pickupPostalCode" maxLength={5} />
      </div>
      <CheckboxField
        id="useResidentialAddressZIP"
        name="useResidentialAddressZIP"
        label="Use current ZIP"
        onChange={(e) => setOriginZipToCurrentZip(e.target.checked)}
      />
      <div className="display-inline-block">
        <TextField
          label="Second origin ZIP (optional)"
          id="secondPickupPostalCode"
          name="secondPickupPostalCode"
          maxLength={5}
        />
      </div>
    </SectionWrapper>
  );
};

OriginZIPInfo.propTypes = {
  setFieldValue: PropTypes.func.isRequired,
  currentZip: PropTypes.string,
};

OriginZIPInfo.defaultProps = {
  currentZip: undefined,
};

export default OriginZIPInfo;
