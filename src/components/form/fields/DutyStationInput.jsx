import { useField } from 'formik';
import React from 'react';
import PropTypes from 'prop-types';

import { DutyStationSearchBox } from 'scenes/ServiceMembers/DutyStationSearchBox';
import { DutyStationShape } from 'types/dutyStation';

// TODO: refactor component when we can to make it more user friendly with Formik
export const DutyStationInput = (props) => {
  //  react/prop-types
  const { label, name, displayAddress } = props;
  //  no-unused-vars
  const [field, meta, helpers] = useField(props);
  return (
    <DutyStationSearchBox
      title={label}
      name={name}
      input={{
        value: field.value,
        onChange: helpers.setValue,
      }}
      errorMsg={meta.error}
      displayAddress={displayAddress}
    />
  );
};

DutyStationInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  // duty station value
  //  react/no-unused-prop-types
  value: DutyStationShape,
  // name is for the input
  name: PropTypes.string.isRequired,
  displayAddress: PropTypes.bool,
};

DutyStationInput.defaultProps = {
  value: {},
  displayAddress: true,
};

export default DutyStationInput;
