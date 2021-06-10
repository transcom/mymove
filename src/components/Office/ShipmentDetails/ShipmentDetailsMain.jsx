import React from 'react';
import * as PropTypes from 'prop-types';

import { formatDate } from 'shared/dates';
import { AddressShape } from 'types';
import { ShipmentShape } from 'types/shipment';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';

const ShipmentDetailsMain = ({ className, shipment, dutyStationAddresses }) => {
  const {
    requestedPickupDate,
    scheduledPickupDate,
    pickupAddress,
    destinationAddress,
    primeEstimatedWeight,
    primeActualWeight,
  } = shipment;
  const { originDutyStationAddress, destinationDutyStationAddress } = dutyStationAddresses;
  return (
    <div className={className}>
      <ImportantShipmentDates
        requestedPickupDate={formatDate(requestedPickupDate)}
        scheduledPickupDate={formatDate(scheduledPickupDate)}
      />
      <ShipmentAddresses
        pickupAddress={pickupAddress}
        destinationAddress={destinationAddress || destinationDutyStationAddress?.postal_code}
        originDutyStation={originDutyStationAddress}
        destinationDutyStation={destinationDutyStationAddress}
      />
      <ShipmentWeightDetails estimatedWeight={primeEstimatedWeight} actualWeight={primeActualWeight} />
    </div>
  );
};

ShipmentDetailsMain.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  dutyStationAddresses: PropTypes.shape({
    originDutyStationAddress: AddressShape,
    destinationDutyStationAddress: AddressShape,
  }).isRequired,
};

ShipmentDetailsMain.defaultProps = {
  className: '',
};

export default ShipmentDetailsMain;
