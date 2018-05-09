import { debounce, sortBy } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import AsyncSelect from 'react-select/lib/Async';
import Alert from 'shared/Alert';
import { components } from 'react-select';
import Highlighter from 'react-highlight-words';

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
          <Highlighter
            searchWords={[this.state.inputValue]}
            textToHighlight={props.label}
          />
        </components.Option>
      </div>
    );
  }
  render() {
    return (
      <Fragment>
        {this.state.error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {this.state.error.message}
            </Alert>
          </div>
        )}
        <label>Name of Duty Station</label>
        <AsyncSelect
          cacheOptions
          inputValue={this.state.inputValue}
          getOptionLabel={getOptionName}
          getOptionValue={getOptionName}
          loadOptions={this.getDebouncedOptions}
          onChange={this.localOnChange}
          onInputChange={this.onInputChange}
          components={{ Option: this.renderOption }}
          placeholder="Start typing a duty station..."
        />
        {this.props.input.value && (
          <ul className="station-description">
            <li>{this.props.input.value.name}</li>
            <li>
              {this.props.input.value.address.city},{' '}
              {this.props.input.value.address.state}{' '}
              {this.props.input.value.address.postal_code}
            </li>
          </ul>
        )}
      </Fragment>
    );
  }
}
DutyStationSearchBox.propTypes = {
  onChange: PropTypes.func,
  existingStation: PropTypes.object,
};

export default DutyStationSearchBox;
