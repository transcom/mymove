import { debounce, sortBy } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import AsyncSelect from 'react-select/lib/Async';
import { components } from 'react-select';
import { connect } from 'react-redux';
import Highlighter from 'react-highlight-words';
import { bindActionCreators } from 'redux';

import { loadServiceMember } from './ducks';
import { SearchDutyStations } from './api.js';

import './DutyStation.css';

const inputDebounceTime = 200;
const minSearchLength = 2;
const getOptionName = option => (option ? option.name : '');

export class DutyStationSearchBox extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: null,
      inputValue: '',
    };
    this.loadOptions = this.loadOptions.bind(this);
    this.getDebouncedOptions = this.getDebouncedOptions.bind(this);
    this.debouncedLoadOptions = this.debouncedLoadOptions.bind(this);
    this.localOnChange = this.localOnChange.bind(this);
    this.onInputChange = this.onInputChange.bind(this);
    this.renderOption = this.renderOption.bind(this);
  }
  static getDerivedStateFromProps(nextProps, prevState) {
    if (prevState.value === null) {
      return {
        value: nextProps.existingStation,
        inputValue: getOptionName(nextProps.existingStation),
      };
    }
    return {};
  }

  loadOptions(inputValue, callback) {
    if (
      this.props.currentServiceMember &&
      inputValue &&
      inputValue.length >= minSearchLength
    ) {
      return SearchDutyStations(
        this.props.currentServiceMember.affiliation,
        inputValue,
      )
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
      this.setState({ value });
      this.props.onChange(value);
      return value;
    } else {
      this.setState({ value: null });
      this.props.onChange(null);
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
        {this.state.value && (
          <ul className="station-description">
            <li>{this.state.value.name}</li>
            <li>
              {this.state.value.address.city}, {this.state.value.address.state}{' '}
              {this.state.value.address.postal_code}
            </li>
          </ul>
        )}
      </Fragment>
    );
  }
}
DutyStationSearchBox.propTypes = {
  currentServiceMember: PropTypes.object,
  onChange: PropTypes.func.isRequired,
  existingStation: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadServiceMember }, dispatch);
}
function mapStateToProps(state) {
  const currentServiceMember = state.serviceMember.currentServiceMember;
  const dutyStation =
    currentServiceMember && currentServiceMember.current_station
      ? currentServiceMember.current_station
      : null;
  const props = {
    existingStation: dutyStation,
    ...state.serviceMember,
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(
  DutyStationSearchBox,
);
