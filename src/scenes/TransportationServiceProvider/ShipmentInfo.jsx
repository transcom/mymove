import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';
import { get, capitalize } from 'lodash';

import { NavLink } from 'react-router-dom';

import Alert from 'shared/Alert';
import { withContext } from 'shared/AppContext';

import {
  loadShipmentDependencies,
  patchShipment,
  acceptShipment,
  generateGBL,
  rejectShipment,
} from './ducks';
import PremoveSurvey from 'shared/PremoveSurvey';
import { formatDate } from 'shared/formatters';
import ConfirmWithReasonButton from 'shared/ConfirmWithReasonButton';
import ServiceAgents from './ServiceAgents';
import Weights from './Weights';

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
        />
      </div>
    );
  }
}

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

  render() {
    const last_name = get(this.props.shipment, 'service_member.last_name');
    const first_name = get(this.props.shipment, 'service_member.first_name');
    const locator = get(this.props.shipment, 'move.locator');
    const awarded = this.props.shipment.status === 'AWARDED';

    if (this.state.redirectToHome) {
      return <Redirect to="/" />;
    }

    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds">
            <h1>
              Shipment Info: {last_name}, {first_name}
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
              <li className="Todo-phase2">GBL# OHAI9999999</li>
              <li>Locator# {locator}</li>
              <li>
                {this.props.shipment.source_gbloc} to{' '}
                {this.props.shipment.destination_gbloc}
              </li>
              <li>
                Requested Move date{' '}
                {formatDate(this.props.shipment.requested_pickup_date)}
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
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      loadShipmentDependencies,
      patchShipment,
      acceptShipment,
      generateGBL,
      rejectShipment,
    },
    dispatch,
  );

export default withContext(
  connect(mapStateToProps, mapDispatchToProps)(ShipmentInfo),
);
