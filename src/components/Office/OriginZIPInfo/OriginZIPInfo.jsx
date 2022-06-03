import React, { useState } from 'react';
import { PropTypes } from 'prop-types';
import { useField } from 'formik';

import styles from './OriginZIPInfo.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField, DatePickerInput } from 'components/form/fields';
import { UnsupportedZipCodePPMErrorMsg } from 'utils/validation';

const OriginZIPInfo = ({ currentZip, postalCodeValidator }) => {
  const [postalCodeValid, setPostalCodeValid] = useState({});
  const [isChecked, setIsChecked] = useState(false);
  const [, , postalCodeHelperProps] = useField('pickupPostalCode');
  const [, , checkBoxHelperProps] = useField('useResidentialAddressZIP');

  const setOriginZipToCurrentZip = (checkboxValue) => {
    if (checkboxValue) {
      postalCodeHelperProps.setValue(currentZip);
      checkBoxHelperProps.setValue('checked');
      setIsChecked(true);
    } else {
      postalCodeHelperProps.setValue('');
      checkBoxHelperProps.setValue('');
      setIsChecked(false);
    }
  };

  const handlePrefillPostalCodeChange = (value) => {
    if (isChecked && value !== currentZip) {
      checkBoxHelperProps.setValue('checked');
    }
    postalCodeHelperProps.setValue(value);
  };

  const postalCodeValidate = async (value, location, name) => {
    if (value?.length !== 5) {
      return undefined;
    }
    // only revalidate if the value has changed, editing other fields will re-validate unchanged ones
    if (postalCodeValid[`${name}`]?.value !== value) {
      const response = await postalCodeValidator(value, location, UnsupportedZipCodePPMErrorMsg);
      setPostalCodeValid((state) => {
        return {
          ...state,
          [name]: { value, isValid: !response },
        };
      });
      return response;
    }
    return postalCodeValid[`${name}`]?.isValid ? undefined : UnsupportedZipCodePPMErrorMsg;
  };

  return (
    <SectionWrapper className={styles.OriginZIPInfo}>
      <h2>Origin info</h2>
      <DatePickerInput label="Planned departure date" name="expectedDepartureDate" required />
      <div className="display-inline-block">
        <TextField
          label="Origin ZIP"
          id="pickupPostalCode"
          name="pickupPostalCode"
          maxLength={5}
          onChange={(e) => {
            handlePrefillPostalCodeChange(e.target.value);
          }}
          validate={(value) => postalCodeValidate(value, 'origin', 'pickupPostalCode')}
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
          label="Second origin ZIP"
          id="secondPickupPostalCode"
          name="secondPickupPostalCode"
          maxLength={5}
          optional
          validate={(value) => postalCodeValidate(value, 'origin', 'secondaryPickupPostalCode')}
        />
      </div>
    </SectionWrapper>
  );
};

OriginZIPInfo.propTypes = {
  currentZip: PropTypes.string,
  postalCodeValidator: PropTypes.func.isRequired,
};

OriginZIPInfo.defaultProps = {
  currentZip: undefined,
};

export default OriginZIPInfo;
