/* eslint-disable react/jsx-props-no-spreading */
import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { FormGroup, Label } from '@trussworks/react-uswds';
import AsyncSelect from 'react-select/async';
import classNames from 'classnames';
import { debounce } from 'lodash';

import styles from './DutyStationSearchBox.module.scss';
import { SearchDutyStations, ShowAddress } from './api';

import { DutyStationShape } from 'types';

const getOptionName = (option) => option.name;

const formatOptionLabel = (option, input) => {
  const { name } = option;
  const { inputValue } = input;

  const optionLabel = name || '';
  const inputText = inputValue || '';

  const searchIndex = optionLabel.toLowerCase().indexOf(inputText.toLowerCase());

  if (searchIndex === -1) {
    return <span>{optionLabel}</span>;
  }

  return (
    <span>
      {optionLabel.substr(0, searchIndex)}
      <mark>{optionLabel.substr(searchIndex, inputText.length)}</mark>
      {optionLabel.substr(searchIndex + inputText.length)}
    </span>
  );
};

const uswdsBlack = '#565c65';
const uswdsBlue = '#2491ff';

const MIN_SEARCH_LENGTH = 2;
const DEBOUNCE_TIMER_MS = 200;

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
  valueContainer: (provided) => ({
    ...provided,
    display: 'flex',
  }),
};

export const DutyStationSearchBoxComponent = (props) => {
  const { searchDutyStations, showAddress, title, input, name, errorMsg, displayAddress } = props;
  const { value, onChange, name: inputName } = input;

  const [inputValue, setInputValue] = useState('');

  const loadOptions = debounce((query, callback) => {
    if (!query || query.length < MIN_SEARCH_LENGTH) {
      callback(null);
      return undefined;
    }

    searchDutyStations(query)
      .then((stations) => {
        callback(stations);
      })
      .catch(() => {
        callback(null);
      });

    return undefined;
  }, DEBOUNCE_TIMER_MS);

  const selectOption = async (selectedValue) => {
    const address = await showAddress(selectedValue.address_id);
    const newValue = {
      ...selectedValue,
      address,
    };

    onChange(newValue);
    return newValue;
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
    [styles.dutyInputBoxError]: errorMsg,
  });

  const noOptionsMessage = () => (inputValue.length ? 'No Options' : '');
  const hasStation = !!value && !!value.address;

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
          formatOptionLabel={formatOptionLabel}
          getOptionValue={getOptionName}
          loadOptions={loadOptions}
          onChange={selectOption}
          onInputChange={changeInputText}
          placeholder="Start typing a duty station..."
          value={hasStation ? value : null}
          noOptionsMessage={noOptionsMessage}
          styles={customStyles}
        />
      </div>
      {displayAddress && hasStation && (
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
