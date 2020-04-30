import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { get } from 'lodash';
import ShipmentContainer from '../../components/Office/ShipmentContainer';
import ShipmentHeading from '../../components/Office/ShipmentHeading';
import ImportantShipmentDates from '../../components/Office/ImportantShipmentDates';
import RequestedServiceItemsTable from '../../components/Office/RequestedServiceItemsTable';
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
    // get service items
  }

  render() {
    // eslint-disable-next-line react/prop-types
    const { mtoShipments } = this.props;
    const serviceItems = [
      {
        id: 'abc-123',
        dateRequested: '20 Nov 2020',
        serviceItem: 'Dom. Origin 1st Day SIT',
        code: 'DOMSIT',
        details: {
          text: {
            ZIP: '60612',
            Reason: "here's the reason",
          },
          imgURL: null,
        },
      },
      {
        id: 'abc-1234',
        dateRequested: '22 Nov 2020',
        serviceItem: 'Dom. Destination 1st Day SIT',
        code: 'DDFSIT',
        details: {
          text: {
            'First available delivery date': '22 Nov 2020',
            'First customer contact': '22 Nov 2020 12:00pm',
            'Second customer contact': '22 Nov 2020 12:00pm',
          },
          imgURL: null,
        },
      },
      {
        id: 'cba-123',
        dateRequested: '22 Nov 2020',
        serviceItem: 'Dom. Origin Shuttle Service',
        code: 'DOSHUT',
        details: {
          text: {
            'Reason for request': "Here's the reason",
            'Estimated weight': '3,500lbs',
          },
          imgURL: null,
        },
      },
      {
        id: 'cba-1234',
        dateRequested: '22 Nov 2020',
        serviceItem: 'Dom. Destination Shuttle Service',
        code: 'DDSHUT',
        details: {
          text: {
            'Reason for request': "Here's the reason",
            'Estimated weight': '3,500lbs',
          },
          imgURL: null,
        },
      },
      {
        id: 'abc12345',
        dateRequested: '22 Nov 2020',
        serviceItem: 'Dom. Crating',
        code: 'DCRT',
        details: {
          text: {
            Description: "Here's the description",
            'Item dimensions': '84"x26"x42"',
            'Crate dimensions': '110"x36"x54"',
          },
          imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
        },
      },
    ];

    return (
      <div style={{ display: 'flex' }}>
        <div className="" style={{ width: '85%' }} data-cy="too-shipment-container">
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
                <RequestedServiceItemsTable serviceItems={serviceItems} />
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
  return {
    mtoShipments: selectMTOShiomentsByMTOId(state, moveTaskOrderId),
  };
};

const mapDispatchToProps = {
  getMTOShipments,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MoveTaskOrder));
