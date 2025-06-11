/* eslint-disable react/jsx-props-no-spreading */
import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { ErrorMessage, FormGroup, Label } from '@trussworks/react-uswds';
import AsyncSelect from 'react-select/async';
import classNames from 'classnames';
import { debounce } from 'lodash';

import styles from './LocationSearchBox.module.scss';
import { SearchDutyLocations, ShowAddress } from './api';

import { DutyLocationShape } from 'types';
import RequiredAsterisk from 'components/form/RequiredAsterisk';

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

const formatLocation = (option, input) => {
  const { inputValue } = input;
  const outputLabel = `${option?.city || ''}, ${option?.state || ''} ${option?.postalCode || ''} (${
    option?.county || ''
  })`;
  const inputText = inputValue || '';

  const searchIndex = outputLabel.toLowerCase().indexOf(inputText.toLowerCase());

  if (searchIndex === -1) {
    return <span>{outputLabel}</span>;
  }

  return (
    <span>
      {outputLabel.substr(0, searchIndex)}
      <mark>{outputLabel.substr(searchIndex, inputText.length)}</mark>
      {outputLabel.substr(searchIndex + inputText.length)}
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
    maxWidth: '100%',
    '@media (max-width: 768px)': {
      maxWidth: '32em',
    },
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
  // fixes a bug with AsyncSelect highlighting all results blue
  option: (provided, state) => ({
    ...provided,
    backgroundColor: state.isFocused ? '#f0f0f0' : 'white', // Change background color on focus
    color: 'black', // Change text color
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
  handleLocationOnChange,
  showRequiredAsterisk,
}) => {
  const { value, onChange, locationState, name: inputName } = input;

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
    if (!selectedValue.address && !handleLocationOnChange) {
      const address = await showAddress(selectedValue.address_id);
      const newValue = {
        ...selectedValue,
        address,
      };
      locationState(newValue);
      onChange(newValue);
      return newValue;
    }

    locationState(selectedValue);
    onChange(selectedValue);

    if (handleLocationOnChange !== null) {
      handleLocationOnChange(selectedValue);
    }

    return selectedValue;
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

  const handleKeyDown = (event) => {
    if (event.key === 'Backspace' && !inputValue) {
      if (handleLocationOnChange) {
        handleLocationOnChange(null);
      } else {
        onChange(null);
      }
    }
  };

  const handleFocus = () => {
    if (handleLocationOnChange) {
      handleLocationOnChange(null);
    } else {
      onChange(null);
    }
  };

  const noOptionsMessage = () => (inputValue.length ? 'No Options' : '');
  const hasLocation = !!value && !!value.address;

  return (
    <FormGroup>
      <div className="labelWrapper">
        <Label hint={hint} htmlFor={inputId} className={labelClasses} data-testid={`${name}-label`}>
          <span aria-label="Required">
            {title} {showRequiredAsterisk && <RequiredAsterisk />}
          </span>
        </Label>
      </div>
      <div className={inputContainerClasses}>
        <AsyncSelect
          name={name}
          data-testid={inputId}
          inputId={inputId}
          className={dutyInputClasses}
          cacheOptions
          formatOptionLabel={handleLocationOnChange ? formatLocation : formatOptionLabel}
          getOptionValue={getOptionName}
          getOptionLabel={(option) => option.name}
          loadOptions={loadOptions}
          onChange={selectOption}
          onKeyDown={handleKeyDown}
          onInputChange={changeInputText}
          placeholder={placeholder || 'Start typing a duty location...'}
          value={
            (handleLocationOnChange && !!value && value.city != null && value.city !== '') ||
            (!handleLocationOnChange && hasLocation)
              ? value
              : ''
          }
          noOptionsMessage={noOptionsMessage}
          onFocus={handleFocus}
          styles={isDisabled ? disabledStyles : customStyles}
          isDisabled={isDisabled}
        />
      </div>
      {displayAddress && hasLocation && (
        <p className={locationClasses}>
          {value.address.city}, {value.address.state} {value.address.postalCode}
        </p>
      )}
      {errorMsg && <ErrorMessage>{errorMsg}</ErrorMessage>}
    </FormGroup>
  );
};

export const LocationSearchBoxContainer = (props) => {
  const { searchLocations } = props;

  return <LocationSearchBoxComponent {...props} searchLocations={searchLocations} showAddress={ShowAddress} />;
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
    locationState: PropTypes.func,
  }),
  hint: PropTypes.node,
  placeholder: PropTypes.string,
  isDisabled: PropTypes.bool,
  searchLocations: PropTypes.func,
  handleLocationOnChange: PropTypes.func,
};

LocationSearchBoxContainer.defaultProps = {
  displayAddress: true,
  title: 'Name of Duty Location:',
  errorMsg: '',
  input: {
    name: '',
    onChange: () => {},
    value: undefined,
    locationState: () => {},
  },
  hint: '',
  placeholder: 'Start typing a duty location...',
  isDisabled: false,
  searchLocations: SearchDutyLocations,
  handleLocationOnChange: null,
};

LocationSearchBoxComponent.propTypes = {
  ...LocationSearchBoxContainer.propTypes,
  searchLocations: PropTypes.func,
  showAddress: PropTypes.func.isRequired,
  isDisabled: PropTypes.bool,
  showRequiredAsterisk: PropTypes.bool,
};

LocationSearchBoxComponent.defaultProps = {
  ...LocationSearchBoxContainer.defaultProps,
  searchLocations: SearchDutyLocations,
  isDisabled: false,
  showRequiredAsterisk: false,
};

export default LocationSearchBoxContainer;
