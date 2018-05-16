import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { loadPpm, createOrUpdatePpm } from './ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './DateAndLocation.css';

const NULL_ZIP = ''; //HACK: until we can figure out how to unset zip
const formName = 'ppp_date_and_location';
const DateAndLocationWizardForm = reduxifyWizardForm(formName);

export class DateAndLocation extends Component {
  state = { showAdditionalPickup: false, showTempStorage: false };
  static getDerivedStateFromProps(nextProps, prevState) {
    const result = {};
    if (
      get(
        nextProps,
        'formData.values.additional_pickup_postal_code',
        NULL_ZIP,
      ) !== NULL_ZIP
    )
      result.showAdditionalPickup = true;
    if (get(nextProps, 'formData.values.days_in_storage', 0) > 0)
      result.showTempStorage = true;
    return result;
  }
  setShowAdditionalPickup = show => {
    this.setState({ showAdditionalPickup: show }, () => {
      if (!show) console.log(show);
      // Redux form isn't being connected to reducer, can't access change prop
      // this.props.change('additional_pickup_postal_code', NULL_ZIP);
    });
  };
  setShowTempStorage = show => {
    this.setState({ showTempStorage: show }, () => {
      if (!show) console.log(show);
      // if (!show) this.props.change('days_in_storage', '0');
    });
  };
  componentDidMount() {
    document.title = 'Transcom PPP: Date & Locations';
    const moveId = this.props.match.params.moveId;
    this.props.loadPpm(moveId);
  }
  handleSubmit = () => {
    const pendingValues = Object.assign({}, this.props.formData.values);
    if (pendingValues) {
      const moveId = this.props.match.params.moveId;
      pendingValues[
        'has_additional_postal_code'
      ] = this.state.showAdditionalPickup;
      pendingValues['has_sit'] = this.state.showTempStorage;
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
    } = this.props;
    const initialValues = currentPpm ? currentPpm : null;
    const { showAdditionalPickup, showTempStorage } = this.state;
    return (
      <DateAndLocationWizardForm
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
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
        <p>Do you have stuff at another pickup location?</p>
        <YesNoBoolean
          value={showAdditionalPickup}
          onChange={this.setShowAdditionalPickup}
        />
        {this.state.showAdditionalPickup && (
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
        <p>
          Are you going to put your stuff in temporary storage before moving
          into your new home?
        </p>
        <YesNoBoolean
          value={showTempStorage}
          onChange={this.setShowTempStorage}
        />
        {this.state.showTempStorage && (
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
    formData: state.form[formName],
    enableReinitialize: true,
  };
  return props;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadPpm, createOrUpdatePpm }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DateAndLocation);
