import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, capitalize } from 'lodash';

import { NavLink } from 'react-router-dom';

import Alert from 'shared/Alert'; // eslint-disable-line
import { withContext } from 'shared/AppContext';

import {
  loadShipmentDependencies,
  acceptShipment,
  rejectShipment,
} from './ducks';
import { formatDate } from 'shared/formatters';

class AcceptPanel extends Component {
  state = {
    displayState: 'Button',
    'origin-agent-name': '',
    'origin-agent-phone-number': '',
    'origin-agent-email': '',
    'destination-agent-name': '',
    'destination-agent-phone-number': '',
    'destination-agent-email': '',
  };

  setConfirmState = () => {
    this.setState({ displayState: 'Confirm' });
  };

  setButtonState = () => {
    this.setState({ displayState: 'Button' });
  };

  handleChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };

  acceptShipment = event => {
    event.preventDefault();

    var originShippingAgent = {
      name: this.state['origin-agent-name'],
      phone_number: this.state['origin-agent-phone-number'],
      email: this.state['origin-agent-email'],
    };
    var destinationShippingAgent = {
      name: this.state['destination-agent-name'],
      phone_number: this.state['destination-agent-phone-number'],
      email: this.state['destination-agent-email'],
    };
    console.log(originShippingAgent, destinationShippingAgent);
    this.props.acceptShipment(originShippingAgent, destinationShippingAgent);
    this.setState({ displayState: 'Button' });
  };

  render() {
    if (this.state.displayState === 'Confirm') {
      return (
        <div className="cancel-panel">
          <h2 className="extras usa-heading">Accept Shipment</h2>
          <div className="extras content">
            Enter the Origin and Destination Shipping Agent Information
            <div>
              <label htmlFor="origin-agent-name">Origin Agent Name</label>
              <input
                id="origin-agent-name"
                name="origin-agent-name"
                type="text"
                onChange={this.handleChange}
              />
              <label htmlFor="origin-agent-phone-number">
                Origin Agent Phone
              </label>
              <input
                id="origin-agent-phone-number"
                name="origin-agent-phone-number"
                type="text"
                onChange={this.handleChange}
              />
              <label htmlFor="origin-agent-email">Origin Agent Email</label>
              <input
                id="origin-agent-email"
                name="origin-agent-email"
                type="text"
                onChange={this.handleChange}
              />
            </div>
            <div>
              <label htmlFor="destination-agent-name">
                Destination Agent Name
              </label>
              <input
                id="destination-agent-name"
                name="destination-agent-name"
                type="text"
                onChange={this.handleChange}
              />
              <label htmlFor="destination-agent-phone-number">
                Destination Agent Phone
              </label>
              <input
                id="destination-agent-phone-number"
                name="destination-agent-phone-number"
                type="text"
                onChange={this.handleChange}
              />
              <label htmlFor="destination-agent-email">
                Destination Agent Email
              </label>
              <input
                id="destination-agent-email"
                name="destination-agent-email"
                type="text"
                onChange={this.handleChange}
              />
            </div>
            <div className="usa-grid">
              <div className="usa-width-one-whole extras options">
                <a onClick={this.setButtonState}>Never mind</a>
              </div>
              <div className="usa-width-one-whole extras options">
                <button onClick={this.acceptShipment}>Submit</button>
              </div>
            </div>
          </div>
        </div>
      );
    } else if (this.state.displayState === 'Button') {
      return (
        <button className="usa-button-primary" onClick={this.setConfirmState}>
          Accept Shipment
        </button>
      );
    }
  }
}

class RejectPanel extends Component {
  state = {
    displayState: 'Button',
    rejectReason: '',
  };

  setConfirmState = () => {
    this.setState({ displayState: 'Confirm' });
  };

  setRejectState = () => {
    if (this.state.rejectReason !== '') {
      this.setState({ displayState: 'Reject' });
    }
  };

  setButtonState = () => {
    this.setState({ displayState: 'Button' });
  };

  handleChange = event => {
    this.setState({ rejectReason: event.target.value });
  };

  rejectShipment = event => {
    event.preventDefault();
    this.props.rejectShipment(this.state.rejectReason);
    this.setState({ displayState: 'Button' });
  };

  render() {
    if (this.state.displayState === 'Reject') {
      return (
        <div className="reject-panel">
          <h2 className="extras usa-heading">Reject Shipment</h2>
          <div className="extras content">
            <Alert type="warning" heading="Rejection Warning">
              Are you sure you want to reject the entire move? This will affect
              your quality score.
            </Alert>
            <div className="usa-grid">
              <div className="usa-width-one-whole extras options">
                <a onClick={this.setButtonState}>No, never mind</a>
              </div>
              <div className="usa-width-one-whole extras options">
                <button onClick={this.rejectMove}>Yes, reject shipment</button>
              </div>
            </div>
          </div>
        </div>
      );
    } else if (this.state.displayState === 'Confirm') {
      return (
        <div className="reject-panel">
          <h2 className="extras usa-heading">Reject Shipment</h2>
          <div className="extras content">
            Why is the shipment being rejected?
            <textarea required onChange={this.handleChange} />
            <div className="usa-grid">
              <div className="usa-width-one-whole extras options">
                <a onClick={this.setButtonState}>Never mind</a>
              </div>
              <div className="usa-width-one-whole extras options">
                <button onClick={this.setRejectState}>Reject shipment</button>
              </div>
            </div>
          </div>
        </div>
      );
    } else if (this.state.displayState === 'Button') {
      return (
        <button className="usa-button-secondary" onClick={this.setConfirmState}>
          Reject Shipment
        </button>
      );
    }
  }
}

class ShipmentInfo extends Component {
  componentDidMount() {
    this.props.loadShipmentDependencies(this.props.match.params.shipmentId);
  }

  acceptShipment = (originShippingAgent, destinationShippingAgent) => {
    this.props.acceptShipment(
      this.props.shipment.id,
      originShippingAgent,
      destinationShippingAgent,
    );
  };

  rejectShipment = rejectReason => {
    this.props.rejectShipment(this.props.shipment.id, rejectReason);
  };

  render() {
    var move = this.props.shipment.move;
    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds Todo-phase2">
            {/* This comes from the Server Member model which is not yet on Shipments */}
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
            <ul className="move-info-header-meta">
              {/* Not clear where this comes from yet */}
              <li className="Todo-phase2">GBL# KKFA9999999</li>
              <li>Locator# {move && move.locator}</li>
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
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds Todo-phase2">
            <div className="infoPanel-Header">Dates &amp; Tracking</div>
            <div className="infoPanel-Body usa-grid-full">
              <div className="usa-width-one-half" />
              <div className="usa-width-one-half" />
            </div>
          </div>
          <div className="usa-width-one-third" />
          <AcceptPanel acceptShipment={this.acceptShipment} />
          <RejectPanel rejectShipment={this.rejectShipment} />
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
      acceptShipment,
      rejectShipment,
    },
    dispatch,
  );

export default withContext(
  connect(mapStateToProps, mapDispatchToProps)(ShipmentInfo),
);
