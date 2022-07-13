import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatPrimeAPIFullAddress } from 'utils/formatters';
import { ShipmentShape } from 'types/shipment';

const EvaluationReportShipmentInfo = ({ shipment }) => {
  let heading = '????';
  let pickupAddress = '????';
  let destinationAddress = '????';

  // TODO
  switch (shipment.shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      heading = 'HHG';
      pickupAddress = formatPrimeAPIFullAddress(shipment.pickupAddress);
      destinationAddress = formatPrimeAPIFullAddress(shipment.destinationAddress);
      break;
    case SHIPMENT_OPTIONS.NTS:
      heading = 'NTS';
      pickupAddress = formatPrimeAPIFullAddress(shipment.pickupAddress);
      destinationAddress = shipment.storageFacility.facilityName;
      break;
    case SHIPMENT_OPTIONS.NTSR:
      heading = 'NTS-Release';
      pickupAddress = shipment.storageFacility.facilityName;
      destinationAddress = formatPrimeAPIFullAddress(shipment.destinationAddress);
      break;
    case SHIPMENT_OPTIONS.PPM:
      heading = 'PPM';
      pickupAddress = shipment.ppmShipment.pickupPostalCode;
      destinationAddress = shipment.ppmShipment.destinationPostalCode;
      break;
    default:
      heading = '????!!!';
      break;
  }

  return (
    <>
      <div />
      <h4>
        {heading} Shipment ID #{shipment.id}
      </h4>
      <small>
        {pickupAddress} <FontAwesomeIcon icon="arrow-right" /> {destinationAddress}
      </small>
    </>
  );
};
EvaluationReportShipmentInfo.propTypes = {
  shipment: ShipmentShape,
};

EvaluationReportShipmentInfo.defaultProps = {
  shipment: {},
};

export default EvaluationReportShipmentInfo;
