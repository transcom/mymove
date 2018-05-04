import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { Field } from 'redux-form';

import { updateServiceMember, loadServiceMember } from './ducks';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';

import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

import './DutyStation.css';

const validateDutyStationForm = (values, form) => {
  let errors = {};

  if (!values.current_station) {
    const newError = {
      current_station: 'Please select a duty station.',
    };
    return newError;
  }
  return errors;
};

const dutyStationFormName = 'duty_station';
const DutyStationWizardForm = reduxifyWizardForm(
  dutyStationFormName,
  validateDutyStationForm,
);

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
  componentDidMount() {
    this.props.loadServiceMember(this.props.match.params.serviceMemberId);
  }

  handleSubmit = (somethings, elses) => {
    const pendingValues = this.props.formData.values;
    if (pendingValues) {
      this.props.updateServiceMember({
        current_station: pendingValues.current_station,
      });
    }
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      existingStation,
    } = this.props;
    // TODO: make sure isvalid is accurate

    let initialValues = null;
    if (existingStation) {
      initialValues = { current_station: existingStation };
    }
    return (
      <DutyStationWizardForm
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        initialValues={initialValues}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
      >
        <h1 className="sm-heading">Current Duty Station</h1>
        <Field
          name="current_station"
          component={DutyStationSearchBox}
          affiliation={this.props.affiliation}
        />
      </DutyStationWizardForm>
    );
  }
}
DutyStation.propTypes = {
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
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
  const affiliation = currentServiceMember
    ? currentServiceMember.affiliation
    : null;
  const props = {
    affiliation,
    formData: state.form[dutyStationFormName],
    existingStation: dutyStation,
    ...state.serviceMember,
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(DutyStation);
