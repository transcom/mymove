import { useField } from 'formik';
import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import './DropdownInput.module.scss';
import LocationSearchBox from 'components/LocationSearchBox/LocationSearchBox';
import { searchLocationByZipCity } from 'services/internalApi';
import { searchLocationByZipCity as ghcSearchLocationByZipCity } from 'services/ghcApi';
import { selectLoggedInUser } from 'store/entities/selectors';
import { OfficeUserInfoShape } from 'types/index';

export const ZipCityInput = (props) => {
  const { label, name, displayAddress, hint, placeholder, isDisabled, handleZipCityChange, officeUser } = props;
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
      searchLocations={officeUser?.id ? ghcSearchLocationByZipCity : searchLocationByZipCity}
      handleOnChange={handleZipCityChange}
    />
  );
};

ZipCityInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  // name is for the input
  name: PropTypes.string.isRequired,
  displayAddress: PropTypes.bool,
  hint: PropTypes.node,
  placeholder: PropTypes.string,
  isDisabled: PropTypes.bool,
  handleZipCityChange: PropTypes.func.isRequired,
  officeUser: OfficeUserInfoShape,
};

ZipCityInput.defaultProps = {
  displayAddress: false,
  hint: '',
  placeholder: '',
  isDisabled: false,
  officeUser: {},
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
  };
};

const mapDispatchToProps = {};

export default connect(mapStateToProps, mapDispatchToProps)(ZipCityInput);
