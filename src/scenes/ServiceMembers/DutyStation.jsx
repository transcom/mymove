import { debounce, sortBy } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import AsyncSelect from 'react-select/lib/Async';
import { components } from 'react-select';
import { connect } from 'react-redux';
import Highlighter from 'react-highlight-words';
import { bindActionCreators } from 'redux';

import { updateServiceMember, loadServiceMember } from './ducks';
import { SearchDutyStations } from './api.js';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

import './DutyStation.css';

const inputDebounceTime = 200;
const minSearchLength = 2;
const getOptionName = option => (option ? option.name : '');

export class DutyStation extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: null,
      inputValue: '',
    };
    this.loadOptions = this.loadOptions.bind(this);
    this.getDebouncedOptions = this.getDebouncedOptions.bind(this);
    this.debouncedLoadOptions = this.debouncedLoadOptions.bind(this);
    this.onChange = this.onChange.bind(this);
    this.onInputChange = this.onInputChange.bind(this);
    this.renderOption = this.renderOption.bind(this);
  }

  componentDidMount() {
    this.props.loadServiceMember(this.props.match.params.serviceMemberId);
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

  handleSubmit = () => {
    if (this.state.value) {
      this.props.updateServiceMember({ current_station: this.state.value });
    }
  };

  loadOptions(inputValue, callback) {
    if (
      this.props.currentServiceMember &&
      inputValue &&
      inputValue.length >= minSearchLength
    ) {
      return SearchDutyStations(
        this.props.currentServiceMember.branch,
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

  onChange(value) {
    if (value && value.id) {
      this.setState({ value });
      return value;
    } else {
      this.setState({ value: null });
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
    const { pages, pageKey, hasSubmitSuccess, error } = this.props;
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={this.state.isValid}
        hasSucceeded={hasSubmitSuccess}
        error={error || this.state.error}
      >
        <form className="duty-station" onSubmit={no_op}>
          <h1 className="sm-heading">Current Duty Station</h1>
          <label>Name of Duty Station</label>
          <AsyncSelect
            cacheOptions
            inputValue={this.state.inputValue}
            getOptionLabel={getOptionName}
            getOptionValue={getOptionName}
            loadOptions={this.getDebouncedOptions}
            onChange={this.onChange}
            onInputChange={this.onInputChange}
            components={{ Option: this.renderOption }}
            placeholder="Start typing a duty station..."
          />
          {this.state.value && (
            <ul className="station-description">
              <li>{this.state.value.name}</li>
              <li>
                {this.state.value.address.city},{' '}
                {this.state.value.address.state}{' '}
                {this.state.value.address.postal_code}
              </li>
            </ul>
          )}
        </form>
      </WizardPage>
    );
  }
}
DutyStation.propTypes = {
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { updateServiceMember, loadServiceMember },
    dispatch,
  );
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
export default connect(mapStateToProps, mapDispatchToProps)(DutyStation);
