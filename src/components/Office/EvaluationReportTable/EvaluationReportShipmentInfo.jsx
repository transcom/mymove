import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { SHIPMENT_OPTIONS } from '../../../shared/constants';

import { ShipmentShape } from 'types/shipment';

const EvaluationReportShipmentInfo = ({ shipment }) => {
  let heading = '????';
  switch (shipment.shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      heading = 'HHG';
      break;
    case SHIPMENT_OPTIONS.NTS:
      heading = 'NTS';
      break;
    case SHIPMENT_OPTIONS.NTSR:
      heading = 'NTS-Release';
      break;
    case SHIPMENT_OPTIONS.PPM:
      heading = 'PPM';
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
        origin <FontAwesomeIcon icon="arrow-right" /> destination
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
