import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { get } from 'lodash';

import ShipmentContainer from '../../components/Office/ShipmentContainer';
import ShipmentHeading from '../../components/Office/ShipmentHeading';
import ImportantShipmentDates from '../../components/Office/ImportantShipmentDates';
import RequestedServiceItemsTable from '../../components/Office/RequestedServiceItemsTable';
import { getMTOShipments, selectMTOShipmentsByMTOId } from '../../shared/Entities/modules/mtoShipments';
import { getMTOServiceItems, selectMTOServiceItemsByMTOId } from '../../shared/Entities/modules/mtoServiceItems';

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
    this.props.getMTOShipments(moveTaskOrderId);
    this.props.getMTOServiceItems(moveTaskOrderId);
    /* eslint-enable react/prop-types,react/destructuring-assignment */
  }

  render() {
    // eslint-disable-next-line react/prop-types
    const { mtoShipments, mtoServiceItems } = this.props;

    return (
      <div style={{ display: 'flex' }}>
        <div className="" style={{ width: '85%' }} data-testid="too-shipment-container">
          {/* eslint-disable-next-line react/prop-types */}
          {mtoShipments.map((mtoShipment) => {
            return (
              <ShipmentContainer>
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
                <ImportantShipmentDates
                  requestedPickupDate={formatShipmentDate(mtoShipment.requestedPickupDate)}
                  scheduledPickupDate={formatShipmentDate(mtoShipment.scheduledPickupDate)}
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
    mtoShipments: selectMTOShipmentsByMTOId(state, moveTaskOrderId),
    mtoServiceItems,
  };
};

const mapDispatchToProps = {
  getMTOShipments,
  getMTOServiceItems,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveTaskOrder));
