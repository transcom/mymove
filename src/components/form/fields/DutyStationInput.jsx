import { useField } from 'formik';
import React from 'react';
import PropTypes from 'prop-types';

import DutyStationSearchBox from 'components/DutyStationSearchBox/DutyStationSearchBox';

// TODO: refactor component when we can to make it more user friendly with Formik
export const DutyStationInput = (props) => {
  const { label, name, displayAddress, hint } = props;
  const [field, meta, helpers] = useField(props);

  const errorString = meta.value?.name ? meta.error?.name || meta.error : '';

  return (
    <DutyStationSearchBox
      title={label}
      name={name}
      input={{
        value: field.value,
        onChange: helpers.setValue,
        name,
      }}
      errorMsg={errorString}
      displayAddress={displayAddress}
      hint={hint}
    />
  );
};

DutyStationInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  // name is for the input
  name: PropTypes.string.isRequired,
  displayAddress: PropTypes.bool,
  hint: PropTypes.node,
};

DutyStationInput.defaultProps = {
  displayAddress: true,
  hint: '',
};

export default DutyStationInput;
