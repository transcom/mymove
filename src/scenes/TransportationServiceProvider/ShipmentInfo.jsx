import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';
import { get, capitalize } from 'lodash';
import { NavLink, Link } from 'react-router-dom';
import { reduxForm } from 'redux-form';

import Alert from 'shared/Alert';
import DocumentList from 'shared/DocumentViewer/DocumentList';
import { withContext } from 'shared/AppContext';
import PremoveSurvey from 'shared/PremoveSurvey';
import ConfirmWithReasonButton from 'shared/ConfirmWithReasonButton';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import {
  getAllShipmentDocuments,
  selectShipmentDocuments,
  getShipmentDocumentsLabel,
} from 'shared/Entities/modules/shipmentDocuments';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';
import faExternalLinkAlt from '@fortawesome/fontawesome-free-solid/faExternalLinkAlt';

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
import Dates from './Dates';
import Locations from './Locations';
import FormButton from './FormButton';
import CustomerInfo from './CustomerInfo';

import './tsp.css';

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
    this.props.getAllShipmentDocuments(
      getShipmentDocumentsLabel,
      this.props.match.params.shipmentId,
    );
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
    const {
      context,
      shipment,
      shipmentDocuments,
      deliveryAddress,
    } = this.props;

    const {
      service_member: serviceMember = {},
      move = {},
      gbl_number: gbl,
    } = shipment;

    const showDocumentViewer = context.flags.documentViewer;
    const awarded = shipment.status === 'AWARDED';
    const approved = shipment.status === 'APPROVED';
    const inTransit = shipment.status === 'IN_TRANSIT';

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
                  <Dates
                    title="Dates"
                    shipment={this.props.shipment}
                    update={this.props.patchShipment}
                  />
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
                  <Locations
                    deliveryAddress={deliveryAddress}
                    title="Locations"
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
                  FormComponent={PickupDateForm}
                  schema={this.props.pickupSchema}
                  onSubmit={this.pickupShipment}
                  buttonTitle="Enter Pickup"
                />
              )}
              {inTransit && (
                <FormButton
                  FormComponent={DeliveryDateForm}
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
              <div className="customer-info">
                <h2 className="extras usa-heading">Customer Info</h2>
                <CustomerInfo />
              </div>
              <div className="documents">
                <h2 className="documents-list-header">
                  Documents
                  {!showDocumentViewer && (
                    <FontAwesomeIcon
                      className="icon"
                      icon={faExternalLinkAlt}
                    />
                  )}
                  {showDocumentViewer && (
                    <Link to={`/moves/${move.id}/documents`} target="_blank">
                      <FontAwesomeIcon
                        className="icon"
                        icon={faExternalLinkAlt}
                      />
                    </Link>
                  )}
                </h2>
                {showDocumentViewer && shipmentDocuments.length ? (
                  <DocumentList
                    detailUrlPrefix={`/moves/${
                      this.props.match.params.shipmentId
                    }/documents`}
                    moveDocuments={shipmentDocuments}
                  />
                ) : (
                  <p>No orders have been uploaded.</p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  const shipment = get(state, 'tsp.shipment', {});
  const newDutyStation = get(shipment, 'move.new_duty_station.address', {});
  // if they do not have a delivery address, default to the station's address info
  const deliveryAddress = shipment.has_delivery_address
    ? shipment.delivery_address
    : newDutyStation;

  return {
    swaggerError: state.swagger.hasErrored,
    shipment,
    deliveryAddress,
    shipmentDocuments: selectShipmentDocuments(state),
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
    deliverSchema: get(
      state,
      'swagger.spec.definitions.ActualDeliveryDate',
      {},
    ),
  };
};

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
      getAllShipmentDocuments,
    },
    dispatch,
  );

export default withContext(
  connect(mapStateToProps, mapDispatchToProps)(ShipmentInfo),
);
