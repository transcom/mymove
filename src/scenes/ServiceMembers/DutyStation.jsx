import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { getFormValues } from 'redux-form';
import { Field } from 'redux-form';
import { get } from 'lodash';

import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { NULL_UUID } from 'shared/constants';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { selectActiveOrLatestOrders } from 'shared/Entities/modules/orders';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

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
      errorMessage: null,
    };
  }

  handleSubmit = () => {
    const { values, currentServiceMember, updateServiceMember } = this.props;

    if (values) {
      const payload = {
        id: currentServiceMember.id,
        current_station_id: values.current_station.id,
      };

      return patchServiceMember(payload)
        .then((response) => {
          updateServiceMember(response);
        })
        .catch((e) => {
          // TODO - error handling - below is rudimentary error handling to approximate existing UX
          // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to update service member due to server error');
          this.setState({
            errorMessage,
          });
        });
    }

    return Promise.resolve();
  };

  render() {
    const { pages, pageKey, error, existingStation, newDutyStation, currentStation } = this.props;
    const { errorMessage } = this.state;

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
        serverError={error || errorMessage}
      >
        <h1>Current duty station</h1>
        <SectionWrapper>
          <div className="tablet:margin-top-neg-3"></div>
          <Field
            name="current_station"
            title="What is your current duty station?"
            component={DutyStationSearchBox}
            errorMsg={newDutyStationErrorMsg}
          />
        </SectionWrapper>
      </DutyStationWizardForm>
    );
  }
}
DutyStation.propTypes = {
  error: PropTypes.object,
  updateServiceMember: PropTypes.func.isRequired,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

function mapStateToProps(state) {
  const formValues = getFormValues(dutyStationFormName)(state);
  const orders = selectActiveOrLatestOrders(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    values: getFormValues(dutyStationFormName)(state),
    existingStation: serviceMember?.current_station || {},
    // TODO
    ...state.serviceMember,
    //
    currentServiceMember: serviceMember,
    currentStation: get(formValues, 'current_station', {}),
    newDutyStation: get(orders, 'new_duty_station', {}),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(DutyStation);
