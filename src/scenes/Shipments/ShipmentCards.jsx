// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import 'scenes/Shipments/ShipmentCards.css';

const ShipmentCards = ({ shipments }) => {
  if (!shipments) return <LoadingPlaceholder />;
  if (shipments.length === 0) {
    return <h2> There are no shipments at the moment! </h2>;
  }

  const cards = shipments.map(shipment => {
    let awardedStatus, tspID;
    if (shipment.transportation_service_provider_id) {
      awardedStatus = 'awarded';
      tspID = shipment.transportation_service_provider_id.substr(0, 6) + '...';
    } else {
      awardedStatus = 'available';
      tspID = '-';
    }
    const tdlID = shipment.traffic_distribution_list_id.substr(0, 6);
    const className = `shipment-card ${awardedStatus}`;

    return (
      <div key={shipment.id} className={className}>
        <b>
          Shipment: {shipment.id.substr(0, 6)}
          ...
        </b>
        TDL: {tdlID}
        ...
        <br />
        <br />
        Status: <b>{awardedStatus}</b>
        <br />
        TSP: {tspID}
        <br />
        <br />
        Pickup Date: {shipment.pickup_date}
        <br />
        Delivery Date: {shipment.delivery_date}
      </div>
    );
  });

  return <div className="shipment-cards">{cards}</div>;
};

ShipmentCards.propTypes = {
  shipments: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      traffic_distribution_list_id: PropTypes.string.isRequired,
      pickup_date: PropTypes.string.isRequired,
      delivery_date: PropTypes.string.isRequired,
      shipment_id: PropTypes.string,
      transportation_service_provider_id: PropTypes.string,
      administrative_shipment: PropTypes.bool,
    }),
  ),
};

export default ShipmentCards;
