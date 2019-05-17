import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';
import { Field } from 'redux-form';
import { get } from 'lodash';
import { updateServiceMember } from './ducks';
import { NULL_UUID } from 'shared/constants';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';

import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

import './DutyStation.css';

const validateDutyStationForm = (values, form) => {
  let errors = {};
  // api for duty station always returns an object, even when duty station is not set
  // if there is no duty station, that object will have a null uuid
  if (get(values, 'current_station.id', NULL_UUID) === NULL_UUID) {
    const newError = {
      current_station: 'Please select a duty station.',
    };
    return newError;
  }
  return errors;
};

const dutyStationFormName = 'duty_station';
const DutyStationWizardForm = reduxifyWizardForm(dutyStationFormName, validateDutyStationForm);

export class DutyStation extends Component {
  constructor(props) {
    super(props);

    this.state = {
      value: null,
    };
    this.stationOnChange = this.stationOnChange.bind(this);
  }

  stationOnChange = newStation => {
    this.setState({ value: newStation });
  };

  handleSubmit = (somethings, elses) => {
    const pendingValues = this.props.values;
    if (pendingValues) {
      return this.props.updateServiceMember({
        current_station_id: pendingValues.current_station.id,
      });
    }
  };

  render() {
    const { pages, pageKey, error, existingStation, newDutyStation, currentStation } = this.props;

    let initialValues = null;
    if (existingStation.name) {
      initialValues = { current_station: existingStation };
    }

    const newDutyStationErrorMsg =
      newDutyStation.name && newDutyStation.name === currentStation.name
        ? 'You entered the same duty station for your origin and destination. Please change one of them.'
        : '';

    return (
      <DutyStationWizardForm
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        initialValues={initialValues}
        serverError={error}
      >
        <h1 className="sm-heading">Current Duty Station</h1>
        <Field name="current_station" component={DutyStationSearchBox} errorMsg={newDutyStationErrorMsg} />
      </DutyStationWizardForm>
    );
  }
}
DutyStation.propTypes = {
  error: PropTypes.object,
  updateServiceMember: PropTypes.func.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateServiceMember }, dispatch);
}
function mapStateToProps(state) {
  const formValues = getFormValues(dutyStationFormName)(state);
  return {
    values: getFormValues(dutyStationFormName)(state),
    existingStation: get(state, 'serviceMember.currentServiceMember.current_station', {}),
    ...state.serviceMember,
    currentStation: get(formValues, 'current_station', {}),
    newDutyStation: get(state, 'orders.currentOrders.new_duty_station', {}),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(DutyStation);
