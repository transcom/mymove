import { useField, useFormikContext } from 'formik';
import React from 'react';
import PropTypes from 'prop-types';

import LocationSearchBox from 'components/LocationSearchBox/LocationSearchBox';
import './DropdownInput.module.scss';

// TODO: refactor component when we can to make it more user friendly with Formik
export const DutyLocationInput = (props) => {
  const {
    label,
    name,
    displayAddress,
    hint,
    placeholder,
    isDisabled,
    searchLocations,
    metaOverride,
    onDutyLocationChange,
  } = props;
  const [field, meta, helpers] = useField(props);

  const { touched } = useFormikContext();
  const handleDutyLocationChange = (event) => {
    if (onDutyLocationChange) {
      onDutyLocationChange(event);
    }
  };

  let errorString = '';
  if (metaOverride && metaOverride !== '') {
    errorString = metaOverride;
  } else {
    errorString = meta.value?.name ? meta.error?.name || meta.error : '';
  }

  const handleChange = (value) => {
    touched[name] = true;
    helpers.setValue(value);
  };

  return (
    <LocationSearchBox
      title={label}
      name={name}
      input={{
        value: field.value,
        onChange: handleChange,
        locationState: handleDutyLocationChange,
        name,
      }}
      errorMsg={errorString}
      displayAddress={displayAddress}
      hint={hint}
      placeholder={placeholder}
      isDisabled={isDisabled}
      searchLocations={searchLocations}
    />
  );
};

DutyLocationInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  // name is for the input
  name: PropTypes.string.isRequired,
  displayAddress: PropTypes.bool,
  hint: PropTypes.node,
  placeholder: PropTypes.string,
  isDisabled: PropTypes.bool,
  searchLocations: PropTypes.func,
  metaOverride: PropTypes.string,
  onDutyLocationChange: PropTypes.func,
};

DutyLocationInput.defaultProps = {
  displayAddress: true,
  hint: '',
  placeholder: '',
  isDisabled: false,
  searchLocations: undefined,
  metaOverride: '',
  onDutyLocationChange: undefined,
};

export default DutyLocationInput;
