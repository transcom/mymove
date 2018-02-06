// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';

import Alert from 'shared/Alert';
import ShipmentCards from 'scenes/Shipments/ShipmentCards';

import { AvailableShipmentsIndex } from 'shared/api.js';

class AvailableShipments extends Component {
  constructor(props) {
    super(props);
    this.state = { shipments: null, hasError: false };
  }
  componentDidMount() {
    document.title = 'Transcom PPP: Available Shipments';
    this.loadAvailableShipments();
  }
  render() {
    const { shipments, hasError } = this.state;
    return (
      <div className="usa-grid">
        <h1>Available Shipments</h1>
        {hasError && (
          <Alert type="error" heading="Server Error">
            There was a problem loading the shipments from the server.
          </Alert>
        )}
        {!hasError && <ShipmentCards shipments={shipments} />}
      </div>
    );
  }
  loadAvailableShipments = async () => {
    try {
      const shipments = await AvailableShipmentsIndex();
      this.setState({ shipments });
    } catch (e) {
      //componentDidCatch will not get fired because this is async
      //todo: how to we want to monitor errors
      console.error(e);
      this.setState({ hasError: true });
    }
  };
}
export default AvailableShipments;
