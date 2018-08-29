import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, capitalize } from 'lodash';

import { NavLink } from 'react-router-dom';

import { withContext } from 'shared/AppContext';

import {
  loadShipmentDependencies,
  patchShipment,
  acceptShipment,
} from './ducks';
import PremoveSurvey from 'shared/PremoveSurvey';
import { formatDate } from 'shared/formatters';
import ServiceAgents from './ServiceAgents';

class AcceptShipmentPanel extends Component {
  rejectShipment = () => {
    this.setState({ displayState: 'Rejected' });
    // TODO (rebecca): Add rejection flow
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
        <button className="usa-button-secondary" onClick={this.rejectShipment}>
          Reject Shipment
        </button>
      </div>
    );
  }
}

class ShipmentInfo extends Component {
  componentDidMount() {
    this.props.loadShipmentDependencies(this.props.match.params.shipmentId);
  }

  acceptShipment = () => {
    return this.props.acceptShipment(this.props.shipment.id);
  };

  render() {
    var last_name = get(this.props.shipment, 'service_member.last_name');
    var first_name = get(this.props.shipment, 'service_member.first_name');
    var locator = get(this.props.shipment, 'move.locator');

    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds">
            <h1>
              Shipment Info: {last_name}, {first_name}
            </h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Shipments Queue</span>
            </NavLink>
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
        <div className="usa-grid grid-wide tabs">
          {this.props.loadTspDependenciesHasSuccess && (
            <div className="office-tab usa-width-two-thirds">
              <PremoveSurvey
                title="Premove Survey"
                shipment={this.props.shipment}
                update={this.props.patchShipment}
              />
              <ServiceAgents
                title="Service Agents"
                shipment={this.props.shipment}
                serviceAgents={this.props.serviceAgents}
              />
            </div>
          )}
          <div className="usa-width-one-third">
            {this.props.shipment.status === 'AWARDED' && (
              <AcceptShipmentPanel
                acceptShipment={this.acceptShipment}
                shipmentStatus={this.props.shipment.status}
              />
            )}
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
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      loadShipmentDependencies,
      patchShipment,
      acceptShipment,
    },
    dispatch,
  );

export default withContext(
  connect(mapStateToProps, mapDispatchToProps)(ShipmentInfo),
);
