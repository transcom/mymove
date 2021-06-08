import React from 'react';
import * as PropTypes from 'prop-types';

import { formatDate } from 'shared/dates';
import { AddressShape } from 'types';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';

const ShipmentDetailsMain = ({ className, shipment, order }) => {
  return (
    <div className={className}>
      <ImportantShipmentDates
        requestedPickupDate={formatDate(shipment.requestedPickupDate)}
        scheduledPickupDate={formatDate(shipment.scheduledPickupDate)}
      />
      <ShipmentAddresses
        pickupAddress={shipment.pickupAddress}
        destinationAddress={shipment.destinationAddress || order.destinationDutyStationAddress?.postal_code}
        originDutyStation={order.originDutyStationAddress}
        destinationDutyStation={order.destinationDutyStationAddress}
      />
      <ShipmentWeightDetails estimatedWeight={shipment.estimatedWeight} actualWeight={shipment.actualWeight} />
    </div>
  );
};

ShipmentDetailsMain.propTypes = {
  className: PropTypes.string,
  shipment: PropTypes.shape({
    requestedPickupDate: PropTypes.string,
    scheduledPickupDate: PropTypes.string,
    pickupAddress: AddressShape,
    destinationAddress: AddressShape,
    estimatedWeight: PropTypes.number,
    actualWeight: PropTypes.number,
  }).isRequired,
  order: PropTypes.shape({
    originDutyStationAddress: AddressShape,
    destinationDutyStationAddress: AddressShape,
  }).isRequired,
};

ShipmentDetailsMain.defaultProps = {
  className: '',
};

export default ShipmentDetailsMain;
