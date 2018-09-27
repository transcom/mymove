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
    const { moveDocuments, shipment } = this.props;
    return (
      <div>
        <p>{JSON.stringify(shipment)}</p>
        <p>{JSON.stringify(moveDocuments)}</p>
      </div>
    );
  }
}

MoveDocumentView.propTypes = {
  moveDocuments: PropTypes.array.isRequired,
  onDidMount: PropTypes.func.isRequired,
  shipment: PropTypes.object.isRequired,
};

export default MoveDocumentView;
