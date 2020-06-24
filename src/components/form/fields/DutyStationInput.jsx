import { useField } from 'formik';
import React from 'react';
import PropTypes from 'prop-types';

import { DutyStationSearchBox } from 'scenes/ServiceMembers/DutyStationSearchBox';
import { DutyStationShape } from 'types/dutyStation';

// TODO: refactor component when we can to make it more user friendly with Formik
export const DutyStationInput = (props) => {
  // eslint-disable-next-line react/prop-types
  const { label, name } = props;
  // eslint-disable-next-line no-unused-vars
  const [field, meta, helpers] = useField(props);
  return (
    <DutyStationSearchBox
      title={label}
      name={name}
      input={{
        value: field.value,
        onChange: helpers.setValue,
      }}
    />
  );
};

DutyStationInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  // duty station value
  // eslint-disable-next-line react/no-unused-prop-types
  value: DutyStationShape,
  // name is for the input
  name: PropTypes.string.isRequired,
};

DutyStationInput.defaultProps = {
  value: {},
};

export default DutyStationInput;
