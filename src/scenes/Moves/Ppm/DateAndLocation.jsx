import { debounce, get, bind, cloneDeep } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import {
  createOrUpdatePpm,
  getDestinationPostalCode,
  getPpmSitEstimate,
  isHHGPPMComboMove,
  setInitialFormValues,
} from './ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { loadEntitlementsFromState } from 'shared/entitlements';
import Alert from 'shared/Alert';
import WizardHeader from '../WizardHeader';
import ppmBlack from 'shared/icon/ppm-black.svg';
import './DateAndLocation.css';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';

const sitEstimateDebounceTime = 300;
const formName = 'ppp_date_and_location';
const DateAndLocationWizardForm = reduxifyWizardForm(formName);

export class DateAndLocation extends Component {
  state = { showInfo: false };

  componentDidMount() {
    if (!this.props.currentPpm && this.props.isHHGPPMComboMove) {
      const { originalMoveDate, pickupPostalCode, destinationPostalCode } = this.props.defaultValues;
      this.props.setInitialFormValues(originalMoveDate, pickupPostalCode, destinationPostalCode);
    }
  }

  openInfo = () => {
    this.setState({ showInfo: true });
  };
  closeInfo = () => {
    this.setState({ showInfo: false });
  };

  handleSubmit = () => {
    const pendingValues = Object.assign({}, this.props.formValues);
    if (pendingValues) {
      pendingValues.has_additional_postal_code = pendingValues.has_additional_postal_code || false;
      pendingValues.has_sit = pendingValues.has_sit || false;
      if (!pendingValues.has_sit) {
        pendingValues.days_in_storage = null;
      }
      const moveId = this.props.match.params.moveId;
      return this.props.createOrUpdatePpm(moveId, pendingValues);
    }
  };

  getSitEstimate = (moveDate, sitDays, pickupZip, destZip, weight) => {
    if (!pickupZip || !destZip) return;
    if (sitDays <= 90 && pickupZip.length === 5 && destZip.length === 5) {
      this.props.getPpmSitEstimate(moveDate, sitDays, pickupZip, destZip, weight);
    }
  };

  debouncedSitEstimate = debounce(bind(this.getSitEstimate, this), sitEstimateDebounceTime);

  getDebouncedSitEstimate = (e, value, _, field) => {
    const { formValues, entitlement } = this.props;
    const estimateValues = cloneDeep(formValues);
    estimateValues[field] = value; // eslint-disable-line security/detect-object-injection
    this.debouncedSitEstimate(
      estimateValues.original_move_date,
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
      error,
      currentOrders,
      initialValues,
      sitReimbursement,
      hasEstimateError,
      isHHGPPMComboMove,
    } = this.props;

    return (
      <div>
        {isHHGPPMComboMove && (
          <WizardHeader
            icon={ppmBlack}
            title="Move Setup"
            right={
              <ProgressTimeline>
                <ProgressTimelineStep name="Move Setup" current />
                <ProgressTimelineStep name="Review" />
              </ProgressTimeline>
            }
          />
        )}
        <DateAndLocationWizardForm
          handleSubmit={this.handleSubmit}
          pageList={pages}
          pageKey={pageKey}
          serverError={error}
          initialValues={initialValues}
          enableReinitialize={true} //this is needed as the pickup_postal_code value needs to be initialized to the users residential address
        >
          <h2>PPM Dates & Locations</h2>
          {isHHGPPMComboMove && <div>Great! Let's review your pickup and destination information.</div>}
          <h3> Move Date </h3>
          <SwaggerField
            fieldName="original_move_date"
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
          {!isHHGPPMComboMove && (
            <SwaggerField fieldName="has_additional_postal_code" swagger={this.props.schema} component={YesNoBoolean} />
          )}
          {get(this.props, 'formValues.has_additional_postal_code', false) && (
            <Fragment>
              <SwaggerField fieldName="additional_pickup_postal_code" swagger={this.props.schema} required />
              <span className="grey">
                Making additional stops may decrease your PPM incentive. <a onClick={this.openInfo}>Why</a>
              </span>
              {this.state.showInfo && (
                <Alert type="info" heading="">
                  Your PPM incentive is based primarily off two factors -- the weight of your household goods and the
                  base rate it would cost the government to transport your household goods between your destination and
                  origin. When you add additional stops, your overall PPM incentive will change to account for any
                  deviations from the standard route and to account for the fact that not 100% of your household goods
                  travelled the entire way from origin to destination. <a onClick={this.closeInfo}>Close</a>
                </Alert>
              )}
            </Fragment>
          )}
          <h3>Destination Location</h3>
          {!isHHGPPMComboMove && (
            <p>
              Enter the ZIP for your new home if you know it, or for{' '}
              {this.props.currentOrders && this.props.currentOrders.new_duty_station.name} if you don't.
            </p>
          )}
          <SwaggerField
            fieldName="destination_postal_code"
            swagger={this.props.schema}
            onChange={this.getDebouncedSitEstimate}
            required
          />
          <span className="grey">
            The ZIP code for {currentOrders && currentOrders.new_duty_station.name} is{' '}
            {currentOrders && currentOrders.new_duty_station.address.postal_code}{' '}
          </span>
          {!isHHGPPMComboMove && (
            <SwaggerField fieldName="has_sit" swagger={this.props.schema} component={YesNoBoolean} />
          )}
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
                  You can spend up to {sitReimbursement} on private storage. Save your receipts to submit with your PPM
                  paperwork.
                </div>
              )}
              {hasEstimateError && (
                <div className="usa-width-one-whole error-message">
                  <Alert type="warning" heading="Could not retrieve estimate">
                    There was an issue retrieving an estimate for how much you could be reimbursed for private storage.
                    You still qualify but may need to talk with your local PPPO.
                  </Alert>
                </div>
              )}
            </Fragment>
          )}
        </DateAndLocationWizardForm>
      </div>
    );
  }
}

DateAndLocation.propTypes = {
  schema: PropTypes.object.isRequired,
  createOrUpdatePpm: PropTypes.func.isRequired,
  error: PropTypes.object,
};

function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.UpdatePersonallyProcuredMovePayload', {}),
    ...state.ppm,
    currentOrders: state.orders.currentOrders,
    formValues: getFormValues(formName)(state),
    entitlement: loadEntitlementsFromState(state),
    hasEstimateError: state.ppm.hasEstimateError,
    isHHGPPMComboMove: isHHGPPMComboMove(state),
  };
  const defaultPickupZip = get(state.serviceMember, 'currentServiceMember.residential_address.postal_code');
  const currentOrders = state.orders.currentOrders;

  props.initialValues = props.currentPpm
    ? props.currentPpm
    : defaultPickupZip
      ? {
          pickup_postal_code: defaultPickupZip,
        }
      : null;

  if (props.isHHGPPMComboMove) {
    props.defaultValues = {
      pickupPostalCode: defaultPickupZip,
      originalMoveDate: currentOrders.issue_date,
      // defaults to SM's destination address, if none, uses destination duty station zip
      destinationPostalCode: getDestinationPostalCode(state),
    };
  }

  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createOrUpdatePpm, getPpmSitEstimate, setInitialFormValues }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DateAndLocation);
