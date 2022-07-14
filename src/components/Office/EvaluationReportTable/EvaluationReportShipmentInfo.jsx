import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatEvaluationReportShipmentAddress } from 'utils/formatters';
import { ShipmentShape } from 'types/shipment';

const EvaluationReportShipmentInfo = ({ shipment, shipmentNumber }) => {
  let heading = '????';
  let pickupAddress = '????';
  let destinationAddress = '????';

  // TODO
  switch (shipment.shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      heading = 'HHG';
      pickupAddress = formatEvaluationReportShipmentAddress(shipment.pickupAddress);
      destinationAddress = formatEvaluationReportShipmentAddress(shipment.destinationAddress);
      break;
    case SHIPMENT_OPTIONS.NTS:
      heading = 'NTS';
      pickupAddress = formatEvaluationReportShipmentAddress(shipment.pickupAddress);
      destinationAddress = shipment.storageFacility.facilityName;
      break;
    case SHIPMENT_OPTIONS.NTSR:
      heading = 'NTS-Release';
      pickupAddress = shipment.storageFacility.facilityName;
      destinationAddress = formatEvaluationReportShipmentAddress(shipment.destinationAddress);
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
        {heading} Shipment ID #{shipmentNumber}
      </h4>
      <small>
        {pickupAddress} <FontAwesomeIcon icon="arrow-right" /> {destinationAddress}
      </small>
      <Button>Create report</Button>
    </>
  );
};
EvaluationReportShipmentInfo.propTypes = {
  shipment: ShipmentShape.isRequired,
  shipmentNumber: PropTypes.number.isRequired,
};

export default EvaluationReportShipmentInfo;
