import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';
import { useMutation, queryCache } from 'react-query';
import { useLocation, useHistory, useParams } from 'react-router-dom';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatEvaluationReportShipmentAddress } from 'utils/formatters';
import { ShipmentShape } from 'types/shipment';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';
import { createEvaluationReportForShipment } from 'services/ghcApi';
import { SHIPMENT_EVALUATION_REPORTS } from 'constants/queryKeys';

const EvaluationReportShipmentInfo = ({ shipment, shipmentNumber }) => {
  const { moveCode } = useParams();
  const location = useLocation();
  const history = useHistory();

  const [createReportMutation] = useMutation(createEvaluationReportForShipment, {
    onSuccess: () => {
      queryCache.invalidateQueries([SHIPMENT_EVALUATION_REPORTS, moveCode]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const handleCreateClick = async (shipmentID) => {
    const report = await createReportMutation({ body: { shipmentID } });
    const reportId = report?.id;

    history.push(`${location.pathname}/${reportId}`);
  };

  let heading;
  let pickupAddress;
  let destinationAddress;

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
      <Button onClick={() => handleCreateClick(shipment.id)}>Create report</Button>
    </>
  );
};
EvaluationReportShipmentInfo.propTypes = {
  shipment: ShipmentShape.isRequired,
  shipmentNumber: PropTypes.number.isRequired,
};

export default EvaluationReportShipmentInfo;
