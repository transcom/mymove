import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { NavLink } from 'react-router-dom';

import { withContext } from 'shared/AppContext';

import { loadShipmentDependencies } from './ducks';
import { formatDate } from 'shared/formatters';

class ShipmentInfo extends Component {
  componentDidMount() {
    this.props.loadShipmentDependencies(this.props.match.params.shipmentId);
  }

  render() {
    var move = this.props.shipment.move;
    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds">
            <h1>Shipment Info: LastName, FirstName</h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Shipments Queue</span>
            </NavLink>
          </div>
        </div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-one-whole">
            <ul className="move-info-header-meta Todo-phase2">
              <li>GBL# {this.props.shipment.source_gbloc}</li>
              <li>Locator# {move && move.locator}</li>
              <li>KKFA to HAFC</li>
              <li>
                Requested Move date{' '}
                {formatDate(this.props.shipment.requested_pickup_date)}
              </li>
              <li>
                Status: <b>At Destination</b>
              </li>
            </ul>
          </div>
        </div>
        <div className="usa-grid grid-wide tabs">
          <div className="usa-width-two-thirds">
            <p>
              <button className="usa-button-primary">Accept</button>
              <button className="usa-button-secondary">Reject</button>
            </p>
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
