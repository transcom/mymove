import { useField } from 'formik';
import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import './DropdownInput.module.scss';
import { searchCountry } from 'services/internalApi';
import { searchCountry as ghcSearchCountry } from 'services/ghcApi';
import { selectLoggedInUser } from 'store/entities/selectors';
import { OfficeUserInfoShape } from 'types/index';
import CountrySearchBox from 'components/CountrySearchBox/CountrySearchBox';

export const CountryInput = (props) => {
  const { label, name, displayAddress, placeholder, isDisabled, handleCountryChange, officeUser } = props;
  const [field, meta] = useField(props);
  const errorString = meta.value?.name ? meta.error?.name || meta.error : '';

  return (
    <CountrySearchBox
      title={label}
      name={name}
      input={{
        value: field.value,
        onChange: handleCountryChange,
        countryState: () => {},
        name,
      }}
      hint="Required"
      errorMsg={errorString}
      displayAddress={displayAddress}
      placeholder={placeholder}
      isDisabled={isDisabled}
      searchCountries={officeUser?.id ? ghcSearchCountry : searchCountry}
      handleCountryOnChange={handleCountryChange}
    />
  );
};

CountryInput.propTypes = {
  label: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  displayAddress: PropTypes.bool,
  placeholder: PropTypes.string,
  isDisabled: PropTypes.bool,
  handleCountryChange: PropTypes.func.isRequired,
  officeUser: OfficeUserInfoShape,
};

CountryInput.defaultProps = {
  displayAddress: false,
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

export default connect(mapStateToProps, mapDispatchToProps)(CountryInput);
