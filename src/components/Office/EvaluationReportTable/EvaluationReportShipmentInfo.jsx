import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './EvaluationReportShipmentInfo.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatEvaluationReportShipmentAddress, formatShortShipmentID } from 'utils/formatters';
import { ShipmentShape } from 'types/shipment';

const EvaluationReportShipmentInfo = ({ shipment }) => {
  let heading;
  let pickupAddress;
  let destinationAddress;
  let shipmentAccentStyle;

  switch (shipment.shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      heading = 'HHG';
      pickupAddress = formatEvaluationReportShipmentAddress(shipment.pickupAddress);
      destinationAddress = formatEvaluationReportShipmentAddress(shipment.destinationAddress);
      shipmentAccentStyle = styles.hhgShipmentType;
      break;
    case SHIPMENT_OPTIONS.NTS:
      heading = 'NTS';
      pickupAddress = formatEvaluationReportShipmentAddress(shipment.pickupAddress);
      destinationAddress = shipment.storageFacility.facilityName;
      shipmentAccentStyle = styles.ntsShipmentType;
      break;
    case SHIPMENT_OPTIONS.NTSR:
      heading = 'NTS-Release';
      pickupAddress = shipment.storageFacility.facilityName;
      destinationAddress = formatEvaluationReportShipmentAddress(shipment.destinationAddress);
      shipmentAccentStyle = styles.ntsrShipmentType;
      break;
    case SHIPMENT_OPTIONS.PPM:
      heading = 'PPM';
      pickupAddress = shipment.ppmShipment.pickupPostalCode;
      destinationAddress = shipment.ppmShipment.destinationPostalCode;
      shipmentAccentStyle = styles.ppmShipmentType;
      break;
    default:
      break;
  }

  return (
    <>
      <div className={classnames(styles.shipmentAccent, shipmentAccentStyle)} />
      <div className={styles.shipmentInfoContainer}>
        <div className={styles.shipmentInfo}>
          <h4>
            {heading} Shipment ID {formatShortShipmentID(shipment.id)}
          </h4>
          <small>
            {pickupAddress} <FontAwesomeIcon icon="arrow-right" /> {destinationAddress}
          </small>
        </div>
        <Button>Create report</Button>
      </div>
    </>
  );
};
EvaluationReportShipmentInfo.propTypes = {
  shipment: ShipmentShape.isRequired,
};

export default EvaluationReportShipmentInfo;
