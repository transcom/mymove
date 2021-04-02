import React, { Component, Fragment } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { debounce, get } from 'lodash';
import AsyncSelect from 'react-select/async';
import Alert from 'shared/Alert';
import { components } from 'react-select';
import Highlighter from 'react-highlight-words';
import { NULL_UUID } from 'shared/constants';
import { SearchDutyStations, ShowAddress } from './api.js';

import 'pages/MyMove/Profile/DutyStation.css';
import styles from './DutyStationSearchBox.module.scss';

const inputDebounceTime = 200;
const minSearchLength = 2;
const getOptionName = (option) => (option ? option.name : '');

export class DutyStationSearchBox extends Component {
  constructor(props) {
    super(props);

    this.state = {
      inputValue: '',
    };

    this.loadOptions = this.loadOptions.bind(this);
    this.getDebouncedOptions = this.getDebouncedOptions.bind(this);
    this.debouncedLoadOptions = this.debouncedLoadOptions.bind(this);
    this.localOnChange = this.localOnChange.bind(this);
    this.onInputChange = this.onInputChange.bind(this);
    this.renderOption = this.renderOption.bind(this);
    this.noOptionsMessage = this.noOptionsMessage.bind(this);
  }

  loadOptions(inputValue, callback) {
    if (inputValue && inputValue.length >= minSearchLength) {
      return SearchDutyStations(inputValue)
        .then((item) => {
          this.setState({
            error: null,
          });
          callback(item);
        })
        .catch((err) => {
          this.setState({
            error: err,
          });
          callback([]);
        });
    } else {
      callback([]);
    }
  }

  debouncedLoadOptions = debounce(this.loadOptions, inputDebounceTime);

  getDebouncedOptions(inputValue, callback) {
    if (!inputValue) {
      return callback(null);
    }
    this.debouncedLoadOptions(inputValue, callback);
  }

  localOnChange(value) {
    if (value && value.id) {
      return ShowAddress(value.address_id).then((item) => {
        value.address = item;
        this.props.input.onChange(value);
        return value;
      });
    } else {
      this.props.input.onChange(null);
      return null;
    }
  }

  onInputChange(inputValue, { action }) {
    this.setState({ inputValue });
    return inputValue;
  }
  noOptionsMessage(props) {
    if (this.state.inputValue.length < minSearchLength) {
      return <span />;
    }
    return <span>No Options</span>;
  }
  renderOption(props) {
    // React throws an error complaining about use of this property, so we delete it
    delete props.innerProps.innerRef;
    return (
      <div {...props.innerProps}>
        <components.Option {...props}>
          <Highlighter autoEscape searchWords={[this.state.inputValue]} textToHighlight={props.label} />
        </components.Option>
      </div>
    );
  }
  render() {
    const { errorMsg, displayAddress } = this.props;
    const defaultTitle = 'Name of Duty Station:';
    const inputContainerClasses = classNames({ 'usa-input-error': errorMsg });
    const searchBoxHeaderClasses = classNames({ 'duty-station-header': errorMsg });
    const dutyInputClasses = classNames('duty-input-box', {
      [this.props.input.name]: true,
      'duty-input-box-error': errorMsg,
    });
    const locationClasses = classNames({ location: true, 'location-error': errorMsg });
    // api for duty station always returns an object, even when duty station is not set
    // if there is no duty station, that object will have a null uuid
    const isEmptyStation = get(this.props, 'input.value.id', NULL_UUID) === NULL_UUID;
    const title = this.props.title || defaultTitle;
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
      placeholder: (provided) => ({
        color: uswdsBlack,
      }),
    };
    return (
      <Fragment>
        <div className="duty-station-search usa-form-group">
          {this.state.error && (
            <div className="usa-width-one-whole error-message">
              <Alert type="error" heading="An error occurred">
                {this.state.error.message}
              </Alert>
            </div>
          )}
          <div className={inputContainerClasses}>
            <label
              className={`${styles.title} ${searchBoxHeaderClasses} usa-label`}
              htmlFor={`${this.props.name}-input`}
            >
              {errorMsg ? <strong>{title}</strong> : title}
            </label>
            <AsyncSelect
              name={this.props.name}
              inputId={`${this.props.name}-input`}
              className={dutyInputClasses}
              cacheOptions
              getOptionLabel={getOptionName}
              getOptionValue={getOptionName}
              loadOptions={this.getDebouncedOptions}
              onChange={this.localOnChange}
              onInputChange={this.onInputChange}
              components={{ Option: this.renderOption }}
              value={isEmptyStation ? null : this.props.input.value}
              noOptionsMessage={this.noOptionsMessage}
              placeholder="Start typing a duty station..."
              styles={customStyles}
            />
            {displayAddress && !isEmptyStation && (
              <p className={locationClasses}>
                {this.props.input.value.address.city}, {this.props.input.value.address.state}{' '}
                {this.props.input.value.address.postal_code}
              </p>
            )}
            {this.props.errorMsg && <span className="usa-error-message">{this.props.errorMsg}</span>}
          </div>
        </div>
      </Fragment>
    );
  }
}
DutyStationSearchBox.propTypes = {
  onChange: PropTypes.func,
  existingStation: PropTypes.object,
  title: PropTypes.string,
  name: PropTypes.string,
  displayAddress: PropTypes.bool,
  errorMsg: PropTypes.string,
};

DutyStationSearchBox.defaultProps = {
  displayAddress: true,
  errorMsg: undefined,
};

export default DutyStationSearchBox;
