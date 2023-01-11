/* eslint-disable react/jsx-props-no-spreading */
import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { FormGroup, Label } from '@trussworks/react-uswds';
import AsyncSelect from 'react-select/async';
import classNames from 'classnames';
import { debounce } from 'lodash';

import styles from './LocationSearchBox.module.scss';
import { SearchDutyLocations, ShowAddress } from './api';

import Hint from 'components/Hint';
import { DutyLocationShape } from 'types';

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

export const LocationSearchBoxComponent = ({
  searchLocations,
  showAddress,
  title,
  input,
  name,
  errorMsg,
  displayAddress,
  hint,
  placeholder,
  isDisabled,
}) => {
  const { value, onChange, name: inputName } = input;

  const [inputValue, setInputValue] = useState('');
  let disabledStyles = {};
  if (isDisabled) {
    disabledStyles = {
      ...customStyles,
      control: (provided, state) => ({
        ...provided,
        backgroundColor: '#DCDEE0',
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
      singleValue: () => {
        const color = '#1B1B1B';
        return { color };
      },
      valueContainer: (provided) => ({
        ...provided,
        display: 'flex',
        backgroundColor: '#DCDEE0',
      }),
    };
  }

  const loadOptions = debounce((query, callback) => {
    if (!query || query.length < MIN_SEARCH_LENGTH) {
      callback(null);
      return undefined;
    }

    searchLocations(query)
      .then((locations) => {
        callback(locations);
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
  const hasLocation = !!value && !!value.address;
  return (
    <FormGroup>
      <div className="labelWrapper">
        <Label htmlFor={inputId} className={labelClasses}>
          {title}
        </Label>
      </div>
      {hint && <Hint className={styles.hint}>{hint}</Hint>}
      <div className={inputContainerClasses}>
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
          placeholder={placeholder || 'Start typing a duty location...'}
          value={hasLocation ? value : null}
          noOptionsMessage={noOptionsMessage}
          styles={isDisabled ? disabledStyles : customStyles}
          isDisabled={isDisabled}
        />
      </div>
      {displayAddress && hasLocation && (
        <p className={locationClasses}>
          {value.address.city}, {value.address.state} {value.address.postalCode}
        </p>
      )}
      {errorMsg && <span className="usa-error-message">{errorMsg}</span>}
    </FormGroup>
  );
};

export const LocationSearchBoxContainer = (props) => {
  const { isDisabled, searchLocations } = props;
  return (
    <LocationSearchBoxComponent
      {...props}
      searchLocations={searchLocations || SearchDutyLocations}
      showAddress={ShowAddress}
      isDisabled={isDisabled}
    />
  );
};

LocationSearchBoxContainer.propTypes = {
  displayAddress: PropTypes.bool,
  name: PropTypes.string.isRequired,
  errorMsg: PropTypes.string,
  title: PropTypes.string,
  input: PropTypes.shape({
    name: PropTypes.string,
    onChange: PropTypes.func,
    value: DutyLocationShape,
  }),
  hint: PropTypes.node,
  placeholder: PropTypes.string,
  isDisabled: PropTypes.bool,
  searchLocations: PropTypes.func.isRequired,
};

LocationSearchBoxContainer.defaultProps = {
  displayAddress: true,
  title: 'Name of Duty Location:',
  errorMsg: '',
  input: {
    name: '',
    onChange: () => {},
    value: undefined,
  },
  hint: '',
  placeholder: 'Start typing a duty location...',
  isDisabled: false,
};

LocationSearchBoxComponent.propTypes = {
  ...LocationSearchBoxContainer.propTypes,
  searchLocations: PropTypes.func.isRequired,
  showAddress: PropTypes.func.isRequired,
  isDisabled: PropTypes.bool,
};

LocationSearchBoxComponent.defaultProps = {
  ...LocationSearchBoxContainer.defaultProps,
  isDisabled: false,
};

export default LocationSearchBoxContainer;
