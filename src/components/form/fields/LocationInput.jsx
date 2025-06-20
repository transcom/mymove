import { useField } from 'formik';
import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import './DropdownInput.module.scss';
import LocationSearchBox from 'components/LocationSearchBox/LocationSearchBox';
import { searchLocationByZipCityState } from 'services/internalApi';
import { searchLocationByZipCityState as ghcSearchLocationByZipCityState } from 'services/ghcApi';
import { selectLoggedInUser } from 'store/entities/selectors';
import { OfficeUserInfoShape } from 'types/index';

export const LocationInput = (props) => {
  const {
    label,
    name,
    displayAddress,
    placeholder,
    isDisabled,
    handleLocationChange,
    officeUser,
    includePOBoxes,
    showRequiredAsteriskForLocationLookup,
  } = props;
  const [field, meta] = useField(props);
  const errorString = meta.value?.name ? meta.error?.name || meta.error : '';

  return (
    <LocationSearchBox
      title={label}
      name={name}
      input={{
        value: field.value,
        onChange: handleLocationChange,
        locationState: () => {},
        name,
      }}
      required={showRequiredAsteriskForLocationLookup}
      showRequiredAsterisk={showRequiredAsteriskForLocationLookup}
      errorMsg={errorString}
      displayAddress={displayAddress}
      placeholder={placeholder}
      isDisabled={isDisabled}
      searchLocations={officeUser?.id ? ghcSearchLocationByZipCityState : searchLocationByZipCityState}
      handleLocationOnChange={handleLocationChange}
      includePOBoxes={includePOBoxes}
    />
  );
};

LocationInput.propTypes = {
  label: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  displayAddress: PropTypes.bool,
  placeholder: PropTypes.string,
  isDisabled: PropTypes.bool,
  handleLocationChange: PropTypes.func.isRequired,
  officeUser: OfficeUserInfoShape,
  includePOBoxes: PropTypes.bool,
};

LocationInput.defaultProps = {
  displayAddress: false,
  placeholder: '',
  isDisabled: false,
  officeUser: {},
  includePOBoxes: false,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
  };
};

const mapDispatchToProps = {};

export default connect(mapStateToProps, mapDispatchToProps)(LocationInput);
