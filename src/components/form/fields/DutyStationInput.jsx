import { useField } from 'formik';
import React from 'react';
import PropTypes from 'prop-types';

import { DutyStationSearchBox } from 'scenes/ServiceMembers/DutyStationSearchBox';

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
  // label optionally displayed for input
  label: PropTypes.string.isRequired,
  // duty station value
  value: PropTypes.shape({
    address: PropTypes.shape({
      city: PropTypes.string,
      id: PropTypes.string,
      postal_code: PropTypes.string,
      state: PropTypes.string,
      street_address_1: PropTypes.string,
    }),
    address_id: PropTypes.string,
    affiliation: PropTypes.string,
    created_at: PropTypes.string,
    id: PropTypes.string,
    name: PropTypes.string,
    updated_at: PropTypes.string,
  }),
  // name is for the input
  name: PropTypes.string.isRequired,
};

DutyStationInput.defaultProps = {
  value: {},
};

export default DutyStationInput;
