/* eslint-disable react/jsx-props-no-spreading */
import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { FormGroup, Label } from '@trussworks/react-uswds';
import AsyncSelect from 'react-select/async';
import classNames from 'classnames';

import styles from './DutyStationSearchBox.module.scss';

import { SearchDutyStations, ShowAddress } from 'scenes/ServiceMembers/api';
import { DutyStationShape } from 'types';

const getOptionName = (option) => (option ? option.name : '');

const uswdsBlack = '#565c65';
const uswdsBlue = '#2491ff';

const customStyles = {
  control: (provided, state) => ({
    ...provided,
    borderRadius: '0px',
    borderColor: uswdsBlack,
    padding: '0.1rem',
    maxWidth: '32rem',
    ':hover': {
      ...styles[':hover'],
      borderColor: uswdsBlack,
    },
    boxShadow: state.isFocused ? `0 0 0 0.26667rem ${uswdsBlue}` : '',
  }),
  dropdownIndicator: (provided) => ({
    ...provided,
    color: uswdsBlack,
    ':hover': {
      ...styles[':hover'],
      color: uswdsBlack,
    },
  }),
  indicatorSeparator: (provided) => ({
    ...provided,
    backgroundColor: uswdsBlack,
  }),
  placeholder: () => ({
    color: uswdsBlack,
  }),
};

export const DutyStationSearchBoxComponent = (props) => {
  const { searchDutyStations, showAddress, title, input, name, errorMsg, displayAddress } = props;
  const { value, onChange, name: inputName } = input;

  const [inputValue, setInputValue] = useState('');

  const loadOptions = async (query) => {
    if (!query) {
      return null;
    }

    try {
      const stations = await searchDutyStations(query);
      return stations;
    } catch (e) {
      return null;
    }
  };

  const selectOption = async (selectedValue) => {
    if (selectedValue && selectedValue.id) {
      const newValue = {
        ...selectedValue,
        address: await showAddress(selectedValue.address_id),
      };

      onChange(newValue);
      return newValue;
    }

    onChange(null);
    return null;
  };

  const changeInputText = (text) => {
    setInputValue(text);
  };

  const inputId = `${name}-input`;

  const inputContainerClasses = classNames({ 'usa-input-error': errorMsg });
  const locationClasses = classNames('location', { 'location-error': errorMsg });
  const labelClasses = classNames(styles.title, {
    [styles.titleWithError]: errorMsg,
  });
  const dutyInputClasses = classNames('duty-input-box', {
    [inputName]: true,
    'duty-input-box-error': errorMsg,
  });

  const noOptionsMessage = () => (inputValue.length ? 'No Options' : '');

  return (
    <FormGroup>
      <div className={inputContainerClasses}>
        <Label htmlFor={inputId} className={labelClasses}>
          {title}
        </Label>
        <AsyncSelect
          name={name}
          inputId={inputId}
          className={dutyInputClasses}
          cacheOptions
          getOptionLabel={getOptionName}
          getOptionValue={getOptionName}
          loadOptions={loadOptions}
          onChange={selectOption}
          onInputChange={changeInputText}
          placeholder="Start typing a duty station..."
          styles={customStyles}
          value={value}
          noOptionsMessage={noOptionsMessage}
        />
      </div>
      {displayAddress && !!value && (
        <p className={locationClasses}>
          {value.address.city}, {value.address.state} {value.address.postal_code}
        </p>
      )}
      {errorMsg && <span className="usa-error-message">{errorMsg}</span>}
    </FormGroup>
  );
};

export const DutyStationSearchBoxContainer = (props) => {
  return <DutyStationSearchBoxComponent {...props} searchDutyStations={SearchDutyStations} showAddress={ShowAddress} />;
};

DutyStationSearchBoxContainer.propTypes = {
  displayAddress: PropTypes.bool,
  name: PropTypes.string.isRequired,
  errorMsg: PropTypes.string,
  title: PropTypes.string,
  input: PropTypes.shape({
    name: PropTypes.string,
    onChange: PropTypes.func,
    value: DutyStationShape,
  }),
};

DutyStationSearchBoxContainer.defaultProps = {
  displayAddress: true,
  title: 'Name of Duty Station:',
  errorMsg: '',
  input: {
    name: '',
    onChange: () => {},
    value: undefined,
  },
};

DutyStationSearchBoxComponent.propTypes = {
  ...DutyStationSearchBoxContainer.propTypes,
  searchDutyStations: PropTypes.func.isRequired,
  showAddress: PropTypes.func.isRequired,
};

DutyStationSearchBoxComponent.defaultProps = {
  ...DutyStationSearchBoxContainer.defaultProps,
};

export default DutyStationSearchBoxContainer;
