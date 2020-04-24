import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { get } from 'lodash';
import ShipmentContainer from '../../components/Office/ShipmentContainer';
import ShipmentHeading from '../../components/Office/ShipmentHeading';
import { getMTOShipments, selectMTOShiomentsByMTOId } from '../../shared/Entities/modules/mtoShipments';
import '../../index.scss';
import '../../ghc_index.scss';

function formatShipmentType(shipmentType) {
  if (shipmentType === 'HHG') {
    return 'Household Goods';
  }
  return shipmentType;
}

function formatShipmentDate(shipmentDateString) {
  const dateObj = new Date(shipmentDateString);
  const year = new Intl.DateTimeFormat('en', { year: 'numeric' }).format(dateObj);
  const month = new Intl.DateTimeFormat('en', { month: 'short' }).format(dateObj);
  const day = new Intl.DateTimeFormat('en', { day: '2-digit' }).format(dateObj);
  return `${day} ${month} ${year}`;
}

class MoveTaskOrder extends Component {
  componentDidMount() {
    // eslint-disable-next-line react/prop-types,react/destructuring-assignment
    const { moveTaskOrderId } = this.props.match.params;
    // eslint-disable-next-line react/prop-types,react/destructuring-assignment
    this.props.getMTOShipments(moveTaskOrderId);
  }

  render() {
    // eslint-disable-next-line react/prop-types
    const { mtoShipments } = this.props;

    return (
      <div className="grid-container">
        {/* eslint-disable-next-line react/prop-types */}
        {mtoShipments.map((mtoShipment) => {
          return (
            <ShipmentContainer data-cy="too-shipment-container">
              <ShipmentHeading
                key={mtoShipment.id}
                shipmentInfo={{
                  shipmentType: formatShipmentType(mtoShipment.shipmentType),
                  originCity: get(mtoShipment.pickupAddress, 'city'),
                  originState: get(mtoShipment.pickupAddress, 'state'),
                  originPostalCode: get(mtoShipment.pickupAddress, 'postal_code'),
                  destinationCity: get(mtoShipment.destinationAddress, 'city'),
                  destinationState: get(mtoShipment.destinationAddress, 'state'),
                  destinationPostalCode: get(mtoShipment.destinationAddress, 'postal_code'),
                  scheduledPickupDate: formatShipmentDate(mtoShipment.scheduledPickupDate),
                }}
              />
            </ShipmentContainer>
          );
        })}
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const { moveTaskOrderId } = ownProps.match.params;
  return {
    mtoShipments: selectMTOShiomentsByMTOId(state, moveTaskOrderId),
  };
};

const mapDispatchToProps = {
  getMTOShipments,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveTaskOrder));
