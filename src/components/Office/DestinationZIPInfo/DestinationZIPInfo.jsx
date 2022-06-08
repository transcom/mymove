import React, { useState } from 'react';
import { PropTypes } from 'prop-types';
import { useField } from 'formik';

import styles from './DestinationZIPInfo.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField } from 'components/form/fields';
import { UnsupportedZipCodePPMErrorMsg } from 'utils/validation';

const DestinationZIPInfo = ({ dutyZip, postalCodeValidator }) => {
  const [postalCodeValid, setPostalCodeValid] = useState({});
  const [isChecked, setIsChecked] = useState(false);
  const [, , postalCodeHelperProps] = useField('destinationPostalCode');
  const [, , checkBoxHelperProps] = useField('useDutyZIP');

  const setDestinationZipToDutyZip = (checkboxValue) => {
    if (checkboxValue) {
      postalCodeHelperProps.setValue(dutyZip);
      checkBoxHelperProps.setValue('checked');
      setIsChecked(true);
    } else {
      postalCodeHelperProps.setValue('');
      checkBoxHelperProps.setValue('');
      setIsChecked(false);
    }
  };

  const handlePrefillPostalCodeChange = (value) => {
    if (isChecked && value !== dutyZip) {
      checkBoxHelperProps.setValue('');
      setIsChecked(false);
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
    <SectionWrapper className={styles.DestinationZIPInfo}>
      <h2>Destination info</h2>
      <div className="display-inline-block" data-testid="destinationZIP">
        <TextField
          label="Destination ZIP"
          id="destinationPostalCode"
          name="destinationPostalCode"
          maxLength={5}
          onChange={(e) => {
            handlePrefillPostalCodeChange(e.target.value);
          }}
          validate={(value) => postalCodeValidate(value, 'destination', 'secondaryDestinationPostalCode')}
        />
      </div>
      <CheckboxField
        id="useDutyZIP"
        name="useDutyZIP"
        label="Use ZIP for new duty location"
        onChange={(e) => setDestinationZipToDutyZip(e.target.checked)}
      />
      <div className="display-inline-block" data-testid="secondDestinationZIP">
        <TextField
          label="Second destination ZIP"
          id="secondDestinationPostalCode"
          name="secondDestinationPostalCode"
          maxLength={5}
          optional
          validate={(value) => postalCodeValidate(value, 'destination', 'secondaryDestinationPostalCode')}
        />
      </div>
    </SectionWrapper>
  );
};

DestinationZIPInfo.propTypes = {
  dutyZip: PropTypes.string,
  postalCodeValidator: PropTypes.func.isRequired,
};

DestinationZIPInfo.defaultProps = {
  dutyZip: undefined,
};

export default DestinationZIPInfo;
