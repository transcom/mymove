import { debounce, sortBy, get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import AsyncSelect from 'react-select/lib/Async';
import Alert from 'shared/Alert';
import { components } from 'react-select';
import Highlighter from 'react-highlight-words';
import { NULL_UUID } from 'shared/constants';
import { SearchDutyStations } from './api.js';

import './DutyStation.css';

const inputDebounceTime = 200;
const minSearchLength = 2;
const getOptionName = option => (option ? option.name : '');

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
  }

  loadOptions(inputValue, callback) {
    if (inputValue && inputValue.length >= minSearchLength) {
      return SearchDutyStations(inputValue)
        .then(item => {
          this.setState({
            error: null,
          });
          const sorted = sortBy(item, 'name');
          callback(sorted);
        })
        .catch(err => {
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
      this.props.input.onChange(value);
      return value;
    } else {
      this.props.input.onChange(null);
      return null;
    }
  }

  onInputChange(inputValue, { action }) {
    this.setState({ inputValue });
    return inputValue;
  }

  renderOption(props) {
    // React throws an error complaining about use of this property, so we delete it
    delete props.innerProps.innerRef;
    return (
      <div {...props.innerProps}>
        <components.Option {...props}>
          <Highlighter searchWords={[this.state.inputValue]} textToHighlight={props.label} />
        </components.Option>
      </div>
    );
  }
  render() {
    const defaultTitle = 'Name of Duty Station:';
    // api for duty station always returns an object, even when duty station is not set
    // if there is no duty station, that object will have a null uuid
    const isEmptyStation = get(this.props, 'input.value.id', NULL_UUID) === NULL_UUID;
    return (
      <Fragment>
        <div className="duty-station-search">
          {this.state.error && (
            <div className="usa-width-one-whole error-message">
              <Alert type="error" heading="An error occurred">
                {this.state.error.message}
              </Alert>
            </div>
          )}
          <p>{this.props.title || defaultTitle}</p>
          <AsyncSelect
            className="duty-input-box"
            cacheOptions
            getOptionLabel={getOptionName}
            getOptionValue={getOptionName}
            loadOptions={this.getDebouncedOptions}
            onChange={this.localOnChange}
            onInputChange={this.onInputChange}
            components={{ Option: this.renderOption }}
            value={isEmptyStation ? null : this.props.input.value}
            placeholder="Start typing a duty station..."
          />
          {!isEmptyStation && (
            <p className="location">
              {this.props.input.value.address.city}, {this.props.input.value.address.state}{' '}
              {this.props.input.value.address.postal_code}
            </p>
          )}
        </div>
      </Fragment>
    );
  }
}
DutyStationSearchBox.propTypes = {
  onChange: PropTypes.func,
  existingStation: PropTypes.object,
  title: PropTypes.string,
};

export default DutyStationSearchBox;
