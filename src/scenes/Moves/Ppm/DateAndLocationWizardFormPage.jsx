import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { loadPpm, createOrUpdatePpm } from './ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './DateAndLocation.css';

const formName = 'ppp_date_and_location';
const DateAndLocationWizardForm = reduxifyWizardForm(formName);

export class DateAndLocation extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Date & Locations';
    const moveId = this.props.match.params.moveId;
    this.props.loadPpm(moveId);
  }
  handleSubmit = () => {
    const pendingValues = Object.assign({}, this.props.formValues);
    if (pendingValues) {
      pendingValues['has_additional_postal_code'] =
        pendingValues.has_additional_postal_code || false;
      pendingValues['has_sit'] = pendingValues.has_sit || false;
      const moveId = this.props.match.params.moveId;
      this.props.createOrUpdatePpm(moveId, pendingValues);
    }
  };
  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentOrders,
      currentPpm,
      initialValues,
    } = this.props;
    return (
      <DateAndLocationWizardForm
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
        enableReinitialize={true} //this is needed as the pickup_postal_code value needs to be initialized to the users residential address
      >
        <h1 className="sm-heading">PPM Dates & Locations</h1>
        <h3> Move Date </h3>
        <SwaggerField
          fieldName="planned_move_date"
          swagger={this.props.schema}
          required
        />
        <h3>Pickup Location</h3>
        <SwaggerField
          fieldName="pickup_postal_code"
          swagger={this.props.schema}
          required
        />
        <SwaggerField
          fieldName="has_additional_postal_code"
          swagger={this.props.schema}
          component={YesNoBoolean}
        />
        {get(this.props, 'formValues.has_additional_postal_code', false) && (
          <Fragment>
            <SwaggerField
              fieldName="additional_pickup_postal_code"
              swagger={this.props.schema}
              required
            />
            <p>Making additional stops may decrease your PPM incentive.</p>
          </Fragment>
        )}
        <h3>Destination Location</h3>
        <p>
          Enter the ZIP for your new home if you know it, or for{' '}
          {this.props.currentOrders &&
            this.props.currentOrders.new_duty_station.name}{' '}
          if you don't.
        </p>
        <SwaggerField
          fieldName="destination_postal_code"
          swagger={this.props.schema}
          required
        />
        The ZIP code for {currentOrders && currentOrders.new_duty_station.name}{' '}
        is {currentOrders && currentOrders.new_duty_station.address.postal_code}{' '}
        <SwaggerField
          fieldName="has_sit"
          swagger={this.props.schema}
          component={YesNoBoolean}
        />
        {get(this.props, 'formValues.has_sit', false) && (
          <Fragment>
            <SwaggerField
              fieldName="days_in_storage"
              swagger={this.props.schema}
              required
            />{' '}
            <p>You can choose up to 90 days.</p>
          </Fragment>
        )}
      </DateAndLocationWizardForm>
    );
  }
}

DateAndLocation.propTypes = {
  schema: PropTypes.object.isRequired,
  loadPpm: PropTypes.func.isRequired,
  createOrUpdatePpm: PropTypes.func.isRequired,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  const props = {
    schema: get(
      state,
      'swagger.spec.definitions.UpdatePersonallyProcuredMovePayload',
      {},
    ),
    ...state.ppm,
    ...state.loggedInUser,
    currentOrders:
      get(state.loggedInUser, 'loggedInUser.service_member.orders[0]') ||
      get(state.orders, 'currentOrders'),
    currentPpm: get(state.ppm, 'currentPpm'),
    formValues: getFormValues(formName)(state),
  };
  const defaultPickupZip = get(
    state.loggedInUser,
    'loggedInUser.service_member.residential_address.postal_code',
  );
  props.initialValues = props.currentPpm
    ? props.currentPpm
    : defaultPickupZip
      ? {
          pickup_postal_code: defaultPickupZip,
        }
      : null;
  return props;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadPpm, createOrUpdatePpm }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DateAndLocation);
