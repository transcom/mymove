import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, capitalize } from 'lodash';

import { NavLink } from 'react-router-dom';

import { withContext } from 'shared/AppContext';
import Alert from 'shared/Alert'; // eslint-disable-line

import { AcceptShipment } from './api.js';
import { loadShipmentDependencies } from './ducks';
import { formatDate } from 'shared/formatters';

class AcceptShipmentPanel extends Component {
  state = { displayState: 'Awarded', acceptError: false };

  rejectShipment = () => {
    this.setState({ displayState: 'Rejected' });
    // TODO (rebecca): Add rejection flow
  };

  acceptShipment = () => {
    AcceptShipment(this.props.shipmentId)
      .then(shipment => {
        this.setState({ displayState: 'Accepted', acceptError: false });
      })
      .catch(err => {
        console.log(err);
        this.setState({ displayState: 'Awarded', acceptError: true });
      });
  };

  render() {
    if (this.state.displayState === 'Awarded') {
      return (
        <div>
          {this.state.acceptError ? (
            <Alert type="error" heading="Unable to accept shipment" />
          ) : null}
          <button className="usa-button-primary" onClick={this.acceptShipment}>
            Accept Shipment
          </button>
          <button
            className="usa-button-secondary"
            onClick={this.rejectShipment}
          >
            Reject Shipment
          </button>
        </div>
      );
    } else if (this.state.displayState === 'Accepted') {
      return (
        <div>
          <Alert type="info" heading="Shipment accepted" />
        </div>
      );
    } else if (this.state.displayState === 'Rejected') {
      return (
        <div>
          <Alert type="error" heading="Shipment rejected" />
        </div>
      );
    }
  }
}

class ShipmentInfo extends Component {
  componentDidMount() {
    this.props.loadShipmentDependencies(this.props.match.params.shipmentId);
  }

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
          <div className="usa-width-two-thirds">
            {this.props.shipment.status === 'AWARDED' ? (
              <AcceptShipmentPanel shipmentId={this.props.shipment.id} />
            ) : null}
          </div>
          <div className="usa-width-one-third" />
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
  shipment: get(state, 'tsp.shipment', {}),
  loadTspDependenciesHasSuccess: get(
    state,
    'tsp.loadTspDependenciesHasSuccess',
  ),
  loadTspDependenciesHasError: get(state, 'tsp.loadTspDependenciesHasError'),
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      loadShipmentDependencies,
    },
    dispatch,
  );

export default withContext(
  connect(mapStateToProps, mapDispatchToProps)(ShipmentInfo),
);
