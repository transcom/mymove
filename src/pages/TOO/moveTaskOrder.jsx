import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import ShipmentContainer from '../../components/Office/ShipmentContainer';
import ShipmentHeading from '../../components/Office/ShipmentHeading';
import ImportantShipmentDates from '../../components/Office/ImportantShipmentDates';
import RequestedServiceItemsTable from '../../components/Office/RequestedServiceItemsTable';
import { getMTOShipments, selectMTOShipmentsByMTOId } from '../../shared/Entities/modules/mtoShipments';
import { getMTOServiceItems, selectMTOServiceItemsByMTOId } from '../../shared/Entities/modules/mtoServiceItems';
import { getMoveTaskOrder } from '../../shared/Entities/modules/moveTaskOrders';

import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import { selectMoveTaskOrder } from 'shared/Entities/modules/moveTaskOrders';

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

    /* eslint-disable react/prop-types,react/destructuring-assignment */
    this.props.getMoveTaskOrder(moveTaskOrderId);
    this.props.getMTOShipments(moveTaskOrderId);
    this.props.getMTOServiceItems(moveTaskOrderId);
    /* eslint-enable react/prop-types,react/destructuring-assignment */
  }

  render() {
    // eslint-disable-next-line react/prop-types
    const { moveTaskOrder, mtoShipments, mtoServiceItems } = this.props;

    return (
      <div style={{ display: 'flex' }}>
        <div className="" style={{ width: '85%' }} data-testid="too-shipment-container">
          {/* eslint-disable-next-line react/prop-types */}
          {mtoShipments.map((mtoShipment) => {
            const {
              shipmentType,
              pickupAddress,
              destinationAddress,
              scheduledPickupDate,
              requestedPickupDate,
            } = mtoShipment;
            return (
              <ShipmentContainer>
                <ShipmentHeading
                  key={mtoShipment.id}
                  shipmentInfo={{
                    shipmentType: formatShipmentType(shipmentType),
                    originCity: pickupAddress?.city,
                    originState: pickupAddress?.state,
                    // eslint-disable-next-line camelcase
                    originPostalCode: pickupAddress?.postal_code,
                    destinationCity: destinationAddress?.city,
                    destinationState: destinationAddress?.state,
                    // eslint-disable-next-line camelcase
                    destinationPostalCode: destinationAddress?.postal_code,
                    scheduledPickupDate: formatShipmentDate(scheduledPickupDate),
                  }}
                />
                <ImportantShipmentDates
                  requestedPickupDate={formatShipmentDate(requestedPickupDate)}
                  scheduledPickupDate={formatShipmentDate(scheduledPickupDate)}
                />
                <ShipmentAddresses
                  pickupAddress={pickupAddress}
                  destinationAddress={destinationAddress}
                  // eslint-disable-next-line react/prop-types
                  originDutyStation={moveTaskOrder.originDutyStation}
                  // eslint-disable-next-line react/prop-types
                  destinationDutyStation={moveTaskOrder.destinationDutyStation}
                />
                <RequestedServiceItemsTable serviceItems={mtoServiceItems} />
              </ShipmentContainer>
            );
          })}
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const { moveTaskOrderId } = ownProps.match.params;
  const mtoServiceItems = selectMTOServiceItemsByMTOId(state, moveTaskOrderId).map((item) => {
    const detailText = { ZIP: item.pickupPostalCode, Reason: item.reason };
    /* eslint-disable no-param-reassign */
    item.serviceItem = item.reServiceName;
    item.details = { text: detailText, imgURL: '' };
    /* eslint-enable no-param-reassign */
    return item;
  });

  return {
    moveTaskOrder: selectMoveTaskOrder(state, moveTaskOrderId),
    mtoShipments: selectMTOShipmentsByMTOId(state, moveTaskOrderId),
    mtoServiceItems,
  };
};

const mapDispatchToProps = {
  getMoveTaskOrder,
  getMTOShipments,
  getMTOServiceItems,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveTaskOrder));
