import { useField } from 'formik';
import React from 'react';
import PropTypes from 'prop-types';

import './DropdownInput.module.scss';
import LocationSearchBox from 'components/LocationSearchBox/LocationSearchBox';
import { searchLocationByZipCity } from 'services/internalApi';

export const ZipCityInput = (props) => {
  const { label, name, displayAddress, hint, placeholder, isDisabled, handleZipCityChange } = props;
  const [field, meta, helpers] = useField(props);
  const errorString = meta.value?.name ? meta.error?.name || meta.error : '';

  return (
    <LocationSearchBox
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
      placeholder={placeholder}
      isDisabled={isDisabled}
      searchLocations={searchLocationByZipCity}
      handleZipCityOnChange={handleZipCityChange}
    />
  );
};

ZipCityInput.propTypes = {
  label: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  displayAddress: PropTypes.bool,
  hint: PropTypes.node,
  placeholder: PropTypes.string,
  isDisabled: PropTypes.bool,
  handleZipCityChange: PropTypes.func.isRequired,
};

ZipCityInput.defaultProps = {
  displayAddress: true,
  hint: '',
  placeholder: '',
  isDisabled: false,
};

export default ZipCityInput;
