import { get, isEmpty } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { loadEntitlementsFromState } from 'shared/entitlements';
import {
  loadPPMs,
  createPPM,
  selectActivePPMForMove,
  updatePPM,
  updatePPMEstimate,
} from 'shared/Entities/modules/ppms';
import { fetchLatestOrders } from 'shared/Entities/modules/orders';
import Alert from 'shared/Alert';
import { ValidateZipRateData } from 'shared/api';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { selectServiceMemberFromLoggedInUser, selectCurrentOrders, selectCurrentMove } from 'store/entities/selectors';

import './DateAndLocation.css';

const formName = 'ppp_date_and_location';

const UnsupportedZipCodeErrorMsg =
  'Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.';

async function asyncValidate(values, dispatch, props, currentFieldName) {
  const { pickup_postal_code, destination_postal_code } = values;

  // If either postal code is blurred, check both of them for errors. We want to
  // catch these before checking on dates via `GetPpmWeightEstimate`.
  if (['destination_postal_code', 'pickup_postal_code'].includes(currentFieldName)) {
    // eslint-disable-next-line security/detect-object-injection
    const zipValue = values[currentFieldName];
    if (zipValue && zipValue.length < 5) {
      return;
    }
    const pickupZip = pickup_postal_code && pickup_postal_code.slice(0, 5);
    const destinationZip = destination_postal_code && destination_postal_code.slice(0, 5);
    const responseObject = {};
    if (pickupZip) {
      const responseBody = await ValidateZipRateData(pickupZip, 'origin');
      if (!responseBody.valid) {
        responseObject.pickup_postal_code = UnsupportedZipCodeErrorMsg;
      }
    }
    if (destinationZip) {
      const responseBody = await ValidateZipRateData(destinationZip, 'destination');
      if (!responseBody.valid) {
        responseObject.destination_postal_code = UnsupportedZipCodeErrorMsg;
      }
    }
    if (responseObject.pickup_postal_code || responseObject.destination_postal_code) {
      throw responseObject;
    }
  }
}

const DateAndLocationWizardForm = reduxifyWizardForm(formName, null, asyncValidate, [
  'destination_postal_code',
  'pickup_postal_code',
  'original_move_date',
]);

const validateDifferentZip = (value, formValues) => {
  if (value && value === formValues.pickup_postal_code) {
    return 'You entered the same zip code for your origin and destination. Please change one of them.';
  }
};

export class DateAndLocation extends Component {
  componentDidMount() {
    const moveId = this.props.match.params.moveId;
    this.props.loadPPMs(moveId);
    this.props.fetchLatestOrders(this.props.serviceMemberId);
  }

  state = { showInfo: false };

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
      if (isEmpty(this.props.currentPPM)) {
        return this.props
          .createPPM(moveId, pendingValues)
          .then(({ response }) => this.props.updatePPMEstimate(moveId, response.body.id).catch((err) => err));
      } else {
        return this.props
          .updatePPM(moveId, this.props.currentPPM.id, pendingValues)
          .then(({ response }) => this.props.updatePPMEstimate(moveId, response.body.id).catch((err) => err));
      }
    }
  };

  render() {
    const { pages, pageKey, error, currentOrders, initialValues } = this.props;

    return (
      <div>
        <DateAndLocationWizardForm
          reduxFormSubmit={this.handleSubmit}
          pageList={pages}
          pageKey={pageKey}
          serverError={error}
          initialValues={initialValues}
          enableReinitialize={true} //this is needed as the pickup_postal_code value needs to be initialized to the users residential address
        >
          <h1 data-testid="location-page-title">PPM dates & locations</h1>
          <SectionWrapper>
            <h2> Move date </h2>
            <SwaggerField fieldName="original_move_date" swagger={this.props.schema} required />
          </SectionWrapper>
          <SectionWrapper>
            <h2>Pickup location</h2>
            <SwaggerField fieldName="pickup_postal_code" swagger={this.props.schema} required />
            <SwaggerField fieldName="has_additional_postal_code" swagger={this.props.schema} component={YesNoBoolean} />
            {get(this.props, 'formValues.has_additional_postal_code', false) && (
              <Fragment>
                <SwaggerField fieldName="additional_pickup_postal_code" swagger={this.props.schema} required />
                <span className="grey">
                  Making additional stops may decrease your PPM incentive.{' '}
                  <a onClick={this.openInfo} className="usa-link">
                    Why
                  </a>
                </span>
                {this.state.showInfo && (
                  <Alert type="info" heading="">
                    Your PPM incentive is based primarily off two factors -- the weight of your household goods and the
                    base rate it would cost the government to transport your household goods between your destination
                    and origin. When you add additional stops, your overall PPM incentive will change to account for any
                    deviations from the standard route and to account for the fact that not 100% of your household goods
                    travelled the entire way from origin to destination.{' '}
                    <a onClick={this.closeInfo} className="usa-link">
                      Close
                    </a>
                  </Alert>
                )}
              </Fragment>
            )}
          </SectionWrapper>
          <SectionWrapper>
            <h2>Destination location</h2>
            <p>
              Enter the ZIP for your new home if you know it, or for{' '}
              {this.props.currentOrders && this.props.currentOrders.new_duty_station.name} if you don't.
            </p>
            <SwaggerField
              fieldName="destination_postal_code"
              swagger={this.props.schema}
              validate={validateDifferentZip}
              required
            />
            <div style={{ marginTop: '0.5rem' }}>
              <span className="grey">
                The ZIP code for {currentOrders && currentOrders.new_duty_station.name} is{' '}
                {currentOrders && currentOrders.new_duty_station.address.postal_code}.
              </span>
            </div>
            <SwaggerField fieldName="has_sit" swagger={this.props.schema} component={YesNoBoolean} />
            {get(this.props, 'formValues.has_sit', false) && (
              <Fragment>
                <SwaggerField
                  className="days-in-storage"
                  fieldName="days_in_storage"
                  swagger={this.props.schema}
                  required
                />{' '}
                <span className="grey">You can choose up to 90 days.</span>
              </Fragment>
            )}
          </SectionWrapper>
        </DateAndLocationWizardForm>
      </div>
    );
  }
}

DateAndLocation.propTypes = {
  schema: PropTypes.object.isRequired,
  createPPM: PropTypes.func.isRequired,
  updatePPM: PropTypes.func.isRequired,
  error: PropTypes.object,
};

function mapStateToProps(state) {
  const currentMove = selectCurrentMove(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  const defaultPickupZip = serviceMember?.residential_address?.postal_code;
  const originDutyStationZip = serviceMember?.current_station?.address?.postal_code;
  const serviceMemberId = serviceMember?.id;

  const props = {
    serviceMemberId,
    schema: get(state, 'swaggerInternal.spec.definitions.UpdatePersonallyProcuredMovePayload', {}),
    currentPPM: selectActivePPMForMove(state, currentMove?.id),
    currentOrders: selectCurrentOrders(state),
    formValues: getFormValues(formName)(state),
    entitlement: loadEntitlementsFromState(state),
    originDutyStationZip: serviceMember?.current_station?.address?.postal_code,
  };

  props.initialValues = !isEmpty(props.currentPPM)
    ? props.currentPPM
    : defaultPickupZip
    ? {
        pickup_postal_code: defaultPickupZip,
        origin_duty_station_zip: originDutyStationZip,
      }
    : null;

  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadPPMs, createPPM, updatePPM, updatePPMEstimate, fetchLatestOrders }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DateAndLocation);
