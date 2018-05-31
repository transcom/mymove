import { debounce, get, bind, cloneDeep } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { createOrUpdatePpm, getPpmSitEstimate } from './ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { loadEntitlements } from 'scenes/Orders/ducks';
import Alert from 'shared/Alert';

import './DateAndLocation.css';

const sitEstimateDebounceTime = 300;
const formName = 'ppp_date_and_location';
const DateAndLocationWizardForm = reduxifyWizardForm(formName);

export class DateAndLocation extends Component {
  state = { showInfo: false };

  openInfo = () => {
    this.setState({ showInfo: true });
  };
  closeInfo = () => {
    this.setState({ showInfo: false });
  };

  handleSubmit = () => {
    const { sitReimbursement } = this.props;
    const pendingValues = Object.assign({}, this.props.formValues);
    if (pendingValues) {
      pendingValues.has_additional_postal_code =
        pendingValues.has_additional_postal_code || false;
      pendingValues.has_sit = pendingValues.has_sit || false;
      if (pendingValues.has_sit) {
        pendingValues.estimated_storage_reimbursement = sitReimbursement;
      } else {
        pendingValues.days_in_storage = null;
        pendingValues.estimated_storage_reimbursement = null;
      }
      const moveId = this.props.match.params.moveId;
      this.props.createOrUpdatePpm(moveId, pendingValues);
    }
  };

  getSitEstimate = (moveDate, sitDays, pickupZip, destZip, weight) => {
    if (sitDays <= 90 && pickupZip.length === 5 && destZip.length === 5) {
      this.props.getPpmSitEstimate(
        moveDate,
        sitDays,
        pickupZip,
        destZip,
        weight,
      );
    }
  };

  debouncedSitEstimate = debounce(
    bind(this.getSitEstimate, this),
    sitEstimateDebounceTime,
  );

  getDebouncedSitEstimate = (e, value, _, field) => {
    const { formValues, entitlement } = this.props;
    const estimateValues = cloneDeep(formValues);
    estimateValues[field] = value;
    this.debouncedSitEstimate(
      estimateValues.planned_move_date,
      estimateValues.days_in_storage,
      estimateValues.pickup_postal_code,
      estimateValues.destination_postal_code,
      entitlement.sum,
    );
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentOrders,
      initialValues,
      sitReimbursement,
      hasEstimateError,
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
        <h2 className="sm-heading">PPM Dates & Locations</h2>
        <h3> Move Date </h3>
        <SwaggerField
          fieldName="planned_move_date"
          onChange={this.getDebouncedSitEstimate}
          swagger={this.props.schema}
          required
        />
        <h3>Pickup Location</h3>
        <SwaggerField
          fieldName="pickup_postal_code"
          onChange={this.getDebouncedSitEstimate}
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
            <span className="grey">
              Making additional stops may decrease your PPM incentive.{' '}
              <a onClick={this.openInfo}>Why</a>
            </span>
            {this.state.showInfo && (
              <Alert type="info" heading="">
                Your PPM incentive is based primarily off two factors -- the
                weight of your household goods and the base rate it would cost
                the government to transport your household goods between your
                destination and origin. When you add additional stops, your
                overall PPM incentive will change to account for any deviations
                from the standard route and to account for the fact that not
                100% of your household goods travelled the entire way from
                origin to destination. <a onClick={this.closeInfo}>Close</a>
              </Alert>
            )}
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
          onChange={this.getDebouncedSitEstimate}
          required
        />
        <span className="grey">
          The ZIP code for{' '}
          {currentOrders && currentOrders.new_duty_station.name} is{' '}
          {currentOrders && currentOrders.new_duty_station.address.postal_code}{' '}
        </span>
        <SwaggerField
          fieldName="has_sit"
          swagger={this.props.schema}
          component={YesNoBoolean}
        />
        {get(this.props, 'formValues.has_sit', false) && (
          <Fragment>
            <SwaggerField
              className="days-in-storage"
              fieldName="days_in_storage"
              swagger={this.props.schema}
              onChange={this.getDebouncedSitEstimate}
              required
            />{' '}
            <span className="grey">You can choose up to 90 days.</span>
            {sitReimbursement && (
              <div className="storage-estimate">
                You can spend up to {sitReimbursement} on private storage. Save
                your receipts to submit with your PPM paperwork.
              </div>
            )}
            {hasEstimateError && (
              <div className="usa-width-one-whole error-message">
                <Alert type="warning" heading="Could not retrieve estimate">
                  There was an issue retrieving an estimate for how much you
                  could be reimbursed for private storage. You still qualify but
                  may need to talk with your local PPPO.
                </Alert>
              </div>
            )}
          </Fragment>
        )}
      </DateAndLocationWizardForm>
    );
  }
}

DateAndLocation.propTypes = {
  schema: PropTypes.object.isRequired,
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
    currentOrders: state.orders.currentOrders,
    formValues: getFormValues(formName)(state),
    entitlement: loadEntitlements(state),
    hasEstimateError: state.ppm.hasEstimateError,
  };
  const defaultPickupZip = get(
    state.serviceMember,
    'currentServiceMember.residential_address.postal_code',
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
  return bindActionCreators({ createOrUpdatePpm, getPpmSitEstimate }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DateAndLocation);
