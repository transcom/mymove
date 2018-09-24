import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';
import { get, capitalize } from 'lodash';
import { NavLink } from 'react-router-dom';
import { reduxForm } from 'redux-form';

import Alert from 'shared/Alert';
import { withContext } from 'shared/AppContext';
import PremoveSurvey from 'shared/PremoveSurvey';
import ConfirmWithReasonButton from 'shared/ConfirmWithReasonButton';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import {
  loadShipmentDependencies,
  patchShipment,
  acceptShipment,
  generateGBL,
  rejectShipment,
  transportShipment,
  deliverShipment,
} from './ducks';
import ServiceAgents from './ServiceAgents';
import Weights from './Weights';
import FormButton from './FormButton';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';

const attachmentsErrorMessages = {
  400: 'There is already a GBL for this shipment. ',
  417: 'Missing data required to generate a Bill of Lading.',
};

class AcceptShipmentPanel extends Component {
  rejectShipment = reason => {
    this.props.rejectShipment(reason);
  };

  acceptShipment = () => {
    this.props.acceptShipment();
  };

  render() {
    return (
      <div>
        <button className="usa-button-primary" onClick={this.acceptShipment}>
          Accept Shipment
        </button>
        <ConfirmWithReasonButton
          buttonTitle="Reject Shipment"
          reasonPrompt="Why are you rejecting this shipment?"
          warningPrompt="Are you sure you want to reject this shipment?"
          onConfirm={this.rejectShipment}
          buttonDisabled={true}
        />
      </div>
    );
  }
}

let PickupDateForm = props => {
  const { schema, onCancel, handleSubmit, submitting, valid } = props;

  return (
    <form onSubmit={handleSubmit}>
      <SwaggerField fieldName="actual_pickup_date" swagger={schema} required />

      <button onClick={onCancel}>Cancel</button>
      <button type="submit" disabled={submitting || !valid}>
        Done
      </button>
    </form>
  );
};

PickupDateForm = reduxForm({ form: 'pickup_shipment' })(PickupDateForm);

let DeliveryDateForm = props => {
  const { schema, onCancel, handleSubmit, submitting, valid } = props;

  return (
    <form onSubmit={handleSubmit}>
      <SwaggerField
        fieldName="actual_delivery_date"
        swagger={schema}
        required
      />

      <button onClick={onCancel}>Cancel</button>
      <button type="submit" disabled={submitting || !valid}>
        Done
      </button>
    </form>
  );
};

DeliveryDateForm = reduxForm({ form: 'deliver_shipment' })(DeliveryDateForm);

class ShipmentInfo extends Component {
  state = {
    redirectToHome: false,
  };

  componentDidMount() {
    this.props.loadShipmentDependencies(this.props.match.params.shipmentId);
  }

  acceptShipment = () => {
    return this.props.acceptShipment(this.props.shipment.id);
  };

  generateGBL = () => {
    return this.props.generateGBL(this.props.shipment.id);
  };

  rejectShipment = reason => {
    return this.props
      .rejectShipment(this.props.shipment.id, reason)
      .then(() => {
        this.setState({ redirectToHome: true });
      });
  };

  pickupShipment = values =>
    this.props.transportShipment(this.props.shipment.id, values);

  deliverShipment = values =>
    this.props.deliverShipment(this.props.shipment.id, values);

  render() {
    const serviceMember = get(this.props.shipment, 'service_member', {});
    const move = get(this.props.shipment, 'move', {});
    const gbl = get(this.props.shipment, 'gbl_number');

    const awarded = this.props.shipment.status === 'AWARDED';
    const approved = this.props.shipment.status === 'APPROVED';
    const inTransit = this.props.shipment.status === 'IN_TRANSIT';

    if (this.state.redirectToHome) {
      return <Redirect to="/" />;
    }

    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds">
            MOVE INFO - {move.selected_move_type} CODE D
            <h1>
              Shipment Info: {serviceMember.last_name},{' '}
              {serviceMember.first_name}
            </h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            {awarded && (
              <NavLink to="/queues/new" activeClassName="usa-current">
                <span>New Shipments Queue</span>
              </NavLink>
            )}
            {!awarded && (
              <NavLink to="/queues/all" activeClassName="usa-current">
                <span>All Shipments Queue</span>
              </NavLink>
            )}
          </div>
        </div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-one-whole">
            <ul className="move-info-header-meta">
              <li>GBL# {gbl}</li>
              <li>Locator# {move.locator}</li>
              <li>
                {this.props.shipment.source_gbloc} to{' '}
                {this.props.shipment.destination_gbloc}
              </li>
              <li>DoD ID# {serviceMember.edipi}</li>
              <li>
                {serviceMember.telephone}
                {serviceMember.phone_is_preferred && (
                  <FontAwesomeIcon
                    className="icon"
                    icon={faPhone}
                    flip="horizontal"
                  />
                )}
                {serviceMember.text_message_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faComments} />
                )}
                {serviceMember.email_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faEmail} />
                )}
              </li>
              <li>
                Status: <b>{capitalize(this.props.shipment.status)}</b>
              </li>
            </ul>
          </div>
        </div>
        <div className="usa-grid grid-wide panels-body">
          <div className="usa-width-one-whole">
            <div className="usa-width-two-thirds">
              {this.props.loadTspDependenciesHasSuccess && (
                <div className="office-tab">
                  <PremoveSurvey
                    title="Premove Survey"
                    shipment={this.props.shipment}
                    update={this.props.patchShipment}
                  />
                  <ServiceAgents
                    title="ServiceAgents"
                    shipment={this.props.shipment}
                    serviceAgents={this.props.serviceAgents}
                  />
                  <Weights
                    title="Weights & Items"
                    shipment={this.props.shipment}
                    update={this.props.patchShipment}
                  />
                </div>
              )}
            </div>
            <div className="usa-width-one-third">
              {awarded && (
                <AcceptShipmentPanel
                  acceptShipment={this.acceptShipment}
                  rejectShipment={this.rejectShipment}
                  shipmentStatus={this.props.shipment.status}
                />
              )}
              {approved && (
                <FormButton
                  formComponent={PickupDateForm}
                  schema={this.props.pickupSchema}
                  onSubmit={this.pickupShipment}
                  buttonTitle="Enter Pickup"
                />
              )}
              {inTransit && (
                <FormButton
                  formComponent={DeliveryDateForm}
                  schema={this.props.deliverSchema}
                  onSubmit={this.deliverShipment}
                  buttonTitle="Enter Delivery"
                />
              )}
              {this.props.generateGBLError && (
                <Alert type="warning" heading="An error occurred">
                  {attachmentsErrorMessages[this.props.error.statusCode] ||
                    'Something went wrong contacting the server.'}
                </Alert>
              )}
              {this.props.generateGBLSuccess && (
                <Alert type="success" heading="Success!">
                  GBL generated successfully.
                </Alert>
              )}
              <div>
                <button onClick={this.generateGBL}>
                  Generate Bill of Lading
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
  shipment: get(state, 'tsp.shipment', {}),
  serviceAgents: get(state, 'tsp.serviceAgents', []),
  loadTspDependenciesHasSuccess: get(
    state,
    'tsp.loadTspDependenciesHasSuccess',
  ),
  loadTspDependenciesHasError: get(state, 'tsp.loadTspDependenciesHasError'),
  acceptError: get(state, 'tsp.shipmentHasAcceptError'),
  generateGBLError: get(state, 'tsp.generateGBLError'),
  generateGBLSuccess: get(state, 'tsp.generateGBLSuccess'),
  error: get(state, 'tsp.error'),
  pickupSchema: get(state, 'swagger.spec.definitions.ActualPickupDate', {}),
  deliverSchema: get(state, 'swagger.spec.definitions.ActualDeliveryDate', {}),
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      loadShipmentDependencies,
      patchShipment,
      acceptShipment,
      generateGBL,
      rejectShipment,
      transportShipment,
      deliverShipment,
    },
    dispatch,
  );

export default withContext(
  connect(mapStateToProps, mapDispatchToProps)(ShipmentInfo),
);
