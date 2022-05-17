import React, { useState } from 'react';
import { func, bool, string } from 'prop-types';

import styles from './DestinationInfo.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField } from 'components/form/fields';

const DestinationInfo = ({ setFieldValue, dutyZip, isUseDutyZIPChecked, postalCodeValidator }) => {
  const [postalCodeValid, setPostalCodeValid] = useState({});

  const setDestinationZipToDutyZip = (isChecked) => {
    setFieldValue('useDutyZIP', isChecked);
    if (isChecked) {
      setFieldValue('destinationPostalCode', dutyZip);
    } else {
      setFieldValue('destinationPostalCode', '');
    }
  };

  const setDestinationZip = (value) => {
    if (isUseDutyZIPChecked) {
      setFieldValue('useDutyZIP', false);
    }
    setFieldValue('destinationPostalCode', value);
  };

  const postalCodeValidate = async (value, location, name) => {
    if (value?.length !== 5) {
      return undefined;
    }
    if (postalCodeValid[`${name}`]?.value !== value) {
      const response = await postalCodeValidator(value, location, 'Please enter a valid ZIP code');
      setPostalCodeValid((state) => {
        return {
          ...state,
          [name]: { value, isValid: !response },
        };
      });
      return response;
    }
    return postalCodeValid[`${name}`]?.isValid ? undefined : 'Please enter a valid ZIP code';
  };

  return (
    <SectionWrapper className={styles.DestinationInfo}>
      <h2>Destination info</h2>
      <div className="display-inline-block">
        <TextField
          label="Destination ZIP"
          id="destinationPostalCode"
          name="destinationPostalCode"
          maxLength={5}
          onChange={(e) => {
            setDestinationZip(e.target.value);
          }}
          validate={(value) => postalCodeValidate(value, 'destination', 'destinationPostalCode')}
        />
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
          validate={(value) => postalCodeValidate(value, 'destination', 'secondDestinationPostalCode')}
        />
      </div>
    </SectionWrapper>
  );
};

DestinationInfo.propTypes = {
  setFieldValue: func.isRequired,
  dutyZip: string,
  isUseDutyZIPChecked: bool.isRequired,
  postalCodeValidator: func.isRequired,
};

DestinationInfo.defaultProps = {
  dutyZip: undefined,
};

export default DestinationInfo;
