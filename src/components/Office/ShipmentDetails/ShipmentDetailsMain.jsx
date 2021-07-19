import React from 'react';
import * as PropTypes from 'prop-types';

import { formatDate } from 'shared/dates';
import { AddressShape } from 'types';
import { ShipmentShape } from 'types/shipment';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates/ImportantShipmentDates';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';
import ShipmentRemarks from 'components/Office/ShipmentRemarks/ShipmentRemarks';

const ShipmentDetailsMain = ({ className, shipment, dutyStationAddresses, handleDivertShipment }) => {
  const {
    requestedPickupDate,
    scheduledPickupDate,
    pickupAddress,
    destinationAddress,
    primeEstimatedWeight,
    primeActualWeight,
    counselorRemarks,
    customerRemarks,
  } = shipment;
  const { originDutyStationAddress, destinationDutyStationAddress } = dutyStationAddresses;

  return (
    <div className={className}>
      <ImportantShipmentDates
        requestedPickupDate={formatDate(requestedPickupDate)}
        scheduledPickupDate={scheduledPickupDate ? formatDate(scheduledPickupDate) : null}
      />
      <ShipmentAddresses
        pickupAddress={pickupAddress}
        destinationAddress={destinationAddress || destinationDutyStationAddress?.postal_code}
        originDutyStation={originDutyStationAddress}
        destinationDutyStation={destinationDutyStationAddress}
        shipmentInfo={{ shipmentID: shipment.id, ifMatchEtag: shipment.eTag, shipmentStatus: shipment.status }}
        handleDivertShipment={handleDivertShipment}
      />
      <ShipmentWeightDetails estimatedWeight={primeEstimatedWeight} actualWeight={primeActualWeight} />
      {counselorRemarks && <ShipmentRemarks title="Counselor remarks" remarks={counselorRemarks} />}
      {customerRemarks && <ShipmentRemarks title="Customer remarks" remarks={customerRemarks} />}
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
  handleDivertShipment: PropTypes.func.isRequired,
};

ShipmentDetailsMain.defaultProps = {
  className: '',
};

export default ShipmentDetailsMain;
