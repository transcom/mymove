import React, { Component } from 'react';
import PropTypes from 'prop-types';

class MoveDocumentView extends Component {
  componentDidMount() {
    const {
      onDidMount,
      match: {
        params: { shipmentId },
      },
    } = this.props;
    onDidMount(shipmentId);
  }

  render() {
    const { shipment } = this.props;
    return <div>{JSON.stringify(shipment)}</div>;
  }
}

MoveDocumentView.propTypes = {
  onDidMount: PropTypes.func.isRequired,
  shipment: PropTypes.object.isRequired,
};

export default MoveDocumentView;
