import React, { useState } from 'react';
import { func, bool, string } from 'prop-types';

import styles from './OriginInfo.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField, DatePickerInput } from 'components/form/fields';

const OriginInfo = ({ setFieldValue, currentZip, isUseResidentialAddressZIPChecked, postalCodeValidator }) => {
  const [postalCodeValid, setPostalCodeValid] = useState({});

  const setOriginZipToCurrentZip = (isChecked) => {
    setFieldValue('useResidentialAddressZIP', isChecked);
    if (isChecked) {
      setFieldValue('originPostalCode', currentZip);
    } else {
      setFieldValue('originPostalCode', '');
    }
  };

  const setOriginZip = (value) => {
    if (isUseResidentialAddressZIPChecked) {
      setFieldValue('useResidentialAddressZIP', false);
    }
    setFieldValue('originPostalCode', value);
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
    <SectionWrapper className={styles.OriginInfo}>
      <h2>Origin info</h2>
      <DatePickerInput label="Planned departure date" name="plannedDepartureDate" required />
      <div className="display-inline-block">
        <TextField
          label="Origin ZIP"
          id="originPostalCode"
          name="originPostalCode"
          maxLength={5}
          onChange={(e) => {
            setOriginZip(e.target.value);
          }}
          validate={(value) => postalCodeValidate(value, 'origin', 'originPostalCode')}
        />
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
          id="secondOriginPostalCode"
          name="secondOriginPostalCode"
          maxLength={5}
          validate={(value) => postalCodeValidate(value, 'origin', 'secondOriginPostalCode')}
        />
      </div>
    </SectionWrapper>
  );
};

OriginInfo.propTypes = {
  setFieldValue: func.isRequired,
  currentZip: string,
  isUseResidentialAddressZIPChecked: bool.isRequired,
  postalCodeValidator: func.isRequired,
};

OriginInfo.defaultProps = {
  currentZip: undefined,
};

export default OriginInfo;
