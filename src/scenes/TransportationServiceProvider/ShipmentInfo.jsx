import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, capitalize } from 'lodash';

import { getFormValues, reduxForm } from 'redux-form';
import { NavLink } from 'react-router-dom';

import Alert from 'shared/Alert'; // eslint-disable-line
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { withContext } from 'shared/AppContext';

import { loadShipmentDependencies, acceptShipment } from './ducks';
import { formatDate } from 'shared/formatters';

const shipmentAcceptFormName = 'shipment_accept';

let ShipmentAcceptForm = props => {
  const { schema, setButtonState, acceptFormShipment } = props;

  return (
    <form onSubmit={acceptFormShipment}>
      <h3 className="smheading">Origin Shipping Agent</h3>
      <SwaggerField fieldName="origin_agent_name" swagger={schema} required />
      <SwaggerField
        fieldName="origin_agent_phone_number"
        swagger={schema}
        required
      />
      <SwaggerField fieldName="origin_agent_email" swagger={schema} required />

      <h3 className="smheading">Destination Shipping Agent</h3>
      <SwaggerField
        fieldName="destination_agent_name"
        swagger={schema}
        required
      />
      <SwaggerField
        fieldName="destination_agent_phone_number"
        swagger={schema}
        required
      />
      <SwaggerField
        fieldName="destination_agent_email"
        swagger={schema}
        required
      />

      <div className="usa-grid">
        <div className="usa-width-one-whole extras options">
          <a onClick={setButtonState}>Never mind</a>
        </div>
        <div className="usa-width-one-whole extras options">
          <button type="submit">Submit</button>
        </div>
      </div>
    </form>
  );
};

ShipmentAcceptForm = reduxForm({
  form: shipmentAcceptFormName,
})(ShipmentAcceptForm);

class AcceptPanel extends Component {
  state = {
    displayState: 'Button',
  };

  setConfirmState = () => {
    this.setState({ displayState: 'Confirm' });
  };

  setButtonState = () => {
    this.setState({ displayState: 'Button' });
  };

  acceptFormShipment = event => {
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
        <ShipmentAcceptForm
          schema={this.props.schema}
          setButtonState={this.setButtonState}
          acceptFormShipment={this.acceptFormShipment}
        />
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
          <AcceptPanel
            acceptShipment={this.acceptShipment}
            schema={this.props.schemaShipmentAccept}
          />
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
  formValues: getFormValues(shipmentAcceptFormName)(state),
  loadTspDependenciesHasError: get(state, 'tsp.loadTspDependenciesHasError'),
  schemaShipmentAccept: get(
    state,
    'swagger.spec.definitions.ShipmentAccept',
    {},
  ),
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      loadShipmentDependencies,
      acceptShipment,
    },
    dispatch,
  );

export default withContext(
  connect(mapStateToProps, mapDispatchToProps)(ShipmentInfo),
);
