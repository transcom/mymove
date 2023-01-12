import { useField } from 'formik';
import React from 'react';
import PropTypes from 'prop-types';

import './DropdownInput.module.scss';
import LocationSearchBox from 'components/LocationSearchBox/LocationSearchBox';
import { SearchTransportationOffices } from 'services/ghcApi';

// async function showAddress(addressId) {
//   return 'nope';
// }
// TODO: refactor component when we can to make it more user friendly with Formik
export const CloseoutOfficeInput = (props) => {
  const { label, name, displayAddress, hint, placeholder, isDisabled } = props;
  const [field, meta, helpers] = useField(props);

  const errorString = meta.value?.name ? meta.error?.name || meta.error : '';

  // console.log('CloseoutOfficeInput field', name, field.value);
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
      searchLocations={SearchTransportationOffices}
    />
  );
  // return (
  //   <DutyLocationSearchBoxComponent
  //     searchDutyLocations={SearchTransportationOffices}
  //     showAddress={showAddress}
  //     title={label}
  //     name={name}
  //     input={{
  //       value: field.value,
  //       onChange: helpers.setValue,
  //       name,
  //     }}
  //     errorMsg={errorString}
  //     displayAddress={displayAddress}
  //     hint={hint}
  //     placeholder={placeholder}
  //     isDisabled={isDisabled}
  //   />
  // );
};

CloseoutOfficeInput.propTypes = {
  // label displayed for input
  label: PropTypes.string.isRequired,
  // name is for the input
  name: PropTypes.string.isRequired,
  displayAddress: PropTypes.bool,
  hint: PropTypes.node,
  placeholder: PropTypes.string,
  isDisabled: PropTypes.bool,
};

CloseoutOfficeInput.defaultProps = {
  displayAddress: true,
  hint: '',
  placeholder: '',
  isDisabled: false,
};

export default CloseoutOfficeInput;
