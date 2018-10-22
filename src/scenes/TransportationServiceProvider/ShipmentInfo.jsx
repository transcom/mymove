import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';
import { get, capitalize } from 'lodash';
import { NavLink, Link } from 'react-router-dom';
import { reduxForm } from 'redux-form';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';

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
import {
  getAllTariff400ngItems,
  selectTariff400ngItems,
  getTariff400ngItemsLabel,
} from 'shared/Entities/modules/tariff400ngItems';
import {
  getAllShipmentAccessorials,
  selectShipmentAccessorials,
  getShipmentAccessorialsLabel,
} from 'shared/Entities/modules/shipmentAccessorials';

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
  packShipment,
  deliverShipment,
} from './ducks';
import ServiceAgents from './ServiceAgents';
import Weights from './Weights';
import Dates from './Dates';
import LocationsContainer from './LocationsContainer';
import FormButton from './FormButton';
import CustomerInfo from './CustomerInfo';
import PreApprovalPanel from 'shared/PreApprovalRequest/PreApprovalPanel.jsx';

import './tsp.css';

const attachmentsErrorMessages = {
  400: 'There is already a GBL for this shipment. ',
  417: 'Missing data required to generate a GBL.',
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

let PackDateForm = props => {
  const { schema, onCancel, handleSubmit, submitting, valid } = props;

  return (
    <form onSubmit={handleSubmit}>
      <SwaggerField fieldName="actual_pack_date" swagger={schema} required />

      <button onClick={onCancel}>Cancel</button>
      <button type="submit" disabled={submitting || !valid}>
        Done
      </button>
    </form>
  );
};

PackDateForm = reduxForm({ form: 'pack_date_shipment' })(PackDateForm);

let DeliveryDateForm = props => {
  const { schema, onCancel, handleSubmit, submitting, valid } = props;

  return (
    <form onSubmit={handleSubmit}>
      <SwaggerField fieldName="actual_delivery_date" swagger={schema} required />

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
    this.props.loadShipmentDependencies(this.props.match.params.shipmentId).catch(err => {
      this.props.history.replace('/');
    });
  }

  componentDidUpdate(prevProps, prevState) {
    if (!prevProps.shipment.id && this.props.shipment.id) {
      this.props.getAllShipmentDocuments(getShipmentDocumentsLabel, this.props.shipment.id);
      this.props.getAllTariff400ngItems(true, getTariff400ngItemsLabel);
      this.props.getAllShipmentAccessorials(getShipmentAccessorialsLabel, this.props.shipment.id);
    }
  }

  acceptShipment = () => {
    return this.props.acceptShipment(this.props.shipment.id);
  };

  generateGBL = () => {
    return this.props.generateGBL(this.props.shipment.id);
  };

  rejectShipment = reason => {
    return this.props.rejectShipment(this.props.shipment.id, reason).then(() => {
      this.setState({ redirectToHome: true });
    });
  };

  pickupShipment = values => this.props.transportShipment(this.props.shipment.id, values);

  packShipment = values => this.props.packShipment(this.props.shipment.id, values);

  deliverShipment = values => this.props.deliverShipment(this.props.shipment.id, values);

  render() {
    const { context, shipment, shipmentDocuments } = this.props;
    const {
      service_member: serviceMember = {},
      move = {},
      gbl_number: gbl,
      actual_pack_date,
      actual_pickup_date,
    } = shipment;

    const shipmentId = this.props.match.params.shipmentId;
    const newDocumentUrl = `/shipments/${shipmentId}/documents/new`;
    const showDocumentViewer = context.flags.documentViewer;
    const awarded = shipment.status === 'AWARDED';
    const approved = shipment.status === 'APPROVED';
    const inTransit = shipment.status === 'IN_TRANSIT';
    const pmSurveyComplete = Boolean(
      shipment.pm_survey_conducted_date &&
        shipment.pm_survey_method &&
        shipment.pm_survey_planned_pack_date &&
        shipment.pm_survey_planned_pickup_date &&
        shipment.pm_survey_planned_delivery_date &&
        shipment.pm_survey_weight_estimate,
    );
    const gblGenerated =
      shipmentDocuments && shipmentDocuments.find(element => element.move_document_type === 'GOV_BILL_OF_LADING');

    if (this.state.redirectToHome) {
      return <Redirect to="/" />;
    }

    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds">
            MOVE INFO - {move.selected_move_type} CODE D
            <h1>
              Shipment Info: {serviceMember.last_name}, {serviceMember.first_name}
            </h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            {awarded && (
              <NavLink to="/queues/new" activeClassName="usa-current">
                <span>New Shipments Queue</span>
              </NavLink>
            )}
            {approved && (
              <NavLink to="/queues/approved" activeClassName="usa-current">
                <span>Approved Shipments Queue</span>
              </NavLink>
            )}
            {!awarded &&
              !approved && (
                <NavLink to="/queues/all" activeClassName="usa-current">
                  <span>All Shipments Queue</span>
                </NavLink>
              )}
          </div>
        </div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-one-whole">
            <ul className="move-info-header-meta">
              <li>
                GBL# {gbl}
                &nbsp;
              </li>
              <li>
                Locator# {move.locator}
                &nbsp;
              </li>
              <li>
                {this.props.shipment.source_gbloc} to {this.props.shipment.destination_gbloc}
                &nbsp;
              </li>
              <li>
                DoD ID# {serviceMember.edipi}
                &nbsp;
              </li>
              <li>
                {serviceMember.telephone}
                {serviceMember.phone_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faPhone} flip="horizontal" />
                )}
                {serviceMember.text_message_is_preferred && <FontAwesomeIcon className="icon" icon={faComments} />}
                {serviceMember.email_is_preferred && <FontAwesomeIcon className="icon" icon={faEmail} />}
                &nbsp;
              </li>
              <li>
                Status: <b>{capitalize(this.props.shipment.status)}</b>
                &nbsp;
              </li>
            </ul>
          </div>
        </div>
        <div className="usa-grid grid-wide panels-body">
          <div className="usa-width-one-whole">
            <div className="usa-width-two-thirds">
              {awarded && (
                <AcceptShipmentPanel
                  acceptShipment={this.acceptShipment}
                  rejectShipment={this.rejectShipment}
                  shipmentStatus={this.props.shipment.status}
                />
              )}
              {this.props.generateGBLError && (
                <p>
                  <Alert type="warning" heading="An error occurred">
                    {attachmentsErrorMessages[this.props.error.statusCode] ||
                      'Something went wrong contacting the server.'}
                  </Alert>
                </p>
              )}
              {this.props.generateGBLSuccess && (
                <p>
                  <Alert type="success" heading="GBL has been created">
                    <span className="usa-grid usa-alert-no-padding">
                      <span className="usa-width-two-thirds">
                        Click the button to view, print, or download the GBL.
                      </span>
                      <span className="usa-width-one-third">
                        <Link to={`${this.props.gblDocUrl}`} className="usa-alert-right" target="_blank">
                          <button>View GBL</button>
                        </Link>
                      </span>
                    </span>
                  </Alert>
                </p>
              )}
              {approved &&
                pmSurveyComplete &&
                !gblGenerated && (
                  <div>
                    <button onClick={this.generateGBL} disabled={this.props.generateGBLInProgress}>
                      Generate the GBL
                    </button>
                  </div>
                )}
              {this.props.loadTspDependenciesHasSuccess && (
                <div className="office-tab">
                  <Dates title="Dates" shipment={this.props.shipment} update={this.props.patchShipment} />
                  <PremoveSurvey
                    title="Premove Survey"
                    shipment={this.props.shipment}
                    update={this.props.patchShipment}
                  />
                  <PreApprovalPanel shipmentId={this.props.match.params.shipmentId} />
                  <ServiceAgents
                    title="ServiceAgents"
                    shipment={this.props.shipment}
                    serviceAgents={this.props.serviceAgents}
                  />
                  <Weights title="Weights & Items" shipment={this.props.shipment} update={this.props.patchShipment} />
                  <LocationsContainer update={this.props.patchShipment} />
                </div>
              )}
            </div>
            <div className="usa-width-one-third">
              {approved &&
                !actual_pack_date && (
                  <FormButton
                    FormComponent={PackDateForm}
                    schema={this.props.packSchema}
                    onSubmit={this.packShipment}
                    buttonTitle="Enter Packing"
                  />
                )}
              {approved &&
                actual_pack_date &&
                !actual_pickup_date && (
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
              <div className="customer-info">
                <h2 className="extras usa-heading">Customer Info</h2>
                <CustomerInfo />
              </div>
              <div className="documents">
                <h2 className="extras usa-heading">
                  Documents
                  {!showDocumentViewer && <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />}
                  {showDocumentViewer && (
                    <Link to={newDocumentUrl} target="_blank">
                      <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
                    </Link>
                  )}
                </h2>
                {showDocumentViewer && shipmentDocuments.length ? (
                  <DocumentList
                    detailUrlPrefix={`/shipments/${shipmentId}/documents`}
                    moveDocuments={shipmentDocuments}
                  />
                ) : (
                  <Link className="status" to={newDocumentUrl} target="_blank">
                    <span>
                      <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} />
                    </span>
                    Upload new document
                  </Link>
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

  return {
    swaggerError: state.swaggerPublic.hasErrored,
    shipment,
    shipmentDocuments: selectShipmentDocuments(state, shipment.id),
    tariff400ngItems: selectTariff400ngItems(state),
    shipmentAccessorials: selectShipmentAccessorials(state),
    serviceAgents: get(state, 'tsp.serviceAgents', []),
    loadTspDependenciesHasSuccess: get(state, 'tsp.loadTspDependenciesHasSuccess'),
    loadTspDependenciesHasError: get(state, 'tsp.loadTspDependenciesHasError'),
    acceptError: get(state, 'tsp.shipmentHasAcceptError'),
    generateGBLError: get(state, 'tsp.generateGBLError'),
    generateGBLSuccess: get(state, 'tsp.generateGBLSuccess'),
    generateGBLInProgress: get(state, 'tsp.generateGBLInProgress'),
    gblDocUrl: get(state, 'tsp.gblDocUrl'),
    error: get(state, 'tsp.error'),
    pickupSchema: get(state, 'swaggerPublic.spec.definitions.ActualPickupDate', {}),
    packSchema: get(state, 'swaggerPublic.spec.definitions.ActualPackDate', {}),
    deliverSchema: get(state, 'swaggerPublic.spec.definitions.ActualDeliveryDate', {}),
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
      packShipment,
      deliverShipment,
      getAllShipmentDocuments,
      getAllTariff400ngItems,
      getAllShipmentAccessorials,
    },
    dispatch,
  );

export default withContext(connect(mapStateToProps, mapDispatchToProps)(ShipmentInfo));
