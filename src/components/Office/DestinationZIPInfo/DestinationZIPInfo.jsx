import React from 'react';
import { PropTypes } from 'prop-types';

import styles from './DestinationZIPInfo.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField } from 'components/form/fields';

const DestinationZIPInfo = ({ setFieldValue, dutyZip }) => {
  const setDestinationZipToDutyZip = (isChecked) => {
    setFieldValue('useDutyZIP', isChecked);
    if (isChecked) {
      setFieldValue('destinationPostalCode', dutyZip);
    } else {
      setFieldValue('destinationPostalCode', '');
    }
  };

  return (
    <SectionWrapper className={styles.DestinationZIPInfo}>
      <h2>Destination info</h2>
      <div className="display-inline-block">
        <TextField label="Destination ZIP" id="destinationPostalCode" name="destinationPostalCode" maxLength={5} />
      </div>
      <CheckboxField
        id="useDutyZIP"
        name="useDutyZIP"
        label="Use ZIP for new duty location"
        onChange={(e) => setDestinationZipToDutyZip(e.target.checked)}
      />
      <div className="display-inline-block">
        <TextField
          label="Second destination ZIP (optional)"
          id="secondDestinationPostalCode"
          name="secondDestinationPostalCode"
          maxLength={5}
        />
      </div>
    </SectionWrapper>
  );
};

DestinationZIPInfo.propTypes = {
  setFieldValue: PropTypes.func.isRequired,
  dutyZip: PropTypes.string,
};

DestinationZIPInfo.defaultProps = {
  dutyZip: undefined,
};

export default DestinationZIPInfo;
